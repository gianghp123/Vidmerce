package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gianghp123/Vidmerce/backend/internal/configs"
	"github.com/gianghp123/Vidmerce/backend/internal/core"
	"github.com/gianghp123/Vidmerce/backend/internal/core/enums"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ClerkMetadata struct {
	Role enums.UserRole `json:"role"`
}

func customClaimsConstructor(ctx context.Context) any {
	return &ClerkMetadata{}
}

func withCustomClaims(params *clerkhttp.AuthorizationParams) error {
	params.VerifyParams.CustomClaimsConstructor = customClaimsConstructor
	return nil
}

func ClerkAuth() gin.HandlerFunc {
	logger := configs.GetLogger()
	clerkCfg := configs.GetClerkConfig()
	clerk.SetKey(clerkCfg.ClerkSecret)

	return func(c *gin.Context) {
		awsCfg := configs.LoadAWSConfig()
		if !awsCfg.IsLocal {
			c.Next()
			return
		}

		handler := clerkhttp.WithHeaderAuthorization(withCustomClaims)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if !ok {
				authHeader := r.Header.Get("Authorization")
				reason := fmt.Sprintf("missing or invalid clerk session, auth header: %s", authHeader)
				if authHeader == "" {
					reason = "missing clerk authorization header"
				}

				logger.Warn("Unauthorized request",
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
					zap.String("ip", c.ClientIP()),
					zap.String("reason", reason),
				)
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.Unauthorized())
				return
			}

			userID := claims.Subject
			role := enums.UserRoleUser // Default

			if customClaims, ok := claims.Custom.(*ClerkMetadata); ok && customClaims.Role != "" {
				role = customClaims.Role
			}

			c.Set(core.UserIDKey, userID)
			c.Set(core.RoleKey, role)

			ctx := context.WithValue(r.Context(), core.UserIDKey, userID)
			ctx = context.WithValue(ctx, core.RoleKey, role)

			// Replace the request context with our new one
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		}))

		handler.ServeHTTP(c.Writer, c.Request)

		if c.IsAborted() {
			return
		}
	}
}
