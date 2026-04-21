package models

import "github.com/gianghp123/Vidmerce/backend/internal/core/enums"

type UserEntity struct {
	BaseItem
	CreatedAt string         `dynamodbav:"createdAt"`
	UpdatedAt string         `dynamodbav:"updatedAt,omitempty"`
	Email     string         `dynamodbav:"email"`
	Role      enums.UserRole `dynamodbav:"role"`
}
