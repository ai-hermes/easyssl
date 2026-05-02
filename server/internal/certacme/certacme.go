package certacme

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/alidns"
	"github.com/go-acme/lego/v4/providers/dns/tencentcloud"
	legoReg "github.com/go-acme/lego/v4/registration"
	"golang.org/x/net/publicsuffix"
)

const (
	CALEProd    = "letsencrypt"
	CALEStaging = "letsencrypt-staging"
)

type ObtainRequest struct {
	Email                 string
	Domains               []string
	CAProvider            string
	DNSProvider           string
	AliyunAccessKeyID     string
	AliyunAccessKeySecret string
	TencentSecretID       string
	TencentSecretKey      string
	TencentRegion         string
	TencentSessionToken   string
	DNSPropagationTimeout int
	DNSTTL                int
	KeyAlgorithm          string
}

type ObtainResponse struct {
	Certificate string
	PrivateKey  string
}

type account struct {
	email string
	key   crypto.PrivateKey
	reg   *legoReg.Resource
}

func (a *account) GetEmail() string {
	return a.email
}

func (a *account) GetRegistration() *legoReg.Resource {
	return a.reg
}

func (a *account) GetPrivateKey() crypto.PrivateKey {
	return a.key
}

func Obtain(req ObtainRequest) (*ObtainResponse, error) {
	if len(req.Domains) == 0 {
		return nil, fmt.Errorf("domains is required")
	}

	keyType, err := readKeyType(req.KeyAlgorithm)
	if err != nil {
		return nil, err
	}

	privateKey, err := certcrypto.GeneratePrivateKey(keyType)
	if err != nil {
		return nil, fmt.Errorf("generate account key failed: %w", err)
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		autoEmail, err := buildAutoContactEmail(req.Domains)
		if err != nil {
			return nil, fmt.Errorf("contactEmail is required and could not infer from domains: %w", err)
		}
		email = autoEmail
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, fmt.Errorf("invalid contactEmail: %w", err)
	}

	user := &account{email: email, key: privateKey}
	legoCfg := lego.NewConfig(user)
	legoCfg.Certificate.KeyType = keyType
	legoCfg.CADirURL = readCADir(req.CAProvider)

	client, err := lego.NewClient(legoCfg)
	if err != nil {
		return nil, fmt.Errorf("create acme client failed: %w", err)
	}

	reg, err := client.Registration.Register(legoReg.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, fmt.Errorf("register acme account failed: %w", err)
	}
	user.reg = reg

	provider, err := newDNSProvider(req)
	if err != nil {
		return nil, err
	}

	if err := client.Challenge.SetDNS01Provider(provider, dns01.DisableAuthoritativeNssPropagationRequirement()); err != nil {
		return nil, fmt.Errorf("set dns challenge provider failed: %w", err)
	}

	obtainReq := certificate.ObtainRequest{Domains: req.Domains, Bundle: true}
	resp, err := client.Certificate.Obtain(obtainReq)
	if err != nil {
		return nil, fmt.Errorf("obtain certificate failed: %w", err)
	}

	return &ObtainResponse{
		Certificate: strings.TrimSpace(string(resp.Certificate)),
		PrivateKey:  strings.TrimSpace(string(resp.PrivateKey)),
	}, nil
}

func readCADir(provider string) string {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "", CALEProd:
		return lego.LEDirectoryProduction
	case CALEStaging:
		return lego.LEDirectoryStaging
	default:
		return lego.LEDirectoryProduction
	}
}

func readKeyType(name string) (certcrypto.KeyType, error) {
	switch strings.ToUpper(strings.TrimSpace(name)) {
	case "", "RSA2048":
		return certcrypto.RSA2048, nil
	case "RSA4096":
		return certcrypto.RSA4096, nil
	case "EC256":
		return certcrypto.EC256, nil
	case "EC384":
		return certcrypto.EC384, nil
	default:
		return "", fmt.Errorf("unsupported keyAlgorithm: %s", name)
	}
}

func buildAutoContactEmail(domains []string) (string, error) {
	for _, item := range domains {
		d := normalizeDomain(item)
		if d == "" {
			continue
		}

		root, err := publicsuffix.EffectiveTLDPlusOne(d)
		if err != nil || root == "" {
			continue
		}

		buf := make([]byte, 4)
		if _, err := rand.Read(buf); err != nil {
			return "", err
		}
		localPart := "easyssl-" + hex.EncodeToString(buf)
		email := localPart + "@" + root
		if _, err := mail.ParseAddress(email); err != nil {
			continue
		}
		return email, nil
	}
	return "", fmt.Errorf("no valid domain found in %v", domains)
}

func normalizeDomain(domain string) string {
	d := strings.TrimSpace(domain)
	d = strings.TrimPrefix(d, "*.")
	d = strings.TrimPrefix(d, "http://")
	d = strings.TrimPrefix(d, "https://")
	if d == "" {
		return ""
	}
	if u, err := url.Parse("https://" + d); err == nil && u.Hostname() != "" {
		d = u.Hostname()
	}
	return strings.ToLower(strings.TrimSpace(d))
}

func newDNSProvider(req ObtainRequest) (challenge.Provider, error) {
	switch normalizeDNSProvider(req.DNSProvider) {
	case "", "aliyun":
		if req.AliyunAccessKeyID == "" {
			return nil, fmt.Errorf("aliyun.accessKeyId is required")
		}
		if req.AliyunAccessKeySecret == "" {
			return nil, fmt.Errorf("aliyun.accessKeySecret is required")
		}
		dnsCfg := alidns.NewDefaultConfig()
		dnsCfg.APIKey = req.AliyunAccessKeyID
		dnsCfg.SecretKey = req.AliyunAccessKeySecret
		if req.DNSPropagationTimeout > 0 {
			dnsCfg.PropagationTimeout = time.Duration(req.DNSPropagationTimeout) * time.Second
		}
		if req.DNSTTL > 0 {
			dnsCfg.TTL = req.DNSTTL
		}
		provider, err := alidns.NewDNSProviderConfig(dnsCfg)
		if err != nil {
			return nil, fmt.Errorf("create aliyun dns provider failed: %w", err)
		}
		return provider, nil

	case "tencentcloud":
		if req.TencentSecretID == "" {
			return nil, fmt.Errorf("tencentcloud.secretId is required")
		}
		if req.TencentSecretKey == "" {
			return nil, fmt.Errorf("tencentcloud.secretKey is required")
		}
		dnsCfg := tencentcloud.NewDefaultConfig()
		dnsCfg.SecretID = req.TencentSecretID
		dnsCfg.SecretKey = req.TencentSecretKey
		dnsCfg.Region = req.TencentRegion
		dnsCfg.SessionToken = req.TencentSessionToken
		if req.DNSPropagationTimeout > 0 {
			dnsCfg.PropagationTimeout = time.Duration(req.DNSPropagationTimeout) * time.Second
		}
		if req.DNSTTL > 0 {
			dnsCfg.TTL = req.DNSTTL
		}
		provider, err := tencentcloud.NewDNSProviderConfig(dnsCfg)
		if err != nil {
			return nil, fmt.Errorf("create tencentcloud dns provider failed: %w", err)
		}
		return provider, nil

	default:
		return nil, fmt.Errorf("unsupported dns provider: %s", req.DNSProvider)
	}
}

func normalizeDNSProvider(provider string) string {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "", "aliyun", "alidns":
		return "aliyun"
	case "tencent", "tencentcloud", "dnspod":
		return "tencentcloud"
	default:
		return strings.ToLower(strings.TrimSpace(provider))
	}
}
