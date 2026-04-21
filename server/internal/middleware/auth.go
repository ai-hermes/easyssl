package middleware

import (
	"net/http"
	"strings"
	"time"

	"easyssl/server/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	AdminID string `json:"adminId"`
	jwt.RegisteredClaims
}

func Sign(secret, adminID string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		AdminID: adminID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	return t.SignedString([]byte(secret))
}

func RequireAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			util.Err(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")
		tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !tok.Valid {
			util.Err(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		claims, ok := tok.Claims.(*Claims)
		if !ok {
			util.Err(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		c.Set("adminId", claims.AdminID)
		c.Next()
	}
}
