package req

type GetUsersQuery struct {
	Limit   int    `form:"limit,default=20" example:"20"`
	LastKey string `form:"last_key" example:"eyJ2b3RlSWQiOjF9"`
}

type GetUserByIDParam struct {
	ID string `uri:"id" binding:"required" example:"user_123456"`
}

type CreateUserReq struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
	Role  string `json:"role" binding:"required" example:"USER"`
}

type UpdateUserReq struct {
	Role *string `json:"role,omitempty" example:"ADMIN"`
}
