// Code generated from certimate provider metadata by /tmp/gen_provider_catalog.py; DO NOT EDIT.
package providercatalog

var generatedDefinitions = []Definition{
	{ID: "1panel", Label: "1Panel", Kind: "access", AccessProviderID: "1panel", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiVersion", Label: "ApiVersion", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "35cn", Label: "35Cn", Kind: "access", AccessProviderID: "35cn", Capabilities: nil, AccessFields: []Field{
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "apiPassword", Label: "ApiPassword", Type: "password", Required: true, Secret: true},
	}},
	{ID: "51dnscom", Label: "51Dnscom", Kind: "access", AccessProviderID: "51dnscom", Capabilities: nil, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "apiSecret", Label: "ApiSecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "acmeca", Label: "Acmeca", Kind: "access", AccessProviderID: "acmeca", Capabilities: nil, AccessFields: []Field{
		{Name: "endpoint", Label: "Endpoint", Type: "text", Required: true, Secret: false},
	}},
	{ID: "acmedns", Label: "Acmedns", Kind: "access", AccessProviderID: "acmedns", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "credentials", Label: "Credentials", Type: "password", Required: true, Secret: true},
	}},
	{ID: "acmehttpreq", Label: "Acmehttpreq", Kind: "access", AccessProviderID: "acmehttpreq", Capabilities: nil, AccessFields: []Field{
		{Name: "endpoint", Label: "Endpoint", Type: "text", Required: true, Secret: false},
		{Name: "mode", Label: "Mode", Type: "text", Required: false, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: false, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: false, Secret: true},
	}},
	{ID: "actalisssl", Label: "Actalisssl", Kind: "access", AccessProviderID: "actalisssl", Capabilities: nil, AccessFields: []Field{}},
	{ID: "akamai", Label: "Akamai", Kind: "access", AccessProviderID: "akamai", Capabilities: nil, AccessFields: []Field{
		{Name: "host", Label: "Host", Type: "text", Required: true, Secret: false},
		{Name: "clientToken", Label: "ClientToken", Type: "password", Required: true, Secret: true},
		{Name: "clientSecret", Label: "ClientSecret", Type: "password", Required: true, Secret: true},
		{Name: "accessToken", Label: "AccessToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "aliyun", Label: "Aliyun", Kind: "access", AccessProviderID: "aliyun", Capabilities: nil, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "apisix", Label: "Apisix", Kind: "access", AccessProviderID: "apisix", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "arvancloud", Label: "Arvancloud", Kind: "access", AccessProviderID: "arvancloud", Capabilities: nil, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "aws", Label: "AWS", Kind: "access", AccessProviderID: "aws", Capabilities: nil, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "azure", Label: "Azure", Kind: "access", AccessProviderID: "azure", Capabilities: nil, AccessFields: []Field{
		{Name: "tenantId", Label: "TenantId", Type: "text", Required: true, Secret: false},
		{Name: "clientId", Label: "ClientId", Type: "text", Required: true, Secret: false},
		{Name: "clientSecret", Label: "ClientSecret", Type: "password", Required: true, Secret: true},
		{Name: "subscriptionId", Label: "SubscriptionId", Type: "text", Required: false, Secret: false},
		{Name: "resourceGroupName", Label: "ResourceGroupName", Type: "text", Required: false, Secret: false},
		{Name: "cloudName", Label: "CloudName", Type: "text", Required: false, Secret: false},
	}},
	{ID: "baiducloud", Label: "Baiducloud", Kind: "access", AccessProviderID: "baiducloud", Capabilities: nil, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "baishan", Label: "Baishan", Kind: "access", AccessProviderID: "baishan", Capabilities: nil, AccessFields: []Field{
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "baotapanel", Label: "Baotapanel", Kind: "access", AccessProviderID: "baotapanel", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "baotapanelgo", Label: "Baotapanelgo", Kind: "access", AccessProviderID: "baotapanelgo", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "baotawaf", Label: "Baotawaf", Kind: "access", AccessProviderID: "baotawaf", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "bookmyname", Label: "Bookmyname", Kind: "access", AccessProviderID: "bookmyname", Capabilities: nil, AccessFields: []Field{
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: true, Secret: true},
	}},
	{ID: "bunny", Label: "Bunny", Kind: "access", AccessProviderID: "bunny", Capabilities: nil, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "byteplus", Label: "Byteplus", Kind: "access", AccessProviderID: "byteplus", Capabilities: nil, AccessFields: []Field{
		{Name: "accessKey", Label: "AccessKey", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "cachefly", Label: "Cachefly", Kind: "access", AccessProviderID: "cachefly", Capabilities: nil, AccessFields: []Field{
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "cdnfly", Label: "Cdnfly", Kind: "access", AccessProviderID: "cdnfly", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "apiSecret", Label: "ApiSecret", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "cloudflare", Label: "Cloudflare", Kind: "access", AccessProviderID: "cloudflare", Capabilities: nil, AccessFields: []Field{
		{Name: "dnsApiToken", Label: "DnsApiToken", Type: "password", Required: true, Secret: true},
		{Name: "zoneApiToken", Label: "ZoneApiToken", Type: "password", Required: false, Secret: true},
	}},
	{ID: "cloudns", Label: "Cloudns", Kind: "access", AccessProviderID: "cloudns", Capabilities: nil, AccessFields: []Field{
		{Name: "authId", Label: "AuthId", Type: "text", Required: true, Secret: false},
		{Name: "authPassword", Label: "AuthPassword", Type: "password", Required: true, Secret: true},
	}},
	{ID: "cmcccloud", Label: "Cmcccloud", Kind: "access", AccessProviderID: "cmcccloud", Capabilities: nil, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "constellix", Label: "Constellix", Kind: "access", AccessProviderID: "constellix", Capabilities: nil, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "cpanel", Label: "Cpanel", Kind: "access", AccessProviderID: "cpanel", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "ctcccloud", Label: "Ctcccloud", Kind: "access", AccessProviderID: "ctcccloud", Capabilities: nil, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "cucccloud", Label: "Cucccloud", Kind: "access", AccessProviderID: "cucccloud", Capabilities: nil, AccessFields: []Field{}},
	{ID: "desec", Label: "Desec", Kind: "access", AccessProviderID: "desec", Capabilities: nil, AccessFields: []Field{
		{Name: "token", Label: "Token", Type: "password", Required: true, Secret: true},
	}},
	{ID: "digicert", Label: "Digicert", Kind: "access", AccessProviderID: "digicert", Capabilities: nil, AccessFields: []Field{}},
	{ID: "digitalocean", Label: "Digitalocean", Kind: "access", AccessProviderID: "digitalocean", Capabilities: nil, AccessFields: []Field{
		{Name: "accessToken", Label: "AccessToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "dingtalkbot", Label: "Dingtalkbot", Kind: "access", AccessProviderID: "dingtalkbot", Capabilities: nil, AccessFields: []Field{
		{Name: "webhookUrl", Label: "WebhookUrl", Type: "text", Required: true, Secret: false},
		{Name: "secret", Label: "Secret", Type: "password", Required: false, Secret: true},
		{Name: "customPayload", Label: "CustomPayload", Type: "textarea", Required: false, Secret: false},
	}},
	{ID: "discordbot", Label: "Discordbot", Kind: "access", AccessProviderID: "discordbot", Capabilities: nil, AccessFields: []Field{
		{Name: "botToken", Label: "BotToken", Type: "password", Required: true, Secret: true},
		{Name: "channelId", Label: "ChannelId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "dnsexit", Label: "Dnsexit", Kind: "access", AccessProviderID: "dnsexit", Capabilities: nil, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "dnsla", Label: "Dnsla", Kind: "access", AccessProviderID: "dnsla", Capabilities: nil, AccessFields: []Field{
		{Name: "apiId", Label: "ApiId", Type: "text", Required: true, Secret: false},
		{Name: "apiSecret", Label: "ApiSecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "dnsmadeeasy", Label: "Dnsmadeeasy", Kind: "access", AccessProviderID: "dnsmadeeasy", Capabilities: nil, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "apiSecret", Label: "ApiSecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "dogecloud", Label: "Dogecloud", Kind: "access", AccessProviderID: "dogecloud", Capabilities: nil, AccessFields: []Field{
		{Name: "accessKey", Label: "AccessKey", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "dokploy", Label: "Dokploy", Kind: "access", AccessProviderID: "dokploy", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "duckdns", Label: "Duckdns", Kind: "access", AccessProviderID: "duckdns", Capabilities: nil, AccessFields: []Field{
		{Name: "token", Label: "Token", Type: "password", Required: true, Secret: true},
	}},
	{ID: "dynu", Label: "Dynu", Kind: "access", AccessProviderID: "dynu", Capabilities: nil, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "dynv6", Label: "Dynv6", Kind: "access", AccessProviderID: "dynv6", Capabilities: nil, AccessFields: []Field{
		{Name: "httpToken", Label: "HttpToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "email", Label: "Email", Kind: "access", AccessProviderID: "email", Capabilities: nil, AccessFields: []Field{
		{Name: "smtpHost", Label: "SmtpHost", Type: "text", Required: true, Secret: false},
		{Name: "smtpPort", Label: "SmtpPort", Type: "number", Required: true, Secret: false},
		{Name: "smtpTls", Label: "SmtpTls", Type: "checkbox", Required: true, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: true, Secret: true},
		{Name: "senderAddress", Label: "SenderAddress", Type: "text", Required: true, Secret: false},
		{Name: "senderName", Label: "SenderName", Type: "text", Required: true, Secret: false},
		{Name: "receiverAddress", Label: "ReceiverAddress", Type: "text", Required: false, Secret: false},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "fastly", Label: "Fastly", Kind: "access", AccessProviderID: "fastly", Capabilities: nil, AccessFields: []Field{}},
	{ID: "flexcdn", Label: "Flexcdn", Kind: "access", AccessProviderID: "flexcdn", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiRole", Label: "ApiRole", Type: "text", Required: true, Secret: false},
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKey", Label: "AccessKey", Type: "text", Required: true, Secret: false},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "flyio", Label: "Flyio", Kind: "access", AccessProviderID: "flyio", Capabilities: nil, AccessFields: []Field{
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ftp", Label: "FTP", Kind: "access", AccessProviderID: "ftp", Capabilities: nil, AccessFields: []Field{
		{Name: "host", Label: "Host", Type: "text", Required: true, Secret: false},
		{Name: "port", Label: "Port", Type: "number", Required: true, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: false, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: false, Secret: true},
	}},
	{ID: "gandinet", Label: "Gandinet", Kind: "access", AccessProviderID: "gandinet", Capabilities: nil, AccessFields: []Field{
		{Name: "personalAccessToken", Label: "PersonalAccessToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "gcore", Label: "Gcore", Kind: "access", AccessProviderID: "gcore", Capabilities: nil, AccessFields: []Field{
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "globalsignatlas", Label: "Globalsignatlas", Kind: "access", AccessProviderID: "globalsignatlas", Capabilities: nil, AccessFields: []Field{}},
	{ID: "gname", Label: "Gname", Kind: "access", AccessProviderID: "gname", Capabilities: nil, AccessFields: []Field{
		{Name: "appId", Label: "AppId", Type: "text", Required: true, Secret: false},
		{Name: "appKey", Label: "AppKey", Type: "text", Required: true, Secret: false},
	}},
	{ID: "godaddy", Label: "Godaddy", Kind: "access", AccessProviderID: "godaddy", Capabilities: nil, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "apiSecret", Label: "ApiSecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "goedge", Label: "Goedge", Kind: "access", AccessProviderID: "goedge", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiRole", Label: "ApiRole", Type: "text", Required: true, Secret: false},
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKey", Label: "AccessKey", Type: "text", Required: true, Secret: false},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "googletrustservices", Label: "Googletrustservices", Kind: "access", AccessProviderID: "googletrustservices", Capabilities: nil, AccessFields: []Field{}},
	{ID: "hetzner", Label: "Hetzner", Kind: "access", AccessProviderID: "hetzner", Capabilities: nil, AccessFields: []Field{
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "hostingde", Label: "Hostingde", Kind: "access", AccessProviderID: "hostingde", Capabilities: nil, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "hostinger", Label: "Hostinger", Kind: "access", AccessProviderID: "hostinger", Capabilities: nil, AccessFields: []Field{
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "huaweicloud", Label: "Huaweicloud", Kind: "access", AccessProviderID: "huaweicloud", Capabilities: nil, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
		{Name: "enterpriseProjectId", Label: "EnterpriseProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "infomaniak", Label: "Infomaniak", Kind: "access", AccessProviderID: "infomaniak", Capabilities: nil, AccessFields: []Field{
		{Name: "accessToken", Label: "AccessToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ionos", Label: "Ionos", Kind: "access", AccessProviderID: "ionos", Capabilities: nil, AccessFields: []Field{
		{Name: "apiKeyPublicPrefix", Label: "ApiKeyPublicPrefix", Type: "password", Required: true, Secret: true},
		{Name: "apiKeySecret", Label: "ApiKeySecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "jdcloud", Label: "Jdcloud", Kind: "access", AccessProviderID: "jdcloud", Capabilities: nil, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "kong", Label: "Kong", Kind: "access", AccessProviderID: "kong", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: false, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "ksyun", Label: "Ksyun", Kind: "access", AccessProviderID: "ksyun", Capabilities: nil, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "k8s", Label: "K8S", Kind: "access", AccessProviderID: "k8s", Capabilities: nil, AccessFields: []Field{
		{Name: "kubeConfig", Label: "KubeConfig", Type: "text", Required: false, Secret: false},
	}},
	{ID: "larkbot", Label: "Larkbot", Kind: "access", AccessProviderID: "larkbot", Capabilities: nil, AccessFields: []Field{
		{Name: "webhookUrl", Label: "WebhookUrl", Type: "text", Required: true, Secret: false},
		{Name: "secret", Label: "Secret", Type: "password", Required: false, Secret: true},
		{Name: "customPayload", Label: "CustomPayload", Type: "textarea", Required: false, Secret: false},
	}},
	{ID: "lecdn", Label: "Lecdn", Kind: "access", AccessProviderID: "lecdn", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiVersion", Label: "ApiVersion", Type: "text", Required: true, Secret: false},
		{Name: "apiRole", Label: "ApiRole", Type: "text", Required: true, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "letsencrypt", Label: "Letsencrypt", Kind: "access", AccessProviderID: "letsencrypt", Capabilities: nil, AccessFields: []Field{}},
	{ID: "letsencryptstaging", Label: "Letsencryptstaging", Kind: "access", AccessProviderID: "letsencryptstaging", Capabilities: nil, AccessFields: []Field{}},
	{ID: "linode", Label: "Linode", Kind: "access", AccessProviderID: "linode", Capabilities: nil, AccessFields: []Field{
		{Name: "accessToken", Label: "AccessToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "litessl", Label: "Litessl", Kind: "access", AccessProviderID: "litessl", Capabilities: nil, AccessFields: []Field{}},
	{ID: "local", Label: "Local", Kind: "access", AccessProviderID: "local", Capabilities: nil, AccessFields: []Field{}},
	{ID: "mattermost", Label: "Mattermost", Kind: "access", AccessProviderID: "mattermost", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: true, Secret: true},
		{Name: "channelId", Label: "ChannelId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "mohua", Label: "Mohua", Kind: "access", AccessProviderID: "mohua", Capabilities: nil, AccessFields: []Field{
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "apiPassword", Label: "ApiPassword", Type: "password", Required: true, Secret: true},
	}},
	{ID: "namecheap", Label: "Namecheap", Kind: "access", AccessProviderID: "namecheap", Capabilities: nil, AccessFields: []Field{
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "namedotcom", Label: "Namedotcom", Kind: "access", AccessProviderID: "namedotcom", Capabilities: nil, AccessFields: []Field{
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "namesilo", Label: "Namesilo", Kind: "access", AccessProviderID: "namesilo", Capabilities: nil, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "netcup", Label: "Netcup", Kind: "access", AccessProviderID: "netcup", Capabilities: nil, AccessFields: []Field{
		{Name: "customerNumber", Label: "CustomerNumber", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "apiPassword", Label: "ApiPassword", Type: "password", Required: true, Secret: true},
	}},
	{ID: "netlify", Label: "Netlify", Kind: "access", AccessProviderID: "netlify", Capabilities: nil, AccessFields: []Field{
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "nginxproxymanager", Label: "Nginxproxymanager", Kind: "access", AccessProviderID: "nginxproxymanager", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "authMethod", Label: "AuthMethod", Type: "text", Required: true, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: false, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: false, Secret: true},
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: false, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "ns1", Label: "Ns1", Kind: "access", AccessProviderID: "ns1", Capabilities: nil, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ovhcloud", Label: "Ovhcloud", Kind: "access", AccessProviderID: "ovhcloud", Capabilities: nil, AccessFields: []Field{
		{Name: "endpoint", Label: "Endpoint", Type: "text", Required: true, Secret: false},
		{Name: "authMethod", Label: "AuthMethod", Type: "text", Required: true, Secret: false},
		{Name: "applicationKey", Label: "ApplicationKey", Type: "text", Required: false, Secret: false},
		{Name: "applicationSecret", Label: "ApplicationSecret", Type: "password", Required: false, Secret: true},
		{Name: "consumerKey", Label: "ConsumerKey", Type: "text", Required: false, Secret: false},
		{Name: "clientId", Label: "ClientId", Type: "text", Required: false, Secret: false},
		{Name: "clientSecret", Label: "ClientSecret", Type: "password", Required: false, Secret: true},
	}},
	{ID: "porkbun", Label: "Porkbun", Kind: "access", AccessProviderID: "porkbun", Capabilities: nil, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "secretApiKey", Label: "SecretApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "powerdns", Label: "Powerdns", Kind: "access", AccessProviderID: "powerdns", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "proxmoxve", Label: "Proxmoxve", Kind: "access", AccessProviderID: "proxmoxve", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
		{Name: "apiTokenSecret", Label: "ApiTokenSecret", Type: "password", Required: false, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "qiniu", Label: "Qiniu", Kind: "access", AccessProviderID: "qiniu", Capabilities: nil, AccessFields: []Field{
		{Name: "accessKey", Label: "AccessKey", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "qingcloud", Label: "Qingcloud", Kind: "access", AccessProviderID: "qingcloud", Capabilities: nil, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "rainyun", Label: "Rainyun", Kind: "access", AccessProviderID: "rainyun", Capabilities: nil, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ratpanel", Label: "Ratpanel", Kind: "access", AccessProviderID: "ratpanel", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "accessTokenId", Label: "AccessTokenId", Type: "number", Required: true, Secret: false},
		{Name: "accessToken", Label: "AccessToken", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "rfc2136", Label: "Rfc2136", Kind: "access", AccessProviderID: "rfc2136", Capabilities: nil, AccessFields: []Field{
		{Name: "host", Label: "Host", Type: "text", Required: true, Secret: false},
		{Name: "port", Label: "Port", Type: "number", Required: true, Secret: false},
		{Name: "tsigAlgorithm", Label: "TsigAlgorithm", Type: "text", Required: false, Secret: false},
		{Name: "tsigKey", Label: "TsigKey", Type: "text", Required: false, Secret: false},
		{Name: "tsigSecret", Label: "TsigSecret", Type: "password", Required: false, Secret: true},
	}},
	{ID: "s3", Label: "S3", Kind: "access", AccessProviderID: "s3", Capabilities: nil, AccessFields: []Field{
		{Name: "endpoint", Label: "Endpoint", Type: "text", Required: true, Secret: false},
		{Name: "accessKey", Label: "AccessKey", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
		{Name: "signatureVersion", Label: "SignatureVersion", Type: "text", Required: false, Secret: false},
		{Name: "usePathStyle", Label: "UsePathStyle", Type: "checkbox", Required: false, Secret: false},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "safeline", Label: "Safeline", Kind: "access", AccessProviderID: "safeline", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "samwaf", Label: "Samwaf", Kind: "access", AccessProviderID: "samwaf", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "sectigo", Label: "Sectigo", Kind: "access", AccessProviderID: "sectigo", Capabilities: nil, AccessFields: []Field{}},
	{ID: "slackbot", Label: "Slackbot", Kind: "access", AccessProviderID: "slackbot", Capabilities: nil, AccessFields: []Field{
		{Name: "botToken", Label: "BotToken", Type: "password", Required: true, Secret: true},
		{Name: "channelId", Label: "ChannelId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "spaceship", Label: "Spaceship", Kind: "access", AccessProviderID: "spaceship", Capabilities: nil, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "apiSecret", Label: "ApiSecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ssh", Label: "SSH", Kind: "access", AccessProviderID: "ssh", Capabilities: nil, AccessFields: []Field{
		{Name: "host", Label: "Host", Type: "text", Required: true, Secret: false},
		{Name: "port", Label: "Port", Type: "number", Required: true, Secret: false},
		{Name: "authMethod", Label: "AuthMethod", Type: "text", Required: true, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: false, Secret: true},
		{Name: "key", Label: "Key", Type: "textarea", Required: false, Secret: false},
		{Name: "keyPassphrase", Label: "KeyPassphrase", Type: "password", Required: false, Secret: true},
		{Name: "host", Label: "Host", Type: "text", Required: true, Secret: false},
		{Name: "port", Label: "Port", Type: "number", Required: true, Secret: false},
		{Name: "authMethod", Label: "AuthMethod", Type: "text", Required: true, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: false, Secret: true},
		{Name: "key", Label: "Key", Type: "textarea", Required: false, Secret: false},
		{Name: "keyPassphrase", Label: "KeyPassphrase", Type: "password", Required: false, Secret: true},
		{Name: "jumpServers", Label: "}", Type: "text", Required: false, Secret: false},
	}},
	{ID: "sslcom", Label: "Sslcom", Kind: "access", AccessProviderID: "sslcom", Capabilities: nil, AccessFields: []Field{}},
	{ID: "synologydsm", Label: "Synologydsm", Kind: "access", AccessProviderID: "synologydsm", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: true, Secret: true},
		{Name: "totpSecret", Label: "TotpSecret", Type: "password", Required: false, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "technitiumdns", Label: "Technitiumdns", Kind: "access", AccessProviderID: "technitiumdns", Capabilities: nil, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "telegrambot", Label: "Telegrambot", Kind: "access", AccessProviderID: "telegrambot", Capabilities: nil, AccessFields: []Field{
		{Name: "botToken", Label: "BotToken", Type: "password", Required: true, Secret: true},
		{Name: "chatId", Label: "ChatId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "tencentcloud", Label: "Tencent Cloud", Kind: "access", AccessProviderID: "tencentcloud", Capabilities: nil, AccessFields: []Field{
		{Name: "secretId", Label: "SecretId", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "todaynic", Label: "Todaynic", Kind: "access", AccessProviderID: "todaynic", Capabilities: nil, AccessFields: []Field{
		{Name: "userId", Label: "UserId", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ucloud", Label: "Ucloud", Kind: "access", AccessProviderID: "ucloud", Capabilities: nil, AccessFields: []Field{
		{Name: "privateKey", Label: "PrivateKey", Type: "password", Required: true, Secret: true},
		{Name: "publicKey", Label: "PublicKey", Type: "text", Required: true, Secret: false},
		{Name: "projectId", Label: "ProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "unicloud", Label: "Unicloud", Kind: "access", AccessProviderID: "unicloud", Capabilities: nil, AccessFields: []Field{
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: true, Secret: true},
	}},
	{ID: "upyun", Label: "Upyun", Kind: "access", AccessProviderID: "upyun", Capabilities: nil, AccessFields: []Field{
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: true, Secret: true},
	}},
	{ID: "vercel", Label: "Vercel", Kind: "access", AccessProviderID: "vercel", Capabilities: nil, AccessFields: []Field{
		{Name: "apiAccessToken", Label: "ApiAccessToken", Type: "password", Required: true, Secret: true},
		{Name: "teamId", Label: "TeamId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "volcengine", Label: "Volcengine", Kind: "access", AccessProviderID: "volcengine", Capabilities: nil, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "vultr", Label: "Vultr", Kind: "access", AccessProviderID: "vultr", Capabilities: nil, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "wangsu", Label: "Wangsu", Kind: "access", AccessProviderID: "wangsu", Capabilities: nil, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "webhook", Label: "Webhook", Kind: "access", AccessProviderID: "webhook", Capabilities: nil, AccessFields: []Field{
		{Name: "url", Label: "Url", Type: "text", Required: true, Secret: false},
		{Name: "method", Label: "Method", Type: "text", Required: false, Secret: false},
		{Name: "headers", Label: "HeadersString", Type: "text", Required: false, Secret: false},
		{Name: "data", Label: "DataString", Type: "text", Required: false, Secret: false},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "wecombot", Label: "Wecombot", Kind: "access", AccessProviderID: "wecombot", Capabilities: nil, AccessFields: []Field{
		{Name: "webhookUrl", Label: "WebhookUrl", Type: "text", Required: true, Secret: false},
		{Name: "customPayload", Label: "CustomPayload", Type: "textarea", Required: false, Secret: false},
	}},
	{ID: "westcn", Label: "Westcn", Kind: "access", AccessProviderID: "westcn", Capabilities: nil, AccessFields: []Field{
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "apiPassword", Label: "ApiPassword", Type: "password", Required: true, Secret: true},
	}},
	{ID: "xinnet", Label: "Xinnet", Kind: "access", AccessProviderID: "xinnet", Capabilities: nil, AccessFields: []Field{
		{Name: "agentId", Label: "AgentId", Type: "text", Required: true, Secret: false},
		{Name: "apiPassword", Label: "ApiPassword", Type: "password", Required: true, Secret: true},
	}},
	{ID: "zerossl", Label: "Zerossl", Kind: "access", AccessProviderID: "zerossl", Capabilities: nil, AccessFields: []Field{}},
	{ID: "35cn", Label: "35Cn", Kind: "dns", AccessProviderID: "35cn", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "apiPassword", Label: "ApiPassword", Type: "password", Required: true, Secret: true},
	}},
	{ID: "51dnscom", Label: "51Dnscom", Kind: "dns", AccessProviderID: "51dnscom", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "apiSecret", Label: "ApiSecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "acmedns", Label: "Acmedns", Kind: "dns", AccessProviderID: "acmedns", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "credentials", Label: "Credentials", Type: "password", Required: true, Secret: true},
	}},
	{ID: "acmehttpreq", Label: "Acmehttpreq", Kind: "dns", AccessProviderID: "acmehttpreq", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "endpoint", Label: "Endpoint", Type: "text", Required: true, Secret: false},
		{Name: "mode", Label: "Mode", Type: "text", Required: false, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: false, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: false, Secret: true},
	}},
	{ID: "akamai", Label: "Akamai", Kind: "dns", AccessProviderID: "akamai", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "host", Label: "Host", Type: "text", Required: true, Secret: false},
		{Name: "clientToken", Label: "ClientToken", Type: "password", Required: true, Secret: true},
		{Name: "clientSecret", Label: "ClientSecret", Type: "password", Required: true, Secret: true},
		{Name: "accessToken", Label: "AccessToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "akamai-edgedns", Label: "Akamai Edgedns", Kind: "dns", AccessProviderID: "akamai", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "host", Label: "Host", Type: "text", Required: true, Secret: false},
		{Name: "clientToken", Label: "ClientToken", Type: "password", Required: true, Secret: true},
		{Name: "clientSecret", Label: "ClientSecret", Type: "password", Required: true, Secret: true},
		{Name: "accessToken", Label: "AccessToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "aliyun", Label: "Aliyun", Kind: "dns", AccessProviderID: "aliyun", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-dns", Label: "Aliyun Dns", Kind: "dns", AccessProviderID: "aliyun", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-esa", Label: "Aliyun Esa", Kind: "dns", AccessProviderID: "aliyun", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "arvancloud", Label: "Arvancloud", Kind: "dns", AccessProviderID: "arvancloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "aws", Label: "AWS", Kind: "dns", AccessProviderID: "aws", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "aws-route53", Label: "Aws Route53", Kind: "dns", AccessProviderID: "aws", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "azure", Label: "Azure", Kind: "dns", AccessProviderID: "azure", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "tenantId", Label: "TenantId", Type: "text", Required: true, Secret: false},
		{Name: "clientId", Label: "ClientId", Type: "text", Required: true, Secret: false},
		{Name: "clientSecret", Label: "ClientSecret", Type: "password", Required: true, Secret: true},
		{Name: "subscriptionId", Label: "SubscriptionId", Type: "text", Required: false, Secret: false},
		{Name: "resourceGroupName", Label: "ResourceGroupName", Type: "text", Required: false, Secret: false},
		{Name: "cloudName", Label: "CloudName", Type: "text", Required: false, Secret: false},
	}},
	{ID: "azure-dns", Label: "Azure Dns", Kind: "dns", AccessProviderID: "azure", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "tenantId", Label: "TenantId", Type: "text", Required: true, Secret: false},
		{Name: "clientId", Label: "ClientId", Type: "text", Required: true, Secret: false},
		{Name: "clientSecret", Label: "ClientSecret", Type: "password", Required: true, Secret: true},
		{Name: "subscriptionId", Label: "SubscriptionId", Type: "text", Required: false, Secret: false},
		{Name: "resourceGroupName", Label: "ResourceGroupName", Type: "text", Required: false, Secret: false},
		{Name: "cloudName", Label: "CloudName", Type: "text", Required: false, Secret: false},
	}},
	{ID: "baiducloud", Label: "Baiducloud", Kind: "dns", AccessProviderID: "baiducloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "baiducloud-dns", Label: "Baiducloud Dns", Kind: "dns", AccessProviderID: "baiducloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "bookmyname", Label: "Bookmyname", Kind: "dns", AccessProviderID: "bookmyname", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: true, Secret: true},
	}},
	{ID: "bunny", Label: "Bunny", Kind: "dns", AccessProviderID: "bunny", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "cloudflare", Label: "Cloudflare", Kind: "dns", AccessProviderID: "cloudflare", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "dnsApiToken", Label: "DnsApiToken", Type: "password", Required: true, Secret: true},
		{Name: "zoneApiToken", Label: "ZoneApiToken", Type: "password", Required: false, Secret: true},
	}},
	{ID: "cloudns", Label: "Cloudns", Kind: "dns", AccessProviderID: "cloudns", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "authId", Label: "AuthId", Type: "text", Required: true, Secret: false},
		{Name: "authPassword", Label: "AuthPassword", Type: "password", Required: true, Secret: true},
	}},
	{ID: "cmcccloud", Label: "Cmcccloud", Kind: "dns", AccessProviderID: "cmcccloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "cmcccloud-dns", Label: "Cmcccloud Dns", Kind: "dns", AccessProviderID: "cmcccloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "constellix", Label: "Constellix", Kind: "dns", AccessProviderID: "constellix", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "cpanel", Label: "Cpanel", Kind: "dns", AccessProviderID: "cpanel", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "ctcccloud", Label: "Ctcccloud", Kind: "dns", AccessProviderID: "ctcccloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ctcccloud-smartdns", Label: "Ctcccloud Smartdns", Kind: "dns", AccessProviderID: "ctcccloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "desec", Label: "Desec", Kind: "dns", AccessProviderID: "desec", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "token", Label: "Token", Type: "password", Required: true, Secret: true},
	}},
	{ID: "digitalocean", Label: "Digitalocean", Kind: "dns", AccessProviderID: "digitalocean", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessToken", Label: "AccessToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "dnsexit", Label: "Dnsexit", Kind: "dns", AccessProviderID: "dnsexit", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "dnsla", Label: "Dnsla", Kind: "dns", AccessProviderID: "dnsla", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiId", Label: "ApiId", Type: "text", Required: true, Secret: false},
		{Name: "apiSecret", Label: "ApiSecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "dnsmadeeasy", Label: "Dnsmadeeasy", Kind: "dns", AccessProviderID: "dnsmadeeasy", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "apiSecret", Label: "ApiSecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "duckdns", Label: "Duckdns", Kind: "dns", AccessProviderID: "duckdns", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "token", Label: "Token", Type: "password", Required: true, Secret: true},
	}},
	{ID: "dynu", Label: "Dynu", Kind: "dns", AccessProviderID: "dynu", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "dynv6", Label: "Dynv6", Kind: "dns", AccessProviderID: "dynv6", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "httpToken", Label: "HttpToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "gandinet", Label: "Gandinet", Kind: "dns", AccessProviderID: "gandinet", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "personalAccessToken", Label: "PersonalAccessToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "gcore", Label: "Gcore", Kind: "dns", AccessProviderID: "gcore", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "gname", Label: "Gname", Kind: "dns", AccessProviderID: "gname", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "appId", Label: "AppId", Type: "text", Required: true, Secret: false},
		{Name: "appKey", Label: "AppKey", Type: "text", Required: true, Secret: false},
	}},
	{ID: "godaddy", Label: "Godaddy", Kind: "dns", AccessProviderID: "godaddy", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "apiSecret", Label: "ApiSecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "hetzner", Label: "Hetzner", Kind: "dns", AccessProviderID: "hetzner", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "hostingde", Label: "Hostingde", Kind: "dns", AccessProviderID: "hostingde", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "hostinger", Label: "Hostinger", Kind: "dns", AccessProviderID: "hostinger", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "huaweicloud", Label: "Huaweicloud", Kind: "dns", AccessProviderID: "huaweicloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
		{Name: "enterpriseProjectId", Label: "EnterpriseProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "huaweicloud-dns", Label: "Huaweicloud Dns", Kind: "dns", AccessProviderID: "huaweicloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
		{Name: "enterpriseProjectId", Label: "EnterpriseProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "infomaniak", Label: "Infomaniak", Kind: "dns", AccessProviderID: "infomaniak", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessToken", Label: "AccessToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ionos", Label: "Ionos", Kind: "dns", AccessProviderID: "ionos", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiKeyPublicPrefix", Label: "ApiKeyPublicPrefix", Type: "password", Required: true, Secret: true},
		{Name: "apiKeySecret", Label: "ApiKeySecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "jdcloud", Label: "Jdcloud", Kind: "dns", AccessProviderID: "jdcloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "jdcloud-dns", Label: "Jdcloud Dns", Kind: "dns", AccessProviderID: "jdcloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "linode", Label: "Linode", Kind: "dns", AccessProviderID: "linode", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessToken", Label: "AccessToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "namecheap", Label: "Namecheap", Kind: "dns", AccessProviderID: "namecheap", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "namedotcom", Label: "Namedotcom", Kind: "dns", AccessProviderID: "namedotcom", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "namesilo", Label: "Namesilo", Kind: "dns", AccessProviderID: "namesilo", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "netcup", Label: "Netcup", Kind: "dns", AccessProviderID: "netcup", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "customerNumber", Label: "CustomerNumber", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "apiPassword", Label: "ApiPassword", Type: "password", Required: true, Secret: true},
	}},
	{ID: "netlify", Label: "Netlify", Kind: "dns", AccessProviderID: "netlify", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ns1", Label: "Ns1", Kind: "dns", AccessProviderID: "ns1", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ovhcloud", Label: "Ovhcloud", Kind: "dns", AccessProviderID: "ovhcloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "endpoint", Label: "Endpoint", Type: "text", Required: true, Secret: false},
		{Name: "authMethod", Label: "AuthMethod", Type: "text", Required: true, Secret: false},
		{Name: "applicationKey", Label: "ApplicationKey", Type: "text", Required: false, Secret: false},
		{Name: "applicationSecret", Label: "ApplicationSecret", Type: "password", Required: false, Secret: true},
		{Name: "consumerKey", Label: "ConsumerKey", Type: "text", Required: false, Secret: false},
		{Name: "clientId", Label: "ClientId", Type: "text", Required: false, Secret: false},
		{Name: "clientSecret", Label: "ClientSecret", Type: "password", Required: false, Secret: true},
	}},
	{ID: "porkbun", Label: "Porkbun", Kind: "dns", AccessProviderID: "porkbun", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "secretApiKey", Label: "SecretApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "powerdns", Label: "Powerdns", Kind: "dns", AccessProviderID: "powerdns", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "qingcloud", Label: "Qingcloud", Kind: "dns", AccessProviderID: "qingcloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "qingcloud-dns", Label: "Qingcloud Dns", Kind: "dns", AccessProviderID: "qingcloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "rainyun", Label: "Rainyun", Kind: "dns", AccessProviderID: "rainyun", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "rfc2136", Label: "Rfc2136", Kind: "dns", AccessProviderID: "rfc2136", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "host", Label: "Host", Type: "text", Required: true, Secret: false},
		{Name: "port", Label: "Port", Type: "number", Required: true, Secret: false},
		{Name: "tsigAlgorithm", Label: "TsigAlgorithm", Type: "text", Required: false, Secret: false},
		{Name: "tsigKey", Label: "TsigKey", Type: "text", Required: false, Secret: false},
		{Name: "tsigSecret", Label: "TsigSecret", Type: "password", Required: false, Secret: true},
	}},
	{ID: "spaceship", Label: "Spaceship", Kind: "dns", AccessProviderID: "spaceship", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "apiSecret", Label: "ApiSecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "technitiumdns", Label: "Technitiumdns", Kind: "dns", AccessProviderID: "technitiumdns", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "tencentcloud", Label: "Tencent Cloud", Kind: "dns", AccessProviderID: "tencentcloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "secretId", Label: "SecretId", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "tencentcloud-dns", Label: "Tencentcloud Dns", Kind: "dns", AccessProviderID: "tencentcloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "secretId", Label: "SecretId", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "tencentcloud-eo", Label: "Tencentcloud Eo", Kind: "dns", AccessProviderID: "tencentcloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "secretId", Label: "SecretId", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "todaynic", Label: "Todaynic", Kind: "dns", AccessProviderID: "todaynic", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "userId", Label: "UserId", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ucloud", Label: "Ucloud", Kind: "dns", AccessProviderID: "ucloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "privateKey", Label: "PrivateKey", Type: "password", Required: true, Secret: true},
		{Name: "publicKey", Label: "PublicKey", Type: "text", Required: true, Secret: false},
		{Name: "projectId", Label: "ProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "ucloud-udnr", Label: "Ucloud Udnr", Kind: "dns", AccessProviderID: "ucloud", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "privateKey", Label: "PrivateKey", Type: "password", Required: true, Secret: true},
		{Name: "publicKey", Label: "PublicKey", Type: "text", Required: true, Secret: false},
		{Name: "projectId", Label: "ProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "vercel", Label: "Vercel", Kind: "dns", AccessProviderID: "vercel", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiAccessToken", Label: "ApiAccessToken", Type: "password", Required: true, Secret: true},
		{Name: "teamId", Label: "TeamId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "volcengine", Label: "Volcengine", Kind: "dns", AccessProviderID: "volcengine", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "volcengine-dns", Label: "Volcengine Dns", Kind: "dns", AccessProviderID: "volcengine", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "vultr", Label: "Vultr", Kind: "dns", AccessProviderID: "vultr", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "westcn", Label: "Westcn", Kind: "dns", AccessProviderID: "westcn", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "apiPassword", Label: "ApiPassword", Type: "password", Required: true, Secret: true},
	}},
	{ID: "xinnet", Label: "Xinnet", Kind: "dns", AccessProviderID: "xinnet", Capabilities: []string{"dns"}, AccessFields: []Field{
		{Name: "agentId", Label: "AgentId", Type: "text", Required: true, Secret: false},
		{Name: "apiPassword", Label: "ApiPassword", Type: "password", Required: true, Secret: true},
	}},
	{ID: "1panel", Label: "1Panel", Kind: "deploy", AccessProviderID: "1panel", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiVersion", Label: "ApiVersion", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "1panel-console", Label: "1Panel Console", Kind: "deploy", AccessProviderID: "1panel", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiVersion", Label: "ApiVersion", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "aliyun-alb", Label: "Aliyun Alb", Kind: "deploy", AccessProviderID: "aliyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-apigw", Label: "Aliyun Apigw", Kind: "deploy", AccessProviderID: "aliyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-cas", Label: "Aliyun Cas", Kind: "deploy", AccessProviderID: "aliyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-casdeploy", Label: "Aliyun Casdeploy", Kind: "deploy", AccessProviderID: "aliyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-cdn", Label: "Aliyun Cdn", Kind: "deploy", AccessProviderID: "aliyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-clb", Label: "Aliyun Clb", Kind: "deploy", AccessProviderID: "aliyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-dcdn", Label: "Aliyun Dcdn", Kind: "deploy", AccessProviderID: "aliyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-ddospro", Label: "Aliyun Ddospro", Kind: "deploy", AccessProviderID: "aliyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-esa", Label: "Aliyun Esa", Kind: "deploy", AccessProviderID: "aliyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-esasaas", Label: "Aliyun Esasaas", Kind: "deploy", AccessProviderID: "aliyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-fc", Label: "Aliyun Fc", Kind: "deploy", AccessProviderID: "aliyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-ga", Label: "Aliyun Ga", Kind: "deploy", AccessProviderID: "aliyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-live", Label: "Aliyun Live", Kind: "deploy", AccessProviderID: "aliyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-nlb", Label: "Aliyun Nlb", Kind: "deploy", AccessProviderID: "aliyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-oss", Label: "Aliyun Oss", Kind: "deploy", AccessProviderID: "aliyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-vod", Label: "Aliyun Vod", Kind: "deploy", AccessProviderID: "aliyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "aliyun-waf", Label: "Aliyun Waf", Kind: "deploy", AccessProviderID: "aliyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "resourceGroupId", Label: "ResourceGroupId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "apisix", Label: "Apisix", Kind: "deploy", AccessProviderID: "apisix", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "aws-acm", Label: "Aws Acm", Kind: "deploy", AccessProviderID: "aws", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "aws-cloudfront", Label: "Aws Cloudfront", Kind: "deploy", AccessProviderID: "aws", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "aws-iam", Label: "Aws Iam", Kind: "deploy", AccessProviderID: "aws", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "azure-keyvault", Label: "Azure Keyvault", Kind: "deploy", AccessProviderID: "azure", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "tenantId", Label: "TenantId", Type: "text", Required: true, Secret: false},
		{Name: "clientId", Label: "ClientId", Type: "text", Required: true, Secret: false},
		{Name: "clientSecret", Label: "ClientSecret", Type: "password", Required: true, Secret: true},
		{Name: "subscriptionId", Label: "SubscriptionId", Type: "text", Required: false, Secret: false},
		{Name: "resourceGroupName", Label: "ResourceGroupName", Type: "text", Required: false, Secret: false},
		{Name: "cloudName", Label: "CloudName", Type: "text", Required: false, Secret: false},
	}},
	{ID: "baiducloud-appblb", Label: "Baiducloud Appblb", Kind: "deploy", AccessProviderID: "baiducloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "baiducloud-blb", Label: "Baiducloud Blb", Kind: "deploy", AccessProviderID: "baiducloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "baiducloud-cdn", Label: "Baiducloud Cdn", Kind: "deploy", AccessProviderID: "baiducloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "baiducloud-cert", Label: "Baiducloud Cert", Kind: "deploy", AccessProviderID: "baiducloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "baishan-cdn", Label: "Baishan Cdn", Kind: "deploy", AccessProviderID: "baishan", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "baotapanel", Label: "Baotapanel", Kind: "deploy", AccessProviderID: "baotapanel", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "baotapanel-console", Label: "Baotapanel Console", Kind: "deploy", AccessProviderID: "baotapanel", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "baotapanelgo", Label: "Baotapanelgo", Kind: "deploy", AccessProviderID: "baotapanelgo", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "baotapanelgo-console", Label: "Baotapanelgo Console", Kind: "deploy", AccessProviderID: "baotapanelgo", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "baotawaf", Label: "Baotawaf", Kind: "deploy", AccessProviderID: "baotawaf", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "baotawaf-console", Label: "Baotawaf Console", Kind: "deploy", AccessProviderID: "baotawaf", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "bunny-cdn", Label: "Bunny Cdn", Kind: "deploy", AccessProviderID: "bunny", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "byteplus-cdn", Label: "Byteplus Cdn", Kind: "deploy", AccessProviderID: "byteplus", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKey", Label: "AccessKey", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "cachefly", Label: "Cachefly", Kind: "deploy", AccessProviderID: "cachefly", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "cdnfly", Label: "Cdnfly", Kind: "deploy", AccessProviderID: "cdnfly", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "apiSecret", Label: "ApiSecret", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "cpanel", Label: "Cpanel", Kind: "deploy", AccessProviderID: "cpanel", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "ctcccloud-ao", Label: "Ctcccloud Ao", Kind: "deploy", AccessProviderID: "ctcccloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ctcccloud-cdn", Label: "Ctcccloud Cdn", Kind: "deploy", AccessProviderID: "ctcccloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ctcccloud-cms", Label: "Ctcccloud Cms", Kind: "deploy", AccessProviderID: "ctcccloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ctcccloud-elb", Label: "Ctcccloud Elb", Kind: "deploy", AccessProviderID: "ctcccloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ctcccloud-faas", Label: "Ctcccloud Faas", Kind: "deploy", AccessProviderID: "ctcccloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ctcccloud-icdn", Label: "Ctcccloud Icdn", Kind: "deploy", AccessProviderID: "ctcccloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ctcccloud-ldvn", Label: "Ctcccloud Ldvn", Kind: "deploy", AccessProviderID: "ctcccloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "dogecloud-cdn", Label: "Dogecloud Cdn", Kind: "deploy", AccessProviderID: "dogecloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKey", Label: "AccessKey", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "dokploy", Label: "Dokploy", Kind: "deploy", AccessProviderID: "dokploy", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "flexcdn", Label: "Flexcdn", Kind: "deploy", AccessProviderID: "flexcdn", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiRole", Label: "ApiRole", Type: "text", Required: true, Secret: false},
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKey", Label: "AccessKey", Type: "text", Required: true, Secret: false},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "flyio", Label: "Flyio", Kind: "deploy", AccessProviderID: "flyio", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ftp", Label: "FTP", Kind: "deploy", AccessProviderID: "ftp", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "host", Label: "Host", Type: "text", Required: true, Secret: false},
		{Name: "port", Label: "Port", Type: "number", Required: true, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: false, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: false, Secret: true},
	}},
	{ID: "gcore-cdn", Label: "Gcore Cdn", Kind: "deploy", AccessProviderID: "gcore", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "goedge", Label: "Goedge", Kind: "deploy", AccessProviderID: "goedge", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiRole", Label: "ApiRole", Type: "text", Required: true, Secret: false},
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKey", Label: "AccessKey", Type: "text", Required: true, Secret: false},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "huaweicloud-aad", Label: "Huaweicloud Aad", Kind: "deploy", AccessProviderID: "huaweicloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
		{Name: "enterpriseProjectId", Label: "EnterpriseProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "huaweicloud-apig", Label: "Huaweicloud Apig", Kind: "deploy", AccessProviderID: "huaweicloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
		{Name: "enterpriseProjectId", Label: "EnterpriseProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "huaweicloud-cdn", Label: "Huaweicloud Cdn", Kind: "deploy", AccessProviderID: "huaweicloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
		{Name: "enterpriseProjectId", Label: "EnterpriseProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "huaweicloud-elb", Label: "Huaweicloud Elb", Kind: "deploy", AccessProviderID: "huaweicloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
		{Name: "enterpriseProjectId", Label: "EnterpriseProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "huaweicloud-live", Label: "Huaweicloud Live", Kind: "deploy", AccessProviderID: "huaweicloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
		{Name: "enterpriseProjectId", Label: "EnterpriseProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "huaweicloud-obs", Label: "Huaweicloud Obs", Kind: "deploy", AccessProviderID: "huaweicloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
		{Name: "enterpriseProjectId", Label: "EnterpriseProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "huaweicloud-scm", Label: "Huaweicloud Scm", Kind: "deploy", AccessProviderID: "huaweicloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
		{Name: "enterpriseProjectId", Label: "EnterpriseProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "huaweicloud-waf", Label: "Huaweicloud Waf", Kind: "deploy", AccessProviderID: "huaweicloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
		{Name: "enterpriseProjectId", Label: "EnterpriseProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "jdcloud-alb", Label: "Jdcloud Alb", Kind: "deploy", AccessProviderID: "jdcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "jdcloud-cdn", Label: "Jdcloud Cdn", Kind: "deploy", AccessProviderID: "jdcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "jdcloud-live", Label: "Jdcloud Live", Kind: "deploy", AccessProviderID: "jdcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "jdcloud-vod", Label: "Jdcloud Vod", Kind: "deploy", AccessProviderID: "jdcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "jdcloud-ssl", Label: "Jdcloud Ssl", Kind: "deploy", AccessProviderID: "jdcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "jdcloud-waf", Label: "Jdcloud Waf", Kind: "deploy", AccessProviderID: "jdcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
	}},
	{ID: "kong", Label: "Kong", Kind: "deploy", AccessProviderID: "kong", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: false, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "k8s-secret", Label: "K8S Secret", Kind: "deploy", AccessProviderID: "k8s", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "kubeConfig", Label: "KubeConfig", Type: "text", Required: false, Secret: false},
	}},
	{ID: "ksyun-cdn", Label: "Ksyun Cdn", Kind: "deploy", AccessProviderID: "ksyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ksyun-slb", Label: "Ksyun Slb", Kind: "deploy", AccessProviderID: "ksyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "lecdn", Label: "Lecdn", Kind: "deploy", AccessProviderID: "lecdn", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiVersion", Label: "ApiVersion", Type: "text", Required: true, Secret: false},
		{Name: "apiRole", Label: "ApiRole", Type: "text", Required: true, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "local", Label: "Local", Kind: "deploy", AccessProviderID: "local", Capabilities: []string{"deploy"}, AccessFields: []Field{}},
	{ID: "mohua-mvh", Label: "Mohua Mvh", Kind: "deploy", AccessProviderID: "mohua", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "apiPassword", Label: "ApiPassword", Type: "password", Required: true, Secret: true},
	}},
	{ID: "netlify", Label: "Netlify", Kind: "deploy", AccessProviderID: "netlify", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
	}},
	{ID: "nginxproxymanager", Label: "Nginxproxymanager", Kind: "deploy", AccessProviderID: "nginxproxymanager", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "authMethod", Label: "AuthMethod", Type: "text", Required: true, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: false, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: false, Secret: true},
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: false, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "proxmoxve", Label: "Proxmoxve", Kind: "deploy", AccessProviderID: "proxmoxve", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
		{Name: "apiTokenSecret", Label: "ApiTokenSecret", Type: "password", Required: false, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "qiniu-cdn", Label: "Qiniu Cdn", Kind: "deploy", AccessProviderID: "qiniu", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKey", Label: "AccessKey", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "qiniu-kodo", Label: "Qiniu Kodo", Kind: "deploy", AccessProviderID: "qiniu", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKey", Label: "AccessKey", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "qiniu-pili", Label: "Qiniu Pili", Kind: "deploy", AccessProviderID: "qiniu", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKey", Label: "AccessKey", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "rainyun-rcdn", Label: "Rainyun Rcdn", Kind: "deploy", AccessProviderID: "rainyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "rainyun-sslcenter", Label: "Rainyun Sslcenter", Kind: "deploy", AccessProviderID: "rainyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ratpanel", Label: "Ratpanel", Kind: "deploy", AccessProviderID: "ratpanel", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "accessTokenId", Label: "AccessTokenId", Type: "number", Required: true, Secret: false},
		{Name: "accessToken", Label: "AccessToken", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "ratpanel-console", Label: "Ratpanel Console", Kind: "deploy", AccessProviderID: "ratpanel", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "accessTokenId", Label: "AccessTokenId", Type: "number", Required: true, Secret: false},
		{Name: "accessToken", Label: "AccessToken", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "s3", Label: "S3", Kind: "deploy", AccessProviderID: "s3", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "endpoint", Label: "Endpoint", Type: "text", Required: true, Secret: false},
		{Name: "accessKey", Label: "AccessKey", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
		{Name: "signatureVersion", Label: "SignatureVersion", Type: "text", Required: false, Secret: false},
		{Name: "usePathStyle", Label: "UsePathStyle", Type: "checkbox", Required: false, Secret: false},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "safeline", Label: "Safeline", Kind: "deploy", AccessProviderID: "safeline", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiToken", Label: "ApiToken", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "samwaf", Label: "Samwaf", Kind: "deploy", AccessProviderID: "samwaf", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "ssh", Label: "SSH", Kind: "deploy", AccessProviderID: "ssh", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "host", Label: "Host", Type: "text", Required: true, Secret: false},
		{Name: "port", Label: "Port", Type: "number", Required: true, Secret: false},
		{Name: "authMethod", Label: "AuthMethod", Type: "text", Required: true, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: false, Secret: true},
		{Name: "key", Label: "Key", Type: "textarea", Required: false, Secret: false},
		{Name: "keyPassphrase", Label: "KeyPassphrase", Type: "password", Required: false, Secret: true},
		{Name: "host", Label: "Host", Type: "text", Required: true, Secret: false},
		{Name: "port", Label: "Port", Type: "number", Required: true, Secret: false},
		{Name: "authMethod", Label: "AuthMethod", Type: "text", Required: true, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: false, Secret: true},
		{Name: "key", Label: "Key", Type: "textarea", Required: false, Secret: false},
		{Name: "keyPassphrase", Label: "KeyPassphrase", Type: "password", Required: false, Secret: true},
		{Name: "jumpServers", Label: "}", Type: "text", Required: false, Secret: false},
	}},
	{ID: "synologydsm", Label: "Synologydsm", Kind: "deploy", AccessProviderID: "synologydsm", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "serverUrl", Label: "ServerUrl", Type: "text", Required: true, Secret: false},
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: true, Secret: true},
		{Name: "totpSecret", Label: "TotpSecret", Type: "password", Required: false, Secret: true},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
	{ID: "tencentcloud-cdn", Label: "Tencentcloud Cdn", Kind: "deploy", AccessProviderID: "tencentcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "secretId", Label: "SecretId", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "tencentcloud-clb", Label: "Tencentcloud Clb", Kind: "deploy", AccessProviderID: "tencentcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "secretId", Label: "SecretId", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "tencentcloud-cos", Label: "Tencentcloud Cos", Kind: "deploy", AccessProviderID: "tencentcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "secretId", Label: "SecretId", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "tencentcloud-css", Label: "Tencentcloud Css", Kind: "deploy", AccessProviderID: "tencentcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "secretId", Label: "SecretId", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "tencentcloud-ecdn", Label: "Tencentcloud Ecdn", Kind: "deploy", AccessProviderID: "tencentcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "secretId", Label: "SecretId", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "tencentcloud-eo", Label: "Tencentcloud Eo", Kind: "deploy", AccessProviderID: "tencentcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "secretId", Label: "SecretId", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "tencentcloud-gaap", Label: "Tencentcloud Gaap", Kind: "deploy", AccessProviderID: "tencentcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "secretId", Label: "SecretId", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "tencentcloud-scf", Label: "Tencentcloud Scf", Kind: "deploy", AccessProviderID: "tencentcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "secretId", Label: "SecretId", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "tencentcloud-ssl", Label: "Tencentcloud Ssl", Kind: "deploy", AccessProviderID: "tencentcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "secretId", Label: "SecretId", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "tencentcloud-ssldeploy", Label: "Tencentcloud Ssldeploy", Kind: "deploy", AccessProviderID: "tencentcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "secretId", Label: "SecretId", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "tencentcloud-sslupdate", Label: "Tencentcloud Sslupdate", Kind: "deploy", AccessProviderID: "tencentcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "secretId", Label: "SecretId", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "tencentcloud-vod", Label: "Tencentcloud Vod", Kind: "deploy", AccessProviderID: "tencentcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "secretId", Label: "SecretId", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "tencentcloud-waf", Label: "Tencentcloud Waf", Kind: "deploy", AccessProviderID: "tencentcloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "secretId", Label: "SecretId", Type: "text", Required: true, Secret: false},
		{Name: "secretKey", Label: "SecretKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "ucloud-ualb", Label: "Ucloud Ualb", Kind: "deploy", AccessProviderID: "ucloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "privateKey", Label: "PrivateKey", Type: "password", Required: true, Secret: true},
		{Name: "publicKey", Label: "PublicKey", Type: "text", Required: true, Secret: false},
		{Name: "projectId", Label: "ProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "ucloud-ucdn", Label: "Ucloud Ucdn", Kind: "deploy", AccessProviderID: "ucloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "privateKey", Label: "PrivateKey", Type: "password", Required: true, Secret: true},
		{Name: "publicKey", Label: "PublicKey", Type: "text", Required: true, Secret: false},
		{Name: "projectId", Label: "ProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "ucloud-uclb", Label: "Ucloud Uclb", Kind: "deploy", AccessProviderID: "ucloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "privateKey", Label: "PrivateKey", Type: "password", Required: true, Secret: true},
		{Name: "publicKey", Label: "PublicKey", Type: "text", Required: true, Secret: false},
		{Name: "projectId", Label: "ProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "ucloud-uewaf", Label: "Ucloud Uewaf", Kind: "deploy", AccessProviderID: "ucloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "privateKey", Label: "PrivateKey", Type: "password", Required: true, Secret: true},
		{Name: "publicKey", Label: "PublicKey", Type: "text", Required: true, Secret: false},
		{Name: "projectId", Label: "ProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "ucloud-pathx", Label: "Ucloud Pathx", Kind: "deploy", AccessProviderID: "ucloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "privateKey", Label: "PrivateKey", Type: "password", Required: true, Secret: true},
		{Name: "publicKey", Label: "PublicKey", Type: "text", Required: true, Secret: false},
		{Name: "projectId", Label: "ProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "ucloud-us3", Label: "Ucloud Us3", Kind: "deploy", AccessProviderID: "ucloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "privateKey", Label: "PrivateKey", Type: "password", Required: true, Secret: true},
		{Name: "publicKey", Label: "PublicKey", Type: "text", Required: true, Secret: false},
		{Name: "projectId", Label: "ProjectId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "unicloud-webhost", Label: "Unicloud Webhost", Kind: "deploy", AccessProviderID: "unicloud", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: true, Secret: true},
	}},
	{ID: "upyun-cdn", Label: "Upyun Cdn", Kind: "deploy", AccessProviderID: "upyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: true, Secret: true},
	}},
	{ID: "upyun-file", Label: "Upyun File", Kind: "deploy", AccessProviderID: "upyun", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "username", Label: "Username", Type: "text", Required: true, Secret: false},
		{Name: "password", Label: "Password", Type: "password", Required: true, Secret: true},
	}},
	{ID: "vercel", Label: "Vercel", Kind: "deploy", AccessProviderID: "vercel", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "apiAccessToken", Label: "ApiAccessToken", Type: "password", Required: true, Secret: true},
		{Name: "teamId", Label: "TeamId", Type: "text", Required: false, Secret: false},
	}},
	{ID: "volcengine-alb", Label: "Volcengine Alb", Kind: "deploy", AccessProviderID: "volcengine", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "volcengine-apig", Label: "Volcengine Apig", Kind: "deploy", AccessProviderID: "volcengine", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "volcengine-cdn", Label: "Volcengine Cdn", Kind: "deploy", AccessProviderID: "volcengine", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "volcengine-certcenter", Label: "Volcengine Certcenter", Kind: "deploy", AccessProviderID: "volcengine", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "volcengine-clb", Label: "Volcengine Clb", Kind: "deploy", AccessProviderID: "volcengine", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "volcengine-dcdn", Label: "Volcengine Dcdn", Kind: "deploy", AccessProviderID: "volcengine", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "volcengine-imagex", Label: "Volcengine Imagex", Kind: "deploy", AccessProviderID: "volcengine", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "volcengine-live", Label: "Volcengine Live", Kind: "deploy", AccessProviderID: "volcengine", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "volcengine-tos", Label: "Volcengine Tos", Kind: "deploy", AccessProviderID: "volcengine", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "volcengine-vod", Label: "Volcengine Vod", Kind: "deploy", AccessProviderID: "volcengine", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "volcengine-waf", Label: "Volcengine Waf", Kind: "deploy", AccessProviderID: "volcengine", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "secretAccessKey", Label: "SecretAccessKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "wangsu-cdn", Label: "Wangsu Cdn", Kind: "deploy", AccessProviderID: "wangsu", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "wangsu-cdnpro", Label: "Wangsu Cdnpro", Kind: "deploy", AccessProviderID: "wangsu", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "wangsu-certificate", Label: "Wangsu Certificate", Kind: "deploy", AccessProviderID: "wangsu", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "accessKeyId", Label: "AccessKeyId", Type: "text", Required: true, Secret: false},
		{Name: "accessKeySecret", Label: "AccessKeySecret", Type: "password", Required: true, Secret: true},
		{Name: "apiKey", Label: "ApiKey", Type: "password", Required: true, Secret: true},
	}},
	{ID: "webhook", Label: "Webhook", Kind: "deploy", AccessProviderID: "webhook", Capabilities: []string{"deploy"}, AccessFields: []Field{
		{Name: "url", Label: "Url", Type: "text", Required: true, Secret: false},
		{Name: "method", Label: "Method", Type: "text", Required: false, Secret: false},
		{Name: "headers", Label: "HeadersString", Type: "text", Required: false, Secret: false},
		{Name: "data", Label: "DataString", Type: "text", Required: false, Secret: false},
		{Name: "allowInsecureConnections", Label: "AllowInsecureConnections", Type: "checkbox", Required: false, Secret: false},
	}},
}
