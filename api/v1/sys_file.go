package v1

import (
	"github.com/gin-gonic/gin"
	"go-admin/global"
	"go-admin/model/common/request"
	"go-admin/model/common/response"
	"go-admin/model/system"
	"go-admin/model/vo"
	"go.uber.org/zap"
)

type FileUploadApi struct{}

// UploadFile
// 上传文件示例
func (f *FileUploadApi) UploadFile(c *gin.Context) {
	var file system.ExaFileUploadAndDownload
	noSave := c.DefaultQuery("noSave", "0")
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GA_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	file, err = fileUploadService.UploadFile(header, noSave) // 文件上传后拿到文件路径
	if err != nil {
		global.GA_LOG.Error("修改数据库链接失败!", zap.Error(err))
		response.FailWithMessage("修改数据库链接失败", c)
		return
	}
	response.OkWithDetailed(vo.ExaFileResponse{File: file}, "上传成功", c)
}

// EditFileName 编辑文件名或者备注
func (f *FileUploadApi) EditFileName(c *gin.Context) {
	var file system.ExaFileUploadAndDownload
	err := c.ShouldBindJSON(&file)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = fileUploadService.EditFileName(file)
	if err != nil {
		global.GA_LOG.Error("编辑失败!", zap.Error(err))
		response.FailWithMessage("编辑失败", c)
		return
	}
	response.OkWithMessage("编辑成功", c)
}

// DeleteFile
// 删除文件
func (f *FileUploadApi) DeleteFile(c *gin.Context) {
	var file system.ExaFileUploadAndDownload
	err := c.ShouldBindJSON(&file)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := fileUploadService.DeleteFile(file); err != nil {
		global.GA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetFileList
// 分页文件列表
func (f *FileUploadApi) GetFileList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := fileUploadService.GetFileRecordInfoList(pageInfo)
	if err != nil {
		global.GA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
