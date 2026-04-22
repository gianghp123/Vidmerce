package middlewares

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gianghp123/Vidmerce/backend/internal/configs"
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

			// Type assert to your custom metadata
			customClaims, ok := claims.Custom.(*ClerkMetadata)
			if ok {
				c.Set("role", customClaims.Role)
			}

			c.Set("user_id", claims.Subject)
			c.Next()
		}))

		handler.ServeHTTP(c.Writer, c.Request)

		if c.IsAborted() {
			return
		}
	}
}
