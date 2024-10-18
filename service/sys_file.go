package service

import (
	"errors"
	"go-admin/global"
	"go-admin/model/common/request"
	"go-admin/model/system"
	"go-admin/utils/upload"
	"mime/multipart"
	"strings"
)

type FileUploadService struct{}

func (fs *FileUploadService) UploadFile(header *multipart.FileHeader, noSave string) (system.ExaFileUploadAndDownload, error) {
	oss := upload.NewOss()
	filePath, key, uploadErr := oss.UploadFile(header)
	if uploadErr != nil {
		panic(uploadErr)
	}
	s := strings.Split(header.Filename, ".")
	f := system.ExaFileUploadAndDownload{
		Url:  filePath,
		Name: header.Filename,
		Tag:  s[len(s)-1],
		Key:  key,
	}
	if noSave == "0" {
		return f, fs.Upload(f)
	}
	return f, nil
}

func (fs *FileUploadService) Upload(file system.ExaFileUploadAndDownload) error {
	return global.GA_DB.Create(&file).Error
}

func (fs *FileUploadService) EditFileName(file system.ExaFileUploadAndDownload) error {
	var fileFromDb system.ExaFileUploadAndDownload
	return global.GA_DB.Where("id = ?", file.ID).First(&fileFromDb).Update("name", file.Name).Error
}

func (fs *FileUploadService) DeleteFile(file system.ExaFileUploadAndDownload) (err error) {
	var fileFromDb system.ExaFileUploadAndDownload
	fileFromDb, err = fs.FindFile(file.ID)
	if err != nil {
		return
	}
	oss := upload.NewOss()
	if err = oss.DeleteFile(fileFromDb.Key); err != nil {
		return errors.New("文件删除失败")
	}
	err = global.GA_DB.Where("id = ?", file.ID).Unscoped().Delete(&file).Error
	return err
}

func (fs *FileUploadService) FindFile(id uint) (system.ExaFileUploadAndDownload, error) {
	var file system.ExaFileUploadAndDownload
	err := global.GA_DB.Where("id = ?", id).First(&file).Error
	return file, err
}

func (fs *FileUploadService) GetFileRecordInfoList(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	keyword := info.Keyword
	db := global.GA_DB.Model(&system.ExaFileUploadAndDownload{})
	var fileLists []system.ExaFileUploadAndDownload
	if len(keyword) > 0 {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&fileLists).Error
	return fileLists, total, err
}
