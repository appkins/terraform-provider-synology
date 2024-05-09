package filestation

import (
	"github.com/appkins/terraform-provider-synology/synology/client/api"
)

type CreateFolderRequest struct {
	api.ApiRequest

	folderPaths []string `form:"folder_path" url:"folder_path"`
	names       []string `form:"name" url:"name"`
	forceParent bool     `form:"force_parent" url:"force_parent"`
}

type CreateFolderResponse struct {
	api.BaseResponse

	Folders []struct {
		Path  string
		Name  string
		IsDir bool
	}
}

var _ api.Request = (*CreateFolderRequest)(nil)

func NewCreateFolderRequest(version int) *CreateFolderRequest {
	return &CreateFolderRequest{
		ApiRequest: api.ApiRequest{
			Version:   version,
			APIName:   "SYNO.FileStation.CreateFolder",
			APIMethod: "create",
		},
	}
}

func (r *CreateFolderRequest) WithFolderPath(value string) *CreateFolderRequest {
	r.folderPaths = append(r.folderPaths, value)
	return r
}

func (r *CreateFolderRequest) WithName(value string) *CreateFolderRequest {
	r.names = append(r.names, value)
	return r
}

func (r *CreateFolderRequest) WithForceParent(value bool) *CreateFolderRequest {
	r.forceParent = value
	return r
}

func (r CreateFolderResponse) ErrorSummaries() []api.ErrorSummary {
	return []api.ErrorSummary{
		{
			1100: "Failed to create a folder. More information in <errors> object.",
			1101: "The number of folders to the parent folder would exceed the system limitation.",
		},
		commonErrors,
	}
}
