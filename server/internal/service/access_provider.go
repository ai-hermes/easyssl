package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"easyssl/server/internal/accessprovider"
	"easyssl/server/internal/model"
	"easyssl/server/internal/providercatalog"
	"easyssl/server/internal/repository"
)

const (
	accessProviderAliyun       = accessprovider.ProviderAliyun
	accessProviderTencentCloud = accessprovider.ProviderTencentCloud
	accessProviderQiniu        = accessprovider.ProviderQiniu
	accessProviderSSH          = accessprovider.ProviderSSH
	secretMask                 = providercatalog.SecretMask
)

func normalizeProvider(provider string) string {
	return providercatalog.Normalize(provider)
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
	out.Config = providercatalog.SanitizeAccessConfig(access.Provider, access.Config)
	return out
}

func (s *Service) normalizeAndValidateAccessConfig(in model.Access, current *model.Access) (map[string]interface{}, error) {
	var currentConfig *map[string]interface{}
	if current != nil && normalizeProvider(current.Provider) == normalizeProvider(in.Provider) {
		cfg := current.Config
		currentConfig = &cfg
	}

	cfg, err := providercatalog.ValidateAccessConfig(in.Provider, in.Config, currentConfig)
	if err != nil {
		return nil, err
	}

	if normalizeProvider(in.Provider) == accessProviderSSH {
		authMethod := strings.ToLower(stringify(cfg["authMethod"]))
		password := stringify(cfg["password"])
		key := stringify(cfg["key"])
		if authMethod == "" {
			if key != "" {
				authMethod = "key"
			} else {
				authMethod = "password"
			}
			cfg["authMethod"] = authMethod
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
	}

	return cfg, nil
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
