package filestation

import (
	"github.com/appkins/terraform-provider-synology/synology/client/api"
	"github.com/appkins/terraform-provider-synology/synology/models"
)

type CreateShareRequest struct {
	api.ApiRequest

	SortBy     string   `query:"sort_by"`
	FileType   string   `query:"file_type"`
	CheckDir   bool     `query:"check_dir"`
	Additional []string `query:"additional" del:","`
}

type CreateShareResponse struct {
	api.BaseResponse

	Offset int `json:"offset"`

	Total int `json:"total"`
}

func (r CreateShareResponse) ErrorSummaries() []api.ErrorSummary {
	return []api.ErrorSummary{commonErrors}
}

type ListShareRequest struct {
	api.ApiRequest

	SortBy     string   `query:"sort_by"`
	FileType   string   `query:"file_type"`
	CheckDir   bool     `query:"check_dir"`
	Additional []string `query:"additional" del:","`
	GoToPath   string   `query:"goto_path"`
	FolderPath string   `query:"folder_path"`
}

type ListShareResponse struct {
	Offset int `json:"offset"`

	Shares []models.Share `json:"shares"`

	Total int `json:"total"`
}

var _ api.Request = (*ListShareRequest)(nil)

func NewListShareRequest(sortBy string, fileType string, checkDir bool, additional []string, goToPath string, folderPath string) *ListShareRequest {

	if additional == nil {
		additional = []string{"real_path", "owner", "time", "perm", "mount_point_type", "sync_share", "volume_status", "indexed", "hybrid_share", "worm_share"}
	}
	if sortBy == "" {
		sortBy = "name"
	}
	return &ListShareRequest{
		ApiRequest: api.ApiRequest{
			Version:   2,
			APIName:   "SYNO.FileStation.List",
			APIMethod: "list_share",
		},
		SortBy:     sortBy,
		FileType:   fileType,
		CheckDir:   checkDir,
		Additional: additional,
		GoToPath:   goToPath,
		FolderPath: folderPath,
	}
}

func (r ListShareRequest) ErrorSummaries() []api.ErrorSummary {
	return []api.ErrorSummary{commonErrors}
}
