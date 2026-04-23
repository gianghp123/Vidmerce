package middlewares

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gianghp123/Vidmerce/backend/internal/configs"
	"github.com/gianghp123/Vidmerce/backend/internal/core"
	"github.com/gianghp123/Vidmerce/backend/internal/core/enums"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gin-gonic/gin"
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
