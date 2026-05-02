package accessprovider

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"easyssl/server/internal/model"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/client"
	tccommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcerr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tcprofile "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcdnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	gossh "golang.org/x/crypto/ssh"
)

const (
	ProviderAliyun       = "aliyun"
	ProviderTencentCloud = "tencentcloud"
	ProviderQiniu        = "qiniu"
	ProviderSSH          = "ssh"
)

func stringify(v interface{}) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%v", v))
}

func normalizeProvider(provider string) string {
	return strings.ToLower(strings.TrimSpace(provider))
}

func TestAccess(ctx context.Context, access model.Access) error {
	switch normalizeProvider(access.Provider) {
	case ProviderAliyun:
		return testAliyunAccess(ctx, access.Config)
	case ProviderTencentCloud:
		return testTencentCloudAccess(ctx, access.Config)
	case ProviderQiniu:
		return testQiniuAccess(ctx, access.Config)
	case ProviderSSH:
		return testSSHAccess(ctx, access.Config)
	default:
		return fmt.Errorf("unsupported provider: %s", access.Provider)
	}
}

func testAliyunAccess(ctx context.Context, config map[string]interface{}) error {
	_ = ctx

	accessKeyID := stringify(config["accessKeyId"])
	accessKeySecret := stringify(config["accessKeySecret"])
	if accessKeyID == "" {
		return fmt.Errorf("aliyun.accessKeyId is required")
	}
	if accessKeySecret == "" {
		return fmt.Errorf("aliyun.accessKeySecret is required")
	}

	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", accessKeyID, accessKeySecret)
	if err != nil {
		return fmt.Errorf("create aliyun dns client failed: %w", err)
	}

	req := alidns.CreateDescribeDomainsRequest()
	req.PageSize = "1"
	_, err = client.DescribeDomains(req)
	if err != nil {
		return fmt.Errorf("aliyun access test failed: %w", err)
	}
	return nil
}

func testSSHAccess(ctx context.Context, config map[string]interface{}) error {
	host := stringify(config["host"])
	if host == "" {
		host = stringify(config["sshHost"])
	}
	if host == "" {
		return fmt.Errorf("ssh.host is required")
	}

	port := stringify(config["port"])
	if port == "" {
		port = stringify(config["sshPort"])
	}
	if port == "" {
		port = "22"
	}

	username := stringify(config["username"])
	if username == "" {
		username = stringify(config["sshUsername"])
	}
	if username == "" {
		username = "root"
	}

	authMethod := strings.ToLower(stringify(config["authMethod"]))
	if authMethod == "" {
		authMethod = strings.ToLower(stringify(config["sshAuthMethod"]))
	}
	password := stringify(config["password"])
	if password == "" {
		password = stringify(config["sshPassword"])
	}
	key := stringify(config["key"])
	if key == "" {
		key = stringify(config["sshKey"])
	}
	keyPassphrase := stringify(config["keyPassphrase"])
	if keyPassphrase == "" {
		keyPassphrase = stringify(config["sshKeyPassphrase"])
	}

	if authMethod == "" {
		if key != "" {
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
		if key == "" {
			return fmt.Errorf("ssh.key is required")
		}
		var signer gossh.Signer
		var err error
		if keyPassphrase == "" {
			signer, err = gossh.ParsePrivateKey([]byte(key))
		} else {
			signer, err = gossh.ParsePrivateKeyWithPassphrase([]byte(key), []byte(keyPassphrase))
		}
		if err != nil {
			return fmt.Errorf("parse ssh private key failed: %w", err)
		}
		auth = append(auth, gossh.PublicKeys(signer))
	default:
		return fmt.Errorf("unsupported ssh authMethod: %s", authMethod)
	}

	timeout := 5 * time.Second
	sshCfg := &gossh.ClientConfig{
		User:            username,
		Auth:            auth,
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}

	address := net.JoinHostPort(host, port)
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("dial ssh failed: %w", err)
	}

	cc, chans, reqs, err := gossh.NewClientConn(conn, address, sshCfg)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("ssh handshake failed: %w", err)
	}
	client := gossh.NewClient(cc, chans, reqs)
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("create ssh session failed: %w", err)
	}
	defer session.Close()

	if err := session.Run("echo easyssl-ok"); err != nil {
		return fmt.Errorf("run ssh test command failed: %w", err)
	}
	return nil
}

func testTencentCloudAccess(ctx context.Context, config map[string]interface{}) error {
	secretID := stringify(config["secretId"])
	secretKey := stringify(config["secretKey"])
	region := stringify(config["region"])
	sessionToken := stringify(config["sessionToken"])
	if region == "" {
		region = "ap-guangzhou"
	}
	if secretID == "" {
		return fmt.Errorf("tencentcloud.secretId is required")
	}
	if secretKey == "" {
		return fmt.Errorf("tencentcloud.secretKey is required")
	}

	var credential *tccommon.Credential
	if sessionToken != "" {
		credential = tccommon.NewTokenCredential(secretID, secretKey, sessionToken)
	} else {
		credential = tccommon.NewCredential(secretID, secretKey)
	}
	profile := tcprofile.NewClientProfile()
	profile.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	client, err := tcdnspod.NewClient(credential, region, profile)
	if err != nil {
		return fmt.Errorf("create tencentcloud dnspod client failed: %w", err)
	}

	req := tcdnspod.NewDescribeDomainListRequest()
	req.Offset = tccommon.Int64Ptr(0)
	req.Limit = tccommon.Int64Ptr(1)
	if _, err := client.DescribeDomainListWithContext(ctx, req); err != nil {
		if sdkErr, ok := err.(*tcerr.TencentCloudSDKError); ok {
			return fmt.Errorf("tencentcloud access test failed: %s: %s", sdkErr.Code, sdkErr.Message)
		}
		return fmt.Errorf("tencentcloud access test failed: %w", err)
	}
	return nil
}

type qiniuAccessTransport struct {
	http.RoundTripper
	mac *auth.Credentials
}

func newQiniuAccessTransport(mac *auth.Credentials, tr http.RoundTripper) *qiniuAccessTransport {
	if tr == nil {
		tr = http.DefaultTransport
	}
	return &qiniuAccessTransport{RoundTripper: tr, mac: mac}
}

func (t *qiniuAccessTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := t.mac.SignRequestV2(req)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Qiniu "+token)
	return t.RoundTripper.RoundTrip(req)
}

func testQiniuAccess(ctx context.Context, config map[string]interface{}) error {
	accessKey := stringify(config["accessKey"])
	secretKey := stringify(config["secretKey"])
	if accessKey == "" {
		return fmt.Errorf("qiniu.accessKey is required")
	}
	if secretKey == "" {
		return fmt.Errorf("qiniu.secretKey is required")
	}

	mac := auth.New(accessKey, secretKey)
	httpClient := &client.Client{Client: &http.Client{Transport: newQiniuAccessTransport(mac, nil)}}
	resp := struct {
		Code  *int   `json:"code,omitempty"`
		Error string `json:"error,omitempty"`
	}{}
	if err := httpClient.Call(ctx, &resp, http.MethodGet, "https://fusion.qiniuapi.com/sslcert?marker=&limit=1", nil); err != nil {
		return fmt.Errorf("qiniu access test failed: %w", err)
	}
	if resp.Code != nil && *resp.Code != 0 && *resp.Code != 200 {
		return fmt.Errorf("qiniu access test failed: code=%d error=%s", *resp.Code, strings.TrimSpace(resp.Error))
	}
	if resp.Code == nil && strings.TrimSpace(resp.Error) != "" {
		return fmt.Errorf("qiniu access test failed: error=%s", strings.TrimSpace(resp.Error))
	}
	return nil
}
