package vo

import "go-admin/model/system"

type ExaFileResponse struct {
	File system.ExaFileUploadAndDownload `json:"file"`
}
