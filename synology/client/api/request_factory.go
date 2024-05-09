package api

type RequestFactory interface {
	NewRequest(apiName, apiMethod string) Request
}
