package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"easyssl/server/internal/accessprovider"
	"easyssl/server/internal/model"
	"easyssl/server/internal/repository"
)

const (
	accessProviderAliyun       = accessprovider.ProviderAliyun
	accessProviderTencentCloud = accessprovider.ProviderTencentCloud
	accessProviderQiniu        = accessprovider.ProviderQiniu
	accessProviderSSH          = accessprovider.ProviderSSH
	secretMask                 = "********"
)

func normalizeProvider(provider string) string {
	return strings.ToLower(strings.TrimSpace(provider))
}

func stringify(v interface{}) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%v", v))
}

func cloneConfig(in map[string]interface{}) map[string]interface{} {
	if in == nil {
		return map[string]interface{}{}
	}
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func parseInt(v interface{}, fallback int) int {
	s := stringify(v)
	if s == "" {
		return fallback
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return n
}

func sanitizeAccess(access model.Access) model.Access {
	out := access
	out.Config = cloneConfig(access.Config)
	switch out.Provider {
	case accessProviderAliyun:
		if _, ok := out.Config["accessKeySecret"]; ok {
			out.Config["accessKeySecret"] = secretMask
		}
	case accessProviderTencentCloud:
		if _, ok := out.Config["secretKey"]; ok {
			out.Config["secretKey"] = secretMask
		}
		if _, ok := out.Config["sessionToken"]; ok {
			out.Config["sessionToken"] = secretMask
		}
	case accessProviderQiniu:
		if _, ok := out.Config["secretKey"]; ok {
			out.Config["secretKey"] = secretMask
		}
	case accessProviderSSH:
		if _, ok := out.Config["password"]; ok {
			out.Config["password"] = secretMask
		}
		if _, ok := out.Config["key"]; ok {
			out.Config["key"] = secretMask
		}
		if _, ok := out.Config["keyPassphrase"]; ok {
			out.Config["keyPassphrase"] = secretMask
		}
	}
	return out
}

func (s *Service) normalizeAndValidateAccessConfig(in model.Access, current *model.Access) (map[string]interface{}, error) {
	switch in.Provider {
	case accessProviderAliyun:
		cfg := cloneConfig(in.Config)
		accessKeyID := stringify(cfg["accessKeyId"])
		accessKeySecret := stringify(cfg["accessKeySecret"])
		resourceGroupID := stringify(cfg["resourceGroupId"])
		region := stringify(cfg["region"])

		if accessKeyID == "" {
			return nil, fmt.Errorf("aliyun.accessKeyId is required")
		}

		if accessKeySecret == "" || accessKeySecret == secretMask {
			if current != nil && current.Provider == accessProviderAliyun {
				accessKeySecret = stringify(current.Config["accessKeySecret"])
			}
		}
		if accessKeySecret == "" {
			return nil, fmt.Errorf("aliyun.accessKeySecret is required")
		}

		out := map[string]interface{}{
			"accessKeyId":     accessKeyID,
			"accessKeySecret": accessKeySecret,
		}
		if resourceGroupID != "" {
			out["resourceGroupId"] = resourceGroupID
		}
		if region != "" {
			out["region"] = region
		}
		return out, nil

	case accessProviderTencentCloud:
		cfg := cloneConfig(in.Config)
		secretID := stringify(cfg["secretId"])
		secretKey := stringify(cfg["secretKey"])
		region := stringify(cfg["region"])
		sessionToken := stringify(cfg["sessionToken"])

		if secretID == "" {
			return nil, fmt.Errorf("tencentcloud.secretId is required")
		}

		if secretKey == "" || secretKey == secretMask {
			if current != nil && current.Provider == accessProviderTencentCloud {
				secretKey = stringify(current.Config["secretKey"])
			}
		}
		if secretKey == "" {
			return nil, fmt.Errorf("tencentcloud.secretKey is required")
		}

		if sessionToken == secretMask && current != nil && current.Provider == accessProviderTencentCloud {
			sessionToken = stringify(current.Config["sessionToken"])
		}

		out := map[string]interface{}{
			"secretId":  secretID,
			"secretKey": secretKey,
		}
		if region != "" {
			out["region"] = region
		}
		if sessionToken != "" {
			out["sessionToken"] = sessionToken
		}
		return out, nil

	case accessProviderQiniu:
		cfg := cloneConfig(in.Config)
		accessKey := stringify(cfg["accessKey"])
		secretKey := stringify(cfg["secretKey"])

		if accessKey == "" {
			return nil, fmt.Errorf("qiniu.accessKey is required")
		}
		if secretKey == "" || secretKey == secretMask {
			if current != nil && current.Provider == accessProviderQiniu {
				secretKey = stringify(current.Config["secretKey"])
			}
		}
		if secretKey == "" {
			return nil, fmt.Errorf("qiniu.secretKey is required")
		}

		return map[string]interface{}{
			"accessKey": accessKey,
			"secretKey": secretKey,
		}, nil

	case accessProviderSSH:
		cfg := cloneConfig(in.Config)
		host := stringify(cfg["host"])
		port := parseInt(cfg["port"], 22)
		username := stringify(cfg["username"])
		authMethod := strings.ToLower(stringify(cfg["authMethod"]))
		password := stringify(cfg["password"])
		key := stringify(cfg["key"])
		keyPassphrase := stringify(cfg["keyPassphrase"])

		if host == "" {
			return nil, fmt.Errorf("ssh.host is required")
		}
		if username == "" {
			username = "root"
		}
		if authMethod == "" {
			if key != "" {
				authMethod = "key"
			} else {
				authMethod = "password"
			}
		}

		if current != nil && current.Provider == accessProviderSSH {
			if password == "" || password == secretMask {
				password = stringify(current.Config["password"])
			}
			if key == "" || key == secretMask {
				key = stringify(current.Config["key"])
			}
			if keyPassphrase == "" || keyPassphrase == secretMask {
				keyPassphrase = stringify(current.Config["keyPassphrase"])
			}
		}

		if authMethod == "password" && password == "" {
			return nil, fmt.Errorf("ssh.password is required")
		}
		if authMethod == "key" && key == "" {
			return nil, fmt.Errorf("ssh.key is required")
		}
		if authMethod != "password" && authMethod != "key" {
			return nil, fmt.Errorf("ssh.authMethod must be password or key")
		}

		out := map[string]interface{}{
			"host":       host,
			"port":       port,
			"username":   username,
			"authMethod": authMethod,
		}
		if password != "" {
			out["password"] = password
		}
		if key != "" {
			out["key"] = key
		}
		if keyPassphrase != "" {
			out["keyPassphrase"] = keyPassphrase
		}
		return out, nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", in.Provider)
	}
}

func (s *Service) prepareAccess(ctx context.Context, auth model.AuthContext, in model.Access) (model.Access, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Provider = normalizeProvider(in.Provider)

	if in.Name == "" {
		return in, fmt.Errorf("name is required")
	}
	if in.Provider == "" {
		return in, fmt.Errorf("provider is required")
	}

	var current *model.Access
	if in.ID != "" {
		old, err := s.repo.GetAccessByIDForUser(ctx, in.ID, auth.UserID, auth.Role)
		if err != nil {
			if err == repository.ErrNotFound {
				return in, fmt.Errorf("access not found")
			}
			return in, err
		}
		current = old
	}

	cfg, err := s.normalizeAndValidateAccessConfig(in, current)
	if err != nil {
		return in, err
	}
	in.Config = cfg
	return in, nil
}
