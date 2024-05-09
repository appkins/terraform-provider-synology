package api

type LoginRequest struct {
	*ApiRequest

	Account  string `url:"account"`
	Password string `url:"passwd"`
	Session  string `url:"session,omitempty"`
}

type LoginResponse struct {
	BaseResponse
	DeviceID     string `json:"did,omitempty"`
	SessionID    string `json:"sid"`
	Token        string `json:"synotoken"`
	IsPortalPort bool   `json:"is_portal_port,omitempty"`
}
