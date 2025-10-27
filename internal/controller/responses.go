package controller

// CommonError represents error response
// @Description Common error response
type CommonError struct {
	Error string `json:"error" example:"error message"`
}

// CommonResponse represents success response
// @Description Common success response
type CommonResponse struct {
	Message string `json:"message" example:"success message"`
}
