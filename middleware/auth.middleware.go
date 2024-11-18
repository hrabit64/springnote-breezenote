package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hrabit64/springnote-breezenote/auth"
	"github.com/hrabit64/springnote-breezenote/config"
	"github.com/hrabit64/springnote-breezenote/dto"
	"net/http"
)

type AuthMiddleware interface {
	UseAuthMiddleware() gin.HandlerFunc
}

type authMiddleware struct {
	AuthClient auth.AuthenticateClient
}

func NewAuthMiddleware(authClient auth.AuthenticateClient) AuthMiddleware {
	return &authMiddleware{AuthClient: authClient}
}

// UseAuthMiddleware 인증 미들웨어를 사용합니다.
func (a *authMiddleware) UseAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			errResponse := &dto.ErrorResponse{
				Status:   http.StatusUnauthorized,
				Title:    "Unauthorized",
				Detail:   "missing authorization header",
				Instance: c.Request.Host + c.Request.URL.Path,
			}
			c.JSON(http.StatusUnauthorized, errResponse)
			c.Abort()
			return
		}

		// Firebase 토큰 검증
		idToken := token[len("Bearer "):]
		decodedToken, err := a.AuthClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			errResponse := &dto.ErrorResponse{
				Status:   http.StatusUnauthorized,
				Title:    "Unauthorized",
				Detail:   "invalid token",
				Instance: c.Request.Host + c.Request.URL.Path,
			}
			c.JSON(http.StatusUnauthorized, errResponse)
			c.Abort()
			return
		}

		// uid 가 관리자인지 확인
		uid := decodedToken.UID
		if uid != config.RootConfig.AllowImageUploadUid {
			errResponse := &dto.ErrorResponse{
				Status:   http.StatusUnauthorized,
				Title:    "Unauthorized",
				Detail:   "you are not allowed to upload image",
				Instance: c.Request.Host + c.Request.URL.Path,
			}
			c.JSON(http.StatusUnauthorized, errResponse)
			c.Abort()
			return
		}

		c.Next()

	}
}
