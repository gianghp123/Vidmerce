package models

import "github.com/gianghp123/Vidmerce/backend/internal/core/enums"

type UserEntity struct {
	BaseItem
	// Creation timestamp
	CreatedAt string `dynamodbav:"createdAt"`
	// UserEntity email
	Email string         `dynamodbav:"email"`
	Role  enums.UserRole `dynamodbav:"role"`
}
