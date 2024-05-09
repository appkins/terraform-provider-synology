package filestation

import "github.com/appkins/terraform-provider-synology/synology/client/api"

type DeleteStartRequest struct {
	api.ApiRequest

	Paths            []string `form:"path" url:"path"`
	AccurateProgress bool     `form:"accurate_progress" url:"accurate_progress"`
}

type DeleteStartResponse struct {
	api.BaseResponse

	TaskID string `mapstructure:"taskid" json:"taskid"`
}

var _ api.Request = (*CreateFolderRequest)(nil)

func NewDeleteStartRequest(paths []string, accurateProgress bool) *DeleteStartRequest {
	return &DeleteStartRequest{
		ApiRequest: api.ApiRequest{
			Version:   2,
			APIName:   "SYNO.FileStation.Delete",
			APIMethod: "start",
		},
		Paths:            paths,
		AccurateProgress: accurateProgress,
	}
}

func (r DeleteStartRequest) ErrorSummaries() []api.ErrorSummary {
	return []api.ErrorSummary{commonErrors}
}

type DeleteStatusRequest struct {
	api.ApiRequest

	TaskID string `form:"taskid" url:"taskid"`
}

type DeleteStatusResponse struct {
	api.BaseResponse

	Finished       bool   `mapstructure:"finished" json:"finished"`
	FoundDirNum    int    `mapstructure:"found_dir_num" json:"found_dir_num"`
	FoundFileNum   int    `mapstructure:"found_file_num" json:"found_file_num"`
	FoundFileSize  int    `mapstructure:"found_file_size" json:"found_file_size"`
	HasDir         bool   `mapstructure:"has_dir" json:"has_dir"`
	Path           string `mapstructure:"path" json:"path"`
	ProcessedNum   int    `mapstructure:"processed_num" json:"processed_num"`
	ProcessingPath string `mapstructure:"processing_path" json:"processing_path"`
	Progress       int    `mapstructure:"progress" json:"progress"`
	Total          int    `mapstructure:"total" json:"total"`
}

func NewDeleteStatusRequest(taskID string) *DeleteStatusRequest {
	return &DeleteStatusRequest{
		ApiRequest: api.ApiRequest{
			Version:   1,
			APIName:   "SYNO.FileStation.Delete",
			APIMethod: "status",
		},
		TaskID: taskID,
	}
}

func (r DeleteStatusRequest) ErrorSummaries() []api.ErrorSummary {
	return []api.ErrorSummary{commonErrors}
}
