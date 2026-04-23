package guards

import (
	"context"

	"github.com/gianghp123/Vidmerce/backend/internal/core"
	"github.com/gianghp123/Vidmerce/backend/internal/core/enums"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
)

type AuthContext struct {
	UserID string
	Role   enums.UserRole
}

func FromAuthContext(ctx context.Context) (*AuthContext, error) {

	uid, okUID := ctx.Value(core.UserIDKey).(string)
	role, okRole := ctx.Value(core.RoleKey).(enums.UserRole)

	if !okUID {
		return nil, response.Unauthorized("user identity not found in context")
	}

	if !okRole {
		role = enums.UserRoleUser
	}

	return &AuthContext{
		UserID: uid,
		Role:   role,
	}, nil
}

func GuardOwn(
	ctx context.Context,
	resourceOwnerId string,
) *response.AppError {
	auth, err := FromAuthContext(ctx)
	if err != nil {
		return response.Unauthorized("authentication required")
	}

	// 1. Admin Bypass
	if auth.Role == enums.UserRoleAdmin {
		return nil
	}

	if auth.UserID != resourceOwnerId {
		return response.Forbidden("you do not have permission to access this resource")
	}

	return nil
}
