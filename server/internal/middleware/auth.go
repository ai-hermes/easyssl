package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"easyssl/server/internal/model"
	"easyssl/server/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ContextUserID   = "userId"
	ContextRole     = "role"
	ContextAuthType = "authType"
	ContextAPIKeyID = "apiKeyId"
)

type Claims struct {
	UserID string `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type APIKeyVerifier interface {
	VerifyAPIKey(ctx context.Context, rawKey string) (*model.AuthContext, error)
}

func Sign(secret, userID, role string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	return t.SignedString([]byte(secret))
}

func RequireAuth(secret string, verifier APIKeyVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authHeader := c.GetHeader("Authorization"); strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err == nil && tok.Valid {
				if claims, ok := tok.Claims.(*Claims); ok && claims.UserID != "" {
					role := claims.Role
					if role == "" {
						role = model.RoleAdmin
					}
					setAuthContext(c, model.AuthContext{UserID: claims.UserID, Role: role, AuthType: "jwt"})
					c.Next()
					return
				}
			}
		}

		if verifier != nil {
			if rawKey := strings.TrimSpace(c.GetHeader("X-API-Key")); rawKey != "" {
				authCtx, err := verifier.VerifyAPIKey(c, rawKey)
				if err == nil && authCtx != nil && authCtx.UserID != "" {
					setAuthContext(c, *authCtx)
					c.Next()
					return
				}
			}
		}

		util.Err(c, http.StatusUnauthorized, "unauthorized")
		c.Abort()
	}
}

func RequireAPIKeyAuth(verifier APIKeyVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		if verifier != nil {
			if rawKey := strings.TrimSpace(c.GetHeader("X-API-Key")); rawKey != "" {
				authCtx, err := verifier.VerifyAPIKey(c, rawKey)
				if err == nil && authCtx != nil && authCtx.UserID != "" {
					setAuthContext(c, *authCtx)
					c.Next()
					return
				}
			}
		}
		util.Err(c, http.StatusUnauthorized, "unauthorized")
		c.Abort()
	}
}

func setAuthContext(c *gin.Context, auth model.AuthContext) {
	c.Set(ContextUserID, auth.UserID)
	c.Set(ContextRole, auth.Role)
	if auth.AuthType != "" {
		c.Set(ContextAuthType, auth.AuthType)
	}
	if auth.APIKeyID != "" {
		c.Set(ContextAPIKeyID, auth.APIKeyID)
	}
}

func GetAuthContext(c *gin.Context) model.AuthContext {
	return model.AuthContext{
		UserID:   c.GetString(ContextUserID),
		Role:     c.GetString(ContextRole),
		AuthType: c.GetString(ContextAuthType),
		APIKeyID: c.GetString(ContextAPIKeyID),
	}
}
