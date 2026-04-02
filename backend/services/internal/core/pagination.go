package core

// Meta is shared between layers (repo → service → handler)
type Meta struct {
	// SQL-style pagination
	Page       int   `json:"page,omitempty"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`

	// Cursor-style (DynamoDB, etc.)
	LastKey any  `json:"last_key,omitempty"`
	HasMore bool `json:"has_more,omitempty"`
}

// NewMeta creates pagination metadata safely (SQL-style)
func NewMeta(page, limit int, total int64) Meta {
	if limit <= 0 {
		limit = 10
	}

	if page <= 0 {
		page = 1
	}

	var totalPages int
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit)) // safe ceil division
	}

	return Meta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// NewCursorMeta creates metadata for cursor-based pagination (DynamoDB)
func NewCursorMeta(limit int, lastKey any, hasMore bool) Meta {
	if limit <= 0 {
		limit = 10
	}

	return Meta{
		Limit:   limit,
		LastKey: lastKey,
		HasMore: hasMore,
	}
}

// PaginatedResult is used in repository layer
type PaginatedResult[T any] struct {
	Data []T
	Meta Meta
}

// PaginatedResponse is returned to client
type PaginatedResponse[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}
