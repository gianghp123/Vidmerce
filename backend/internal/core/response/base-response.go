package response

type BaseResponse[T any] struct {
	Data  T     `json:"data,omitempty"`
	Error any   `json:"error,omitempty"`
	Meta  *Meta `json:"meta,omitempty"` // optional
}

// Success creates a success response
func Success[T any](data T) BaseResponse[T] {
	return BaseResponse[T]{
		Data: data,
	}
}

// SuccessWithMeta (optional, useful for edge cases)
func SuccessWithMeta[T any](data T, meta Meta) BaseResponse[T] {
	return BaseResponse[T]{
		Data: data,
		Meta: &meta,
	}
}

// Fail creates an error response
func Fail(err any) BaseResponse[any] {
	return BaseResponse[any]{
		Error: err,
	}
}
