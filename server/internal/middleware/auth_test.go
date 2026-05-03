package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"easyssl/server/internal/model"

	"github.com/gin-gonic/gin"
)

type stubVerifier struct {
	auth *model.AuthContext
}

func (s stubVerifier) VerifyAPIKey(ctx context.Context, rawKey string) (*model.AuthContext, error) {
	_ = ctx
	_ = rawKey
	return s.auth, nil
}

func TestRequireAPIKeyAuthRejectsBearerOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequireAPIKeyAuth(stubVerifier{auth: &model.AuthContext{UserID: "u1", Role: model.RoleUser, AuthType: "api_key"}}))
	r.GET("/open", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/open", nil)
	req.Header.Set("Authorization", "Bearer dummy")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected http %d, got %d", http.StatusOK, w.Code)
	}
	var resp struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected response code %d, got %d", http.StatusUnauthorized, resp.Code)
	}
}

func TestRequireAPIKeyAuthAcceptsXAPIKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequireAPIKeyAuth(stubVerifier{auth: &model.AuthContext{UserID: "u1", Role: model.RoleUser, AuthType: "api_key", APIKeyID: "k1"}}))
	r.GET("/open", func(c *gin.Context) {
		if got := c.GetString(ContextUserID); got != "u1" {
			t.Fatalf("unexpected user id: %s", got)
		}
		if got := c.GetString(ContextAuthType); got != "api_key" {
			t.Fatalf("unexpected auth type: %s", got)
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/open", nil)
	req.Header.Set("X-API-Key", "esk_test")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
}
