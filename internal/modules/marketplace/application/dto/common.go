// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/application/dto/common.go
package dto

// ApiResponse is a common response structure
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse is a paginated response
type PaginatedResponse struct {
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
	HasMore    bool        `json:"has_more"`
	Data       interface{} `json:"data"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, total, page, perPage int) PaginatedResponse {
	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	return PaginatedResponse{
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
		HasMore:    page < totalPages,
		Data:       data,
	}
}
