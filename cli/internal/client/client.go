package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"

	"easyssl/cli/internal/config"
)

type AuthMode int

const (
	AuthAuto AuthMode = iota
	AuthBearer
	AuthAPIKey
	AuthNone
)

// Options controls HTTP client behavior.
type Options struct {
	Timeout time.Duration
	Verbose bool
	Trace   bool
	Stderr  io.Writer
}

// Client is a thin HTTP client for the EasySSL API.
type Client struct {
	baseURL *url.URL
	token   string
	apiKey  string
	http    *http.Client
	opts    Options
}

// APIResponse is the server envelope.
type APIResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// Result is a decoded API call result.
type Result struct {
	Envelope APIResponse
	RawBody  []byte
}

// New creates a new API client from the current config.
func New(cfg config.Config, opts Options) (*Client, error) {
	if opts.Timeout <= 0 {
		opts.Timeout = 30 * time.Second
	}
	if opts.Stderr == nil {
		opts.Stderr = io.Discard
	}

	base, err := normalizeBaseURL(cfg.Server)
	if err != nil {
		return nil, err
	}
	return &Client{
		baseURL: base,
		token:   cfg.Token,
		apiKey:  cfg.APIKey,
		http:    &http.Client{Timeout: opts.Timeout},
		opts:    opts,
	}, nil
}

func normalizeBaseURL(raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		raw = config.DefaultServer
	}
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parse server url: %w", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("invalid server url: %s", raw)
	}
	u.Path = strings.TrimRight(u.Path, "/")
	if u.Path == "/api" || u.Path == "/openapi" {
		u.Path = ""
	}
	return u, nil
}

// SetToken updates the JWT token used for authenticated requests.
func (c *Client) SetToken(token string) {
	c.token = token
}

// SetAPIKey updates the API key used for OpenAPI requests.
func (c *Client) SetAPIKey(key string) {
	c.apiKey = key
}

// Do sends an API request with the standard response envelope.
func (c *Client) Do(method, reqPath string, query map[string]string, body any, auth AuthMode) (Result, error) {
	var out Result
	u := *c.baseURL
	u.Path = resolvePath(c.baseURL.Path, reqPath)

	vals := u.Query()
	for k, v := range query {
		if strings.TrimSpace(v) == "" {
			continue
		}
		vals.Set(k, v)
	}
	u.RawQuery = vals.Encode()

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return out, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return out, fmt.Errorf("create request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	switch auth {
	case AuthAPIKey:
		if c.apiKey != "" {
			req.Header.Set("X-API-Key", c.apiKey)
		}
	case AuthBearer:
		if c.token != "" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}
	case AuthNone:
		// no auth header
	case AuthAuto:
		if c.apiKey != "" {
			req.Header.Set("X-API-Key", c.apiKey)
		} else if c.token != "" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}
	}

	start := time.Now()
	if c.opts.Verbose {
		fmt.Fprintf(c.opts.Stderr, "[easyssl] request %s %s\n", method, u.String())
		fmt.Fprintf(c.opts.Stderr, "[easyssl] auth-mode=%s\n", authModeString(auth))
		headers := safeHeaders(req.Header)
		fmt.Fprintf(c.opts.Stderr, "[easyssl] headers %s\n", headers)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return out, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	elapsed := time.Since(start)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return out, fmt.Errorf("read response: %w", err)
	}
	out.RawBody = data

	ct := strings.ToLower(resp.Header.Get("Content-Type"))
	requestID := firstNonEmpty(
		resp.Header.Get("X-Request-Id"),
		resp.Header.Get("X-Request-ID"),
		resp.Header.Get("X-Correlation-Id"),
		resp.Header.Get("X-Correlation-ID"),
		resp.Header.Get("Eo-Log-Uuid"),
	)
	if c.opts.Verbose {
		fmt.Fprintf(c.opts.Stderr, "[easyssl] response status=%d elapsed=%s content-type=%s request-id=%s\n", resp.StatusCode, elapsed.Round(time.Millisecond), ct, defaultIfEmpty(requestID, "-"))
		if c.opts.Trace {
			fmt.Fprintf(c.opts.Stderr, "[easyssl] response-body=%s\n", abbreviate(string(data), 1000))
		}
	}

	if !strings.Contains(ct, "application/json") || isLikelyHTML(data) {
		return out, fmt.Errorf(
			"unexpected response from server (status=%d contentType=%q url=%q requestId=%q bodySnippet=%q); check --server points to EasySSL API root and not a web page fallback",
			resp.StatusCode,
			ct,
			u.String(),
			defaultIfEmpty(requestID, "-"),
			abbreviate(string(data), 200),
		)
	}

	if err := json.Unmarshal(data, &out.Envelope); err != nil {
		return out, fmt.Errorf("decode response json: %w", err)
	}

	if resp.StatusCode >= 400 || out.Envelope.Code >= 400 {
		msg := out.Envelope.Msg
		if msg == "" {
			msg = string(data)
		}
		if out.Envelope.Code == 0 {
			out.Envelope.Code = resp.StatusCode
		}
		return out, fmt.Errorf("api error code=%d status=%d msg=%s", out.Envelope.Code, resp.StatusCode, msg)
	}

	return out, nil
}

func safeHeaders(h http.Header) string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		v := strings.Join(h.Values(k), ",")
		lk := strings.ToLower(k)
		if lk == "authorization" || lk == "x-api-key" {
			v = "***"
		}
		parts = append(parts, k+"="+v)
	}
	return strings.Join(parts, " ")
}

func abbreviate(s string, limit int) string {
	s = strings.TrimSpace(s)
	if len(s) <= limit {
		return s
	}
	return s[:limit] + "..."
}

func resolvePath(basePath, reqPath string) string {
	basePath = strings.TrimSpace(basePath)
	reqPath = strings.TrimSpace(reqPath)
	if reqPath == "" {
		if basePath == "" {
			return "/"
		}
		if strings.HasPrefix(basePath, "/") {
			return basePath
		}
		return "/" + basePath
	}
	if strings.HasPrefix(reqPath, "/") {
		if basePath == "" || basePath == "/" {
			return reqPath
		}
		if reqPath == basePath || strings.HasPrefix(reqPath, basePath+"/") {
			return reqPath
		}
		return path.Join(basePath, reqPath)
	}
	fullPath := path.Join(basePath, reqPath)
	if !strings.HasPrefix(fullPath, "/") {
		return "/" + fullPath
	}
	return fullPath
}

func isLikelyHTML(data []byte) bool {
	s := strings.ToLower(strings.TrimSpace(string(data)))
	return strings.HasPrefix(s, "<!doctype html") || strings.HasPrefix(s, "<html")
}

func authModeString(m AuthMode) string {
	switch m {
	case AuthAPIKey:
		return "api_key"
	case AuthBearer:
		return "bearer"
	case AuthNone:
		return "none"
	case AuthAuto:
		return "auto"
	default:
		return "unknown"
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func defaultIfEmpty(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}
