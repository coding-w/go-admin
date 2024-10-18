package upload

import (
	"go-admin/global"
	"mime/multipart"
)

// OSS 对象存储接口
type OSS interface {
	UploadFile(file *multipart.FileHeader) (string, string, error)
	DeleteFile(key string) error
}

// NewOss OSS的实例化方法
func NewOss() OSS {
	switch global.GA_CONFIG.System.OssType {
	case "local":
		return &Local{}
	case "aliyun-oss":
		return &AliyunOSS{}
	default:
		return &Local{}
	}
}
