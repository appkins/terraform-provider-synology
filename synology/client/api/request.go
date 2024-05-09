package api

// Request defines a contract for all Request implementations.
type Request interface {
	init(apiName, apiMethod string)
}

type ApiRequest struct {
	Version   int    `form:"version" query:"version"`
	APIName   string `form:"api" query:"version"`
	APIMethod string `form:"method" query:"version"`
}

// init implements Request.
func (b *ApiRequest) init(apiName, apiMethod string) {
	b.Version = ApiVersions[apiName]
	b.APIName = apiName
	b.APIMethod = apiMethod
}

func NewRequest(apiName, apiMethod string) *ApiRequest {
	request := &ApiRequest{}
	request.init(apiName, apiMethod)
	return request
}
