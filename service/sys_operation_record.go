package service

import (
	"go-admin/global"
	"go-admin/model/system"
)

// OperationRecordService
// 创建记录
type OperationRecordService struct{}

func (operationRecordService *OperationRecordService) CreateSysOperationRecord(sysOperationRecord system.SysOperationRecord) (err error) {
	return global.GA_DB.Create(&sysOperationRecord).Error
}
