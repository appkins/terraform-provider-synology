package filestation

import (
	"github.com/appkins/terraform-provider-synology/synology/client/api"
	"github.com/appkins/terraform-provider-synology/synology/client/util/form"
)

type UploadRequest struct {
	*api.ApiRequest

	Path          string     `form:"path" url:"path"`
	CreateParents bool       `form:"create_parents" url:"create_parents"`
	Overwrite     bool       `form:"overwrite" url:"overwrite"`
	File          *form.File `form:"file" kind:"file"`
}

type UploadResponse struct {
	api.BaseResponse
}

var _ api.Request = (*UploadRequest)(nil)

func NewUploadRequest(path string, file *form.File) *UploadRequest {
	return &UploadRequest{
		ApiRequest: api.NewRequest("SYNO.FileStation.Upload", "upload"),
		Path:       path,
		File:       file,
	}
}

func (r *UploadRequest) WithPath(value string) *UploadRequest {
	r.Path = value
	return r
}

func (r *UploadRequest) WithFile(file *form.File) *UploadRequest {
	r.File = file
	return r
}

func (r *UploadRequest) WithCreateParents(value bool) *UploadRequest {
	r.CreateParents = value
	return r
}

func (r *UploadRequest) WithOverwrite(value bool) *UploadRequest {
	r.Overwrite = value
	return r
}

func (r UploadResponse) ErrorSummaries() []api.ErrorSummary {
	return []api.ErrorSummary{
		{
			1100: "Failed to create a folder. More information in <errors> object.",
			1101: "The number of folders to the parent folder would exceed the system limitation.",
		},
		commonErrors,
	}
}
