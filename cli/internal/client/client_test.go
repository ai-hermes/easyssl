package client

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"easyssl/cli/internal/config"
)

func TestDoPathJoinAvoidsDoubleSlash(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/certificates" {
			t.Fatalf("path=%s, want /api/certificates", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"msg":"ok","data":{}}`))
	})
	ts := httptest.NewServer(h)
	defer ts.Close()

	c, err := New(config.Config{Server: ts.URL + "/"}, Options{})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	if _, err := c.Do("GET", "/api/certificates", nil, nil, AuthNone); err != nil {
		t.Fatalf("Do error: %v", err)
	}
}

func TestDoRejectsHTMLResponse(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html>fallback</html>"))
	})
	ts := httptest.NewServer(h)
	defer ts.Close()

	c, err := New(config.Config{Server: ts.URL}, Options{})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	_, err = c.Do("GET", "/api/certificates", nil, nil, AuthNone)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unexpected response from server") {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestDoRejectsHTMLBodyEvenWithJSONContentType(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<!doctype html><html><body>fallback</body></html>"))
	})
	ts := httptest.NewServer(h)
	defer ts.Close()

	c, err := New(config.Config{Server: ts.URL}, Options{})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	_, err = c.Do("GET", "/api/certificates", nil, nil, AuthNone)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unexpected response from server") {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestNormalizeBaseURLStripsAPISuffix(t *testing.T) {
	c, err := New(config.Config{Server: "https://easyssl.example.com/api"}, Options{})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	if got := c.baseURL.Path; got != "" {
		t.Fatalf("base path = %q, want empty", got)
	}
}

func TestResolvePathAvoidsDuplicateBasePath(t *testing.T) {
	got := resolvePath("/api", "/api/certificates")
	if got != "/api/certificates" {
		t.Fatalf("path = %q, want /api/certificates", got)
	}

	got = resolvePath("", "/openapi/accesses")
	if got != "/openapi/accesses" {
		t.Fatalf("path = %q, want /openapi/accesses", got)
	}

	got = resolvePath("/prefix", "/api/accesses")
	if got != "/prefix/api/accesses" {
		t.Fatalf("path = %q, want /prefix/api/accesses", got)
	}
}

func TestDoPathJoinKeepsPrefixedBasePath(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/prefix/api/certificates" {
			t.Fatalf("path=%s, want /prefix/api/certificates", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"msg":"ok","data":{}}`))
	})
	ts := httptest.NewServer(h)
	defer ts.Close()

	c, err := New(config.Config{Server: ts.URL + "/prefix"}, Options{})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	if _, err := c.Do("GET", "/api/certificates", nil, nil, AuthNone); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}
