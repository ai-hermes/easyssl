package providercatalog

import "testing"

func TestListIncludesCertimateProviderKinds(t *testing.T) {
	if got := len(ListByKind("dns")); got < 66 {
		t.Fatalf("dns providers = %d, want at least 66", got)
	}
	if got := len(ListByKind("deploy")); got < 128 {
		t.Fatalf("deploy providers = %d, want at least 128", got)
	}
	if got := len(ListByKind("access")); got < 100 {
		t.Fatalf("access providers = %d, want at least 100", got)
	}
}

func TestValidateAccessConfigPreservesSecrets(t *testing.T) {
	current := map[string]interface{}{"accessKeyId": "id", "accessKeySecret": "old-secret"}
	cfg, err := ValidateAccessConfig("alidns", map[string]interface{}{"accessKeyId": "id", "accessKeySecret": SecretMask}, &current)
	if err != nil {
		t.Fatal(err)
	}
	if cfg["accessKeySecret"] != "old-secret" {
		t.Fatalf("secret was not preserved: %#v", cfg["accessKeySecret"])
	}
}

func TestSanitizeAccessConfigMasksSecretFields(t *testing.T) {
	cfg := SanitizeAccessConfig("qiniu", map[string]interface{}{"accessKey": "ak", "secretKey": "sk"})
	if cfg["accessKey"] != "ak" {
		t.Fatalf("access key should not be masked")
	}
	if cfg["secretKey"] != SecretMask {
		t.Fatalf("secret key should be masked, got %#v", cfg["secretKey"])
	}
}
