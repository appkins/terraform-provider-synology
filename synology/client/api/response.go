package api

// Response defines an interface for all responses from Synology API.
type Response interface {
	ErrorDescriber

	// GetError returns the latest error associated with response, if any.
	GetError() ApiError

	// SetError sets error object for the current response.
	SetError(ApiError)

	// Success reports whether the current request was successful.
	Success() bool
}

// ApiResponse is a concrete Response implementation.
// It is a generic struct with common to all Synology response fields.
type ApiResponse[TData any] struct {
	Success bool     `json:"success"`
	Data    TData    `json:"data,omitempty"`
	Error   ApiError `json:"error,omitempty"`
}

func NewApiResponse[TData any]() *ApiResponse[TData] {
	return &ApiResponse[TData]{}
}

type BaseResponse struct {
	ApiError
}

func (b *BaseResponse) SetError(e ApiError) {
	b.ApiError = e
}

func (b BaseResponse) Success() bool {
	return b.ApiError.Code == 0
}

func (b *BaseResponse) GetError() ApiError {
	return b.ApiError
}

func (b *BaseResponse) ErrorSummaries() []ErrorSummary {
	return []ErrorSummary{
		GlobalErrors,
	}
}
