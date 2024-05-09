package filestation

import (
	"github.com/appkins/terraform-provider-synology/synology/client/api"
)

type FileStationInfoRequest struct {
	*api.ApiRequest
}

type FileStationInfoResponse struct {
	api.BaseResponse

	IsManager              bool
	SupportVirtualProtocol string
	Supportsharing         bool
	Hostname               string
}

var _ api.Request = (*FileStationInfoRequest)(nil)

func NewFileStationInfoRequest(version int) *FileStationInfoRequest {
	return &FileStationInfoRequest{
		ApiRequest: api.NewRequest("SYNO.FileStation.Info", "get"),
	}
}

func (r FileStationInfoResponse) ErrorSummaries() []api.ErrorSummary {
	return []api.ErrorSummary{commonErrors}
}
