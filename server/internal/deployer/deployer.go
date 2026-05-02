package deployer

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	alicas "github.com/alibabacloud-go/cas-20200407/v4/client"
	aliopen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/dara"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/pkg/sftp"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/client"
	gossh "golang.org/x/crypto/ssh"
)

const (
	ProviderAliyunCAS = "aliyun-cas"
	ProviderQiniu     = "qiniu"
	ProviderSSH       = "ssh"
)

type Request struct {
	Provider     string
	AccessConfig map[string]interface{}
	Config       map[string]interface{}
	Certificate  string
	PrivateKey   string
}

func Execute(ctx context.Context, req Request) error {
	switch strings.ToLower(strings.TrimSpace(req.Provider)) {
	case ProviderAliyunCAS:
		return deployAliyunCAS(ctx, req)
	case ProviderQiniu:
		return deployQiniu(ctx, req)
	case ProviderSSH:
		return deploySSH(ctx, req)
	default:
		return fmt.Errorf("unsupported deploy provider: %s", req.Provider)
	}
}

func deployQiniu(ctx context.Context, req Request) error {
	accessKey := readString(req.AccessConfig, "accessKey")
	secretKey := readString(req.AccessConfig, "secretKey")
	if accessKey == "" {
		return fmt.Errorf("qiniu.accessKey is required")
	}
	if secretKey == "" {
		return fmt.Errorf("qiniu.secretKey is required")
	}

	certName := readString(req.Config, "certName")
	if certName == "" {
		certName = fmt.Sprintf("easyssl-%d", time.Now().UnixMilli())
	}

	commonName := readString(req.Config, "commonName")
	if commonName == "" {
		if cn, err := readCommonName(req.Certificate); err == nil {
			commonName = cn
		}
	}
	if commonName == "" {
		return fmt.Errorf("qiniu commonName is required")
	}

	mac := auth.New(accessKey, secretKey)
	httpClient := &client.Client{Client: &http.Client{Transport: newQiniuTransport(mac, nil)}}

	payload := map[string]string{
		"name":        certName,
		"common_name": commonName,
		"ca":          req.Certificate,
		"pri":         req.PrivateKey,
	}
	resp := struct {
		Code   *int   `json:"code,omitempty"`
		Error  string `json:"error,omitempty"`
		CertID string `json:"certID,omitempty"`
	}{}
	if err := httpClient.CallWithJson(ctx, &resp, "POST", "https://fusion.qiniuapi.com/sslcert", nil, payload); err != nil {
		return fmt.Errorf("upload certificate to qiniu failed: %w", err)
	}

	// Qiniu cert APIs may return either:
	// - upload success: {"certID":"..."}
	// - common envelope: {"code":200,"error":""}
	if strings.TrimSpace(resp.CertID) != "" {
		return nil
	}
	if resp.Code != nil && *resp.Code != 0 && *resp.Code != 200 {
		return fmt.Errorf("upload certificate to qiniu failed: code=%d error=%s", *resp.Code, strings.TrimSpace(resp.Error))
	}
	if resp.Code == nil && strings.TrimSpace(resp.Error) != "" {
		return fmt.Errorf("upload certificate to qiniu failed: error=%s", strings.TrimSpace(resp.Error))
	}
	return nil
}

func deployAliyunCAS(ctx context.Context, req Request) error {
	accessKeyID := readString(req.AccessConfig, "accessKeyId")
	accessKeySecret := readString(req.AccessConfig, "accessKeySecret")
	resourceGroupID := readString(req.Config, "resourceGroupId")
	if resourceGroupID == "" {
		resourceGroupID = readString(req.AccessConfig, "resourceGroupId")
	}
	region := readString(req.Config, "region")
	if region == "" {
		region = readString(req.AccessConfig, "region")
	}

	if accessKeyID == "" {
		return fmt.Errorf("aliyun.accessKeyId is required")
	}
	if accessKeySecret == "" {
		return fmt.Errorf("aliyun.accessKeySecret is required")
	}

	endpoint := "cas.aliyuncs.com"
	if region != "" && region != "cn-hangzhou" {
		endpoint = fmt.Sprintf("cas.%s.aliyuncs.com", region)
	}

	client, err := alicas.NewClient(&aliopen.Config{
		Endpoint:        tea.String(endpoint),
		AccessKeyId:     tea.String(accessKeyID),
		AccessKeySecret: tea.String(accessKeySecret),
	})
	if err != nil {
		return fmt.Errorf("create aliyun cas client failed: %w", err)
	}

	uploadReq := &alicas.UploadUserCertificateRequest{
		Name: tea.String(fmt.Sprintf("easyssl_%d", time.Now().UnixMilli())),
		Cert: tea.String(req.Certificate),
		Key:  tea.String(req.PrivateKey),
	}
	if resourceGroupID != "" {
		uploadReq.ResourceGroupId = tea.String(resourceGroupID)
	}
	if _, err := client.UploadUserCertificateWithContext(ctx, uploadReq, &dara.RuntimeOptions{}); err != nil {
		return fmt.Errorf("upload certificate to aliyun cas failed: %w", err)
	}
	return nil
}

func deploySSH(ctx context.Context, req Request) error {
	_ = ctx
	host := readString(req.AccessConfig, "host")
	if host == "" {
		host = readString(req.AccessConfig, "sshHost")
	}
	port := readInt(req.AccessConfig, "port")
	if port == 0 {
		port = readInt(req.AccessConfig, "sshPort")
	}
	if port == 0 {
		port = 22
	}
	username := readString(req.AccessConfig, "username")
	if username == "" {
		username = readString(req.AccessConfig, "sshUsername")
	}
	if username == "" {
		username = "root"
	}

	if host == "" {
		return fmt.Errorf("ssh.host is required")
	}

	authMethod := strings.ToLower(readString(req.AccessConfig, "authMethod"))
	if authMethod == "" {
		authMethod = strings.ToLower(readString(req.AccessConfig, "sshAuthMethod"))
	}
	password := readString(req.AccessConfig, "password")
	if password == "" {
		password = readString(req.AccessConfig, "sshPassword")
	}
	privateKey := readString(req.AccessConfig, "key")
	if privateKey == "" {
		privateKey = readString(req.AccessConfig, "sshKey")
	}
	keyPassphrase := readString(req.AccessConfig, "keyPassphrase")
	if keyPassphrase == "" {
		keyPassphrase = readString(req.AccessConfig, "sshKeyPassphrase")
	}

	if authMethod == "" {
		if privateKey != "" {
			authMethod = "key"
		} else {
			authMethod = "password"
		}
	}

	var auth []gossh.AuthMethod
	switch authMethod {
	case "password":
		if password == "" {
			return fmt.Errorf("ssh.password is required")
		}
		auth = append(auth, gossh.Password(password))
	case "key":
		if privateKey == "" {
			return fmt.Errorf("ssh.key is required")
		}
		signer, err := parseSSHKey(privateKey, keyPassphrase)
		if err != nil {
			return fmt.Errorf("parse ssh private key failed: %w", err)
		}
		auth = append(auth, gossh.PublicKeys(signer))
	default:
		return fmt.Errorf("unsupported ssh authMethod: %s", authMethod)
	}

	clientCfg := &gossh.ClientConfig{
		User:            username,
		Auth:            auth,
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
	}

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	sshClient, err := gossh.Dial("tcp", addr, clientCfg)
	if err != nil {
		return fmt.Errorf("ssh dial failed: %w", err)
	}
	defer sshClient.Close()

	if pre := readString(req.Config, "preCommand"); pre != "" {
		if err := runSSHCommand(sshClient, pre); err != nil {
			return fmt.Errorf("run preCommand failed: %w", err)
		}
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("create sftp client failed: %w", err)
	}
	defer sftpClient.Close()

	certPath := readString(req.Config, "certPath")
	keyPath := readString(req.Config, "keyPath")

	if certPath == "" {
		certPath = "/etc/nginx/ssl/fullchain.pem"
	}
	if keyPath == "" {
		keyPath = "/etc/nginx/ssl/privkey.pem"
	}

	if err := writeRemoteFile(sftpClient, certPath, req.Certificate); err != nil {
		return fmt.Errorf("write cert file failed: %w", err)
	}
	if err := writeRemoteFile(sftpClient, keyPath, req.PrivateKey); err != nil {
		return fmt.Errorf("write key file failed: %w", err)
	}

	if post := readString(req.Config, "postCommand"); post != "" {
		if err := runSSHCommand(sshClient, post); err != nil {
			return fmt.Errorf("run postCommand failed: %w", err)
		}
	}

	return nil
}

func parseSSHKey(key, passphrase string) (gossh.Signer, error) {
	if passphrase == "" {
		return gossh.ParsePrivateKey([]byte(key))
	}
	return gossh.ParsePrivateKeyWithPassphrase([]byte(key), []byte(passphrase))
}

func runSSHCommand(client *gossh.Client, command string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	return session.Run(command)
}

func writeRemoteFile(client *sftp.Client, path string, content string) error {
	f, err := client.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write([]byte(content)); err != nil {
		return err
	}
	return nil
}

func readString(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%v", v))
}

func readInt(m map[string]interface{}, key string) int {
	if m == nil {
		return 0
	}
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}
	s := strings.TrimSpace(fmt.Sprintf("%v", v))
	n, _ := strconv.Atoi(s)
	return n
}

func readCommonName(certPEM string) (string, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return "", fmt.Errorf("invalid cert pem")
	}
	c, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(c.Subject.CommonName), nil
}

type qiniuTransport struct {
	http.RoundTripper
	mac *auth.Credentials
}

func newQiniuTransport(mac *auth.Credentials, tr http.RoundTripper) *qiniuTransport {
	if tr == nil {
		tr = http.DefaultTransport
	}
	return &qiniuTransport{RoundTripper: tr, mac: mac}
}

func (t *qiniuTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := t.mac.SignRequestV2(req)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Qiniu "+token)
	return t.RoundTripper.RoundTrip(req)
}
