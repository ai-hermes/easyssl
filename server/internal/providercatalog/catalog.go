package providercatalog

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const SecretMask = "********"

type Field struct {
	Name        string   `json:"name"`
	Label       string   `json:"label"`
	Type        string   `json:"type"`
	Required    bool     `json:"required"`
	Secret      bool     `json:"secret"`
	Default     any      `json:"default,omitempty"`
	Options     []Option `json:"options,omitempty"`
	Placeholder string   `json:"placeholder,omitempty"`
}

type Option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type Definition struct {
	ID               string   `json:"id"`
	Label            string   `json:"label"`
	Kind             string   `json:"kind"`
	AccessProviderID string   `json:"accessProviderId"`
	Capabilities     []string `json:"capabilities"`
	Aliases          []string `json:"aliases,omitempty"`
	AccessFields     []Field  `json:"accessFields"`
	DeployFields     []Field  `json:"deployFields,omitempty"`
}

var definitions []Definition
var byKindID map[string]Definition
var accessByID map[string]Definition
var aliases map[string]string

func init() {
	aliases = map[string]string{
		"alidns":  "aliyun",
		"tencent": "tencentcloud",
		"dnspod":  "tencentcloud",
		"ssh":     "ssh",
	}

	definitions = make([]Definition, 0, len(generatedDefinitions))
	for _, def := range generatedDefinitions {
		def.ID = strings.ToLower(strings.TrimSpace(def.ID))
		def.AccessProviderID = Normalize(def.AccessProviderID)
		applyOverrides(&def)
		definitions = append(definitions, def)
	}
	indexDefinitions()
}

func applyOverrides(def *Definition) {
	switch def.ID {
	case "aliyun":
		def.Aliases = append(def.Aliases, "alidns")
		ensureField(def, Field{Name: "region", Label: "Region", Type: "text", Required: false})
	case "tencentcloud":
		def.Aliases = append(def.Aliases, "tencent", "dnspod")
		ensureField(def, Field{Name: "region", Label: "Region", Type: "text", Required: false, Default: "ap-guangzhou"})
		ensureField(def, Field{Name: "sessionToken", Label: "SessionToken", Type: "password", Required: false, Secret: true})
	case "ssh":
		def.AccessFields = []Field{
			{Name: "host", Label: "Host", Type: "text", Required: true},
			{Name: "port", Label: "Port", Type: "number", Required: false, Default: 22},
			{Name: "username", Label: "Username", Type: "text", Required: false, Default: "root"},
			{Name: "authMethod", Label: "AuthMethod", Type: "select", Required: true, Default: "password", Options: []Option{{Value: "password", Label: "Password"}, {Value: "key", Label: "Private Key"}}},
			{Name: "password", Label: "Password", Type: "password", Required: false, Secret: true},
			{Name: "key", Label: "PrivateKey", Type: "textarea", Required: false, Secret: true},
			{Name: "keyPassphrase", Label: "KeyPassphrase", Type: "password", Required: false, Secret: true},
		}
	}

	if def.Kind == "deploy" {
		switch def.ID {
		case "aliyun-cas":
			def.DeployFields = []Field{{Name: "region", Label: "Region", Type: "text", Required: false}, {Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false}}
		case "qiniu":
			def.DeployFields = []Field{{Name: "certName", Label: "CertName", Type: "text", Required: false}, {Name: "commonName", Label: "CommonName", Type: "text", Required: false}}
		case "ssh":
			def.DeployFields = []Field{{Name: "certPath", Label: "CertPath", Type: "text", Required: false, Default: "/etc/nginx/ssl/fullchain.pem"}, {Name: "keyPath", Label: "KeyPath", Type: "text", Required: false, Default: "/etc/nginx/ssl/privkey.pem"}, {Name: "preCommand", Label: "PreCommand", Type: "textarea", Required: false}, {Name: "postCommand", Label: "PostCommand", Type: "textarea", Required: false}}
		default:
			def.DeployFields = []Field{{Name: "extendedConfig", Label: "ExtendedConfigJSON", Type: "textarea", Required: false, Placeholder: "Optional JSON object for provider-specific deployment settings"}}
		}
	}
}

func ensureField(def *Definition, field Field) {
	for i := range def.AccessFields {
		if def.AccessFields[i].Name == field.Name {
			return
		}
	}
	def.AccessFields = append(def.AccessFields, field)
}

func indexDefinitions() {
	byKindID = map[string]Definition{}
	accessByID = map[string]Definition{}
	for _, def := range definitions {
		byKindID[def.Kind+":"+def.ID] = def
		if def.Kind == "access" {
			accessByID[def.ID] = def
			for _, alias := range def.Aliases {
				aliases[Normalize(alias)] = def.ID
			}
		}
	}
}

func List() []Definition {
	out := append([]Definition(nil), definitions...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Kind != out[j].Kind {
			return out[i].Kind < out[j].Kind
		}
		return out[i].ID < out[j].ID
	})
	return out
}

func ListByKind(kind string) []Definition {
	kind = strings.ToLower(strings.TrimSpace(kind))
	out := []Definition{}
	for _, def := range definitions {
		if kind == "" || def.Kind == kind {
			out = append(out, def)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func Normalize(provider string) string {
	p := strings.ToLower(strings.TrimSpace(provider))
	if v, ok := aliases[p]; ok {
		return v
	}
	return p
}

func AccessDefinition(provider string) (Definition, bool) {
	provider = Normalize(provider)
	def, ok := accessByID[provider]
	return def, ok
}

func OperationDefinition(kind, provider string) (Definition, bool) {
	provider = Normalize(provider)
	def, ok := byKindID[strings.ToLower(strings.TrimSpace(kind))+":"+provider]
	return def, ok
}

func ValidateAccessConfig(provider string, config map[string]interface{}, current *map[string]interface{}) (map[string]interface{}, error) {
	def, ok := AccessDefinition(provider)
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
	return validateFields(def.ID, def.AccessFields, config, current)
}

func SanitizeAccessConfig(provider string, config map[string]interface{}) map[string]interface{} {
	out := clone(config)
	def, ok := AccessDefinition(provider)
	if !ok {
		return out
	}
	for _, field := range def.AccessFields {
		if field.Secret {
			if _, exists := out[field.Name]; exists {
				out[field.Name] = SecretMask
			}
		}
	}
	return out
}

func validateFields(provider string, fields []Field, in map[string]interface{}, current *map[string]interface{}) (map[string]interface{}, error) {
	cfg := clone(in)
	out := map[string]interface{}{}
	for _, field := range fields {
		value, exists := cfg[field.Name]
		if field.Secret {
			text := asString(value)
			if (!exists || text == "" || text == SecretMask) && current != nil {
				if old, ok := (*current)[field.Name]; ok {
					value = old
					exists = true
				}
			}
			if asString(value) == SecretMask {
				value = ""
				exists = false
			}
		}

		if (!exists || isEmpty(value)) && field.Default != nil {
			value = field.Default
			exists = true
		}
		if field.Required && (!exists || isEmpty(value)) {
			return nil, fmt.Errorf("%s.%s is required", provider, field.Name)
		}
		if !exists || isEmpty(value) {
			continue
		}

		switch field.Type {
		case "number":
			n, err := toInt(value)
			if err != nil {
				return nil, fmt.Errorf("%s.%s must be a number", provider, field.Name)
			}
			out[field.Name] = n
		case "checkbox":
			out[field.Name] = toBool(value)
		default:
			out[field.Name] = asString(value)
		}
	}
	return out, nil
}

func clone(in map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{}
	for k, v := range in {
		out[k] = v
	}
	return out
}

func isEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t) == ""
	default:
		return false
	}
}

func asString(v interface{}) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%v", v))
}

func toInt(v interface{}) (int, error) {
	switch t := v.(type) {
	case int:
		return t, nil
	case int32:
		return int(t), nil
	case int64:
		return int(t), nil
	case float64:
		return int(t), nil
	case float32:
		return int(t), nil
	default:
		return strconv.Atoi(asString(v))
	}
}

func toBool(v interface{}) bool {
	switch t := v.(type) {
	case bool:
		return t
	case string:
		s := strings.ToLower(strings.TrimSpace(t))
		return s == "true" || s == "1" || s == "yes" || s == "on"
	default:
		return false
	}
}
