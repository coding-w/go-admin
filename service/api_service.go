package service

import (
	"errors"
	"fmt"
	"go-admin/global"
	"go-admin/model/common/request"
	"go-admin/model/system"
	"go-admin/model/vo"
	"gorm.io/gorm"
	"strings"
)

type ApiService struct {
}

func (as *ApiService) CreateApi(api system.SysApi) error {
	if !errors.Is(global.GA_DB.Where("path = ? AND method = ?", api.Path, api.Method).
		First(&system.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同api")
	}
	return global.GA_DB.Create(&api).Error
}

func (as *ApiService) SyncApi() (newApis, deleteApis, ignoreApis []system.SysApi, err error) {
	newApis = make([]system.SysApi, 0)
	deleteApis = make([]system.SysApi, 0)
	ignoreApis = make([]system.SysApi, 0)
	var apis []system.SysApi
	// 获取 所有 api
	err = global.GA_DB.Find(&apis).Error
	if err != nil {
		return
	}
	var ignores []system.SysIgnoreApi
	// 获取 忽略 api
	err = global.GA_DB.Find(&ignores).Error
	if err != nil {
		return
	}
	for i := range ignores {
		ignoreApis = append(ignoreApis, system.SysApi{
			Path:        ignores[i].Path,
			Description: "",
			ApiGroup:    "",
			Method:      ignores[i].Method,
		})
	}
	var cacheApis []system.SysApi
	for i := range global.GA_ROUTERS {
		ignoresFlag := false
		for j := range ignores {
			if ignores[j].Path == global.GA_ROUTERS[i].Path && ignores[j].Method == global.GA_ROUTERS[i].Method {
				ignoresFlag = true
			}
		}
		if !ignoresFlag {
			cacheApis = append(cacheApis, system.SysApi{
				Path:   global.GA_ROUTERS[i].Path,
				Method: global.GA_ROUTERS[i].Method,
			})
		}
	}

	// 对比数据库中的api和内存中的api，如果数据库中的api不存在于内存中，则把api放入删除数组
	// 如果内存中的api不存在于数据库中，则把api放入新增数组
	for i := range cacheApis {
		var flag bool
		// 如果存在于内存不存在于api数组中
		for j := range apis {
			if cacheApis[i].Path == apis[j].Path && cacheApis[i].Method == apis[j].Method {
				flag = true
			}
		}
		if !flag {
			newApis = append(newApis, system.SysApi{
				Path:        cacheApis[i].Path,
				Description: "",
				ApiGroup:    "",
				Method:      cacheApis[i].Method,
			})
		}
	}
	for i := range apis {
		var flag bool
		// 如果存在于api数组不存在于内存
		for j := range cacheApis {
			if cacheApis[j].Path == apis[i].Path && cacheApis[j].Method == apis[i].Method {
				flag = true
			}
		}
		if !flag {
			deleteApis = append(deleteApis, apis[i])
		}
	}
	return
}

func (as *ApiService) GetApiGroups() (groups []string, groupApiMap map[string]string, err error) {
	var apis []system.SysApi
	err = global.GA_DB.Find(&apis).Error
	if err != nil {
		return
	}
	groupApiMap = make(map[string]string, 0)
	for i := range apis {
		pathArr := strings.Split(apis[i].Path, "/")
		newGroup := true
		for i2 := range groups {
			if groups[i2] == apis[i].ApiGroup {
				newGroup = false
			}
		}
		if newGroup {
			groups = append(groups, apis[i].ApiGroup)
		}
		groupApiMap[pathArr[1]] = apis[i].ApiGroup
	}
	return
}

func (as *ApiService) IgnoreApi(ignoreApi system.SysIgnoreApi) error {
	if ignoreApi.Flag {
		return global.GA_DB.Create(&ignoreApi).Error
	}
	return global.GA_DB.Unscoped().Delete(&ignoreApi, "path = ? AND method = ?", ignoreApi.Path, ignoreApi.Method).Error
}

func (as *ApiService) EnterSyncApi(syncApis vo.SysSyncApis) error {
	return global.GA_DB.Transaction(func(tx *gorm.DB) error {
		var txErr error
		if syncApis.NewApis != nil && len(syncApis.NewApis) > 0 {
			txErr = tx.Create(&syncApis.NewApis).Error
			if txErr != nil {
				return txErr
			}
		}
		for i := range syncApis.DeleteApis {
			casbinServiceApp.ClearCasbin(1, syncApis.DeleteApis[i].Path, syncApis.DeleteApis[i].Method)
			txErr = tx.Delete(&system.SysApi{}, "path = ? AND method = ?", syncApis.DeleteApis[i].Path, syncApis.DeleteApis[i].Method).Error
			if txErr != nil {
				return txErr
			}
		}
		return nil
	})
}

func (as *ApiService) DeleteApi(api system.SysApi) (err error) {
	var entity system.SysApi
	err = global.GA_DB.First(&entity, "id = ?", api.ID).Error // 根据id查询api记录
	if errors.Is(err, gorm.ErrRecordNotFound) {               // api记录不存在
		return err
	}
	err = global.GA_DB.Delete(&entity).Error
	if err != nil {
		return err
	}
	casbinServiceApp.ClearCasbin(1, entity.Path, entity.Method)
	return nil
}

func (as *ApiService) GetAPIInfoList(api system.SysApi, info request.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GA_DB.Model(&system.SysApi{})
	var apiList []system.SysApi

	if api.Path != "" {
		db = db.Where("path LIKE ?", "%"+api.Path+"%")
	}

	if api.Description != "" {
		db = db.Where("description LIKE ?", "%"+api.Description+"%")
	}

	if api.Method != "" {
		db = db.Where("method = ?", api.Method)
	}

	if api.ApiGroup != "" {
		db = db.Where("api_group = ?", api.ApiGroup)
	}

	err = db.Count(&total).Error

	if err != nil {
		return apiList, total, err
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	if order != "" {
		orderMap := make(map[string]bool, 5)
		orderMap["id"] = true
		orderMap["path"] = true
		orderMap["api_group"] = true
		orderMap["description"] = true
		orderMap["method"] = true
		if !orderMap[order] {
			err = fmt.Errorf("非法的排序字段: %v", order)
			return apiList, total, err
		}
		OrderStr = order
		if desc {
			OrderStr = order + " desc"
		}
	}
	err = db.Order(OrderStr).Find(&apiList).Error
	return apiList, total, err
}

// GetApiById
// 根据id获取api
func (as *ApiService) GetApiById(id int) (api system.SysApi, err error) {
	err = global.GA_DB.First(&api, "id = ?", id).Error
	return
}

func (as *ApiService) UpdateApi(api system.SysApi) (err error) {
	var oldA system.SysApi
	err = global.GA_DB.First(&oldA, "id = ?", api.ID).Error
	if oldA.Path != api.Path || oldA.Method != api.Method {
		var duplicateApi system.SysApi
		if ferr := global.GA_DB.First(&duplicateApi, "path = ? AND method = ?", api.Path, api.Method).Error; ferr != nil {
			if !errors.Is(ferr, gorm.ErrRecordNotFound) {
				return ferr
			}
		} else {
			if duplicateApi.ID != api.ID {
				return errors.New("存在相同api路径")
			}
		}

	}
	if err != nil {
		return err
	}

	err = casbinServiceApp.UpdateCasbinApi(oldA.Path, api.Path, oldA.Method, api.Method)
	if err != nil {
		return err
	}

	return global.GA_DB.Save(&api).Error
}

func (as *ApiService) GetAllApis() (apis []system.SysApi, err error) {
	err = global.GA_DB.Order("id desc").Find(&apis).Error
	return
}

func (as *ApiService) DeleteApisByIds(ids request.IdsReq) (err error) {
	return global.GA_DB.Transaction(func(tx *gorm.DB) error {
		var apis []system.SysApi
		err = tx.Find(&apis, "id in ?", ids.Ids).Error
		if err != nil {
			return err
		}
		err = tx.Delete(&[]system.SysApi{}, "id in ?", ids.Ids).Error
		if err != nil {
			return err
		}
		for _, sysApi := range apis {
			casbinServiceApp.ClearCasbin(1, sysApi.Path, sysApi.Method)
		}
		return err
	})
}
