package system

import "go-admin/global"

type SysApi struct {
	global.GA_Model
	Path        string `json:"path" gorm:"column:path;comment:api路径"`                 // api路径
	Description string `json:"description" gorm:"column:description;comment:api中文描述"` // api中文描述
	ApiGroup    string `json:"apiGroup" gorm:"column:api_group;comment:api组"`         // api组
	Method      string `json:"method" gorm:"column:method;default:POST;comment:方法"`   // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
}

func (SysApi) TableName() string {
	return "sys_apis"
}

type SysIgnoreApi struct {
	global.GA_Model
	Path   string `json:"path" gorm:"column:path;comment:api路径"`               // api路径
	Method string `json:"method" gorm:"column:method;default:POST;comment:方法"` // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
	Flag   bool   `json:"flag" gorm:"-"`                                       // 是否忽略
}

func (SysIgnoreApi) TableName() string {
	return "sys_ignore_apis"
}
