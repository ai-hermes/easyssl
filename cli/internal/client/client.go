package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"easyssl/cli/internal/config"
)

// Client is a thin HTTP client for the EasySSL API.
type Client struct {
	baseURL string
	token   string
	apiKey  string
	http    *http.Client
}

// New creates a new API client from the current config.
func New(cfg config.Config) *Client {
	return &Client{
		baseURL: cfg.Server,
		token:   cfg.Token,
		apiKey:  cfg.APIKey,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

// SetToken updates the JWT token used for authenticated requests.
func (c *Client) SetToken(token string) {
	c.token = token
}

// SetAPIKey updates the API key used for OpenAPI requests.
func (c *Client) SetAPIKey(key string) {
	c.apiKey = key
}

func (c *Client) do(method, path string, body any, headers map[string]string) (*http.Response, error) {
	url := c.baseURL + path
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return c.http.Do(req)
}

// Do performs a raw HTTP request and returns the response body.
func (c *Client) Do(method, path string, body any) ([]byte, error) {
	resp, err := c.do(method, path, body, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

// Login authenticates with email and password and returns the JWT token.
func (c *Client) Login(email, password string) (string, error) {
	resp, err := c.do("POST", "/api/auth/login", map[string]string{
		"email":    email,
		"password": password,
	}, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("login failed: HTTP %d: %s", resp.StatusCode, string(data))
	}
	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}
	if result.Code != 0 && result.Code != 200 {
		return "", fmt.Errorf("login failed: %s", result.Msg)
	}
	return result.Data.Token, nil
}

// Me returns the current user info.
func (c *Client) Me() ([]byte, error) {
	return c.Do("GET", "/api/auth/me", nil)
}

// ListWorkflows returns the list of workflows.
func (c *Client) ListWorkflows() ([]byte, error) {
	return c.Do("GET", "/api/workflows", nil)
}

// ListCertificates returns the list of certificates.
func (c *Client) ListCertificates() ([]byte, error) {
	return c.Do("GET", "/api/certificates", nil)
}

// ApplyCertificate triggers a certificate application via OpenAPI.
func (c *Client) ApplyCertificate(req map[string]any) ([]byte, error) {
	return c.Do("POST", "/openapi/certificates/apply", req)
}

// GetCertificateRun returns the status of a certificate run.
func (c *Client) GetCertificateRun(runID string) ([]byte, error) {
	return c.Do("GET", "/openapi/certificates/runs/"+runID, nil)
}

// DownloadCertificate downloads a certificate archive.
func (c *Client) DownloadCertificate(id string) ([]byte, error) {
	return c.Do("POST", "/openapi/certificates/"+id+"/download", nil)
}
