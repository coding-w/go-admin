package v1

import (
	"github.com/gin-gonic/gin"
	"go-admin/global"
	"go-admin/model/common/request"
	"go-admin/model/common/response"
	"go-admin/model/dto"
	"go-admin/model/system"
	"go-admin/model/vo"
	"go-admin/utils"
	"go.uber.org/zap"
)

type SysApiApi struct{}

// CreateApi
// 创建基础api
func (saa *SysApiApi) CreateApi(c *gin.Context) {
	var api system.SysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(api, utils.ApiVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.CreateApi(api)
	if err != nil {
		global.GA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// SyncApi
// 同步API
func (saa *SysApiApi) SyncApi(c *gin.Context) {
	newApis, deleteApis, ignoreApis, err := apiService.SyncApi()
	if err != nil {
		global.GA_LOG.Error("同步失败!", zap.Error(err))
		response.FailWithMessage("同步失败", c)
		return
	}
	response.OkWithData(gin.H{
		"newApis":    newApis,
		"deleteApis": deleteApis,
		"ignoreApis": ignoreApis,
	}, c)
}

// GetApiGroups
// 获取API分组
func (saa *SysApiApi) GetApiGroups(c *gin.Context) {
	groups, apiGroupMap, err := apiService.GetApiGroups()
	if err != nil {
		global.GA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(gin.H{
		"groups":      groups,
		"apiGroupMap": apiGroupMap,
	}, c)
}

// IgnoreApi
// 忽略API
func (saa *SysApiApi) IgnoreApi(c *gin.Context) {
	var ignoreApi system.SysIgnoreApi
	err := c.ShouldBindJSON(&ignoreApi)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.IgnoreApi(ignoreApi)
	if err != nil {
		global.GA_LOG.Error("忽略失败!", zap.Error(err))
		response.FailWithMessage("忽略失败", c)
		return
	}
	response.Ok(c)
}

// EnterSyncApi
// 确认同步API
func (saa *SysApiApi) EnterSyncApi(c *gin.Context) {
	var syncApi vo.SysSyncApis
	err := c.ShouldBindJSON(&syncApi)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.EnterSyncApi(syncApi)
	if err != nil {
		global.GA_LOG.Error("忽略失败!", zap.Error(err))
		response.FailWithMessage("忽略失败", c)
		return
	}
	response.Ok(c)
}

// DeleteApi
// 删除api
func (saa *SysApiApi) DeleteApi(c *gin.Context) {
	var api system.SysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(api.GA_Model, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.DeleteApi(api)
	if err != nil {
		global.GA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetApiList
// 分页获取API列表
func (saa *SysApiApi) GetApiList(c *gin.Context) {
	var pageInfo dto.SearchApiParams
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := apiService.GetAPIInfoList(pageInfo.SysApi, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
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

// GetApiById
// 根据id获取api
func (saa *SysApiApi) GetApiById(c *gin.Context) {
	var idInfo request.GetById
	err := c.ShouldBindJSON(&idInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(idInfo, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	api, err := apiService.GetApiById(idInfo.ID)
	if err != nil {
		global.GA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(vo.SysAPIResponse{Api: api}, "获取成功", c)
}

// UpdateApi
// 修改基础api
func (saa *SysApiApi) UpdateApi(c *gin.Context) {
	var api system.SysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(api, utils.ApiVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.UpdateApi(api)
	if err != nil {
		global.GA_LOG.Error("修改失败!", zap.Error(err))
		response.FailWithMessage("修改失败", c)
		return
	}
	response.OkWithMessage("修改成功", c)
}

// GetAllApis
// 获取所有的Api 不分页
func (saa *SysApiApi) GetAllApis(c *gin.Context) {
	apis, err := apiService.GetAllApis()
	if err != nil {
		global.GA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(vo.SysAPIListResponse{Apis: apis}, "获取成功", c)
}

// DeleteApisByIds
// 删除选中Api
func (saa *SysApiApi) DeleteApisByIds(c *gin.Context) {
	var ids request.IdsReq
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = apiService.DeleteApisByIds(ids)
	if err != nil {
		global.GA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// FreshCasbin
// 刷新casbin缓存
func (saa *SysApiApi) FreshCasbin(c *gin.Context) {
	err := casbinService.FreshCasbin()
	if err != nil {
		global.GA_LOG.Error("刷新失败!", zap.Error(err))
		response.FailWithMessage("刷新失败", c)
		return
	}
	response.OkWithMessage("刷新成功", c)
}
