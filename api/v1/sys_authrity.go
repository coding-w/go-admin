package v1

import (
	"github.com/gin-gonic/gin"
	"go-admin/global"
	"go-admin/model/common/request"
	"go-admin/model/common/response"
	"go-admin/model/system"
	"go-admin/model/vo"
	"go-admin/utils"
	"go.uber.org/zap"
)

type AuthorityApi struct{}

// CreateAuthority
// 创建角色
func (a *AuthorityApi) CreateAuthority(c *gin.Context) {
	var authority, authBack system.SysAuthority
	var err error
	if err = c.ShouldBindJSON(&authority); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err = utils.Verify(authority, utils.AuthorityVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if authBack, err = authorityService.CreateAuthority(authority); err != nil {
		global.GA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败"+err.Error(), c)
		return
	}
	err = casbinService.FreshCasbin()
	if err != nil {
		global.GA_LOG.Error("创建成功，权限刷新失败。", zap.Error(err))
		response.FailWithMessage("创建成功，权限刷新失败。"+err.Error(), c)
		return
	}
	response.OkWithDetailed(vo.SysAuthorityResponse{Authority: authBack}, "创建成功", c)

}

// CopyAuthority
// 复制角色
func (a *AuthorityApi) CopyAuthority(c *gin.Context) {
	var copyInfo vo.SysAuthorityCopyResponse
	err := c.ShouldBindJSON(&copyInfo)
	// 获取参数信息
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 检验参数
	err = utils.Verify(copyInfo, utils.OldAuthorityVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(copyInfo.Authority, utils.AuthorityVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	authBack, err := authorityService.CopyAuthority(copyInfo)
	if err != nil {
		global.GA_LOG.Error("拷贝失败!", zap.Error(err))
		response.FailWithMessage("拷贝失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(vo.SysAuthorityResponse{Authority: authBack}, "拷贝成功", c)
}

// DeleteAuthority
// 删除角色
func (a *AuthorityApi) DeleteAuthority(c *gin.Context) {
	var authority system.SysAuthority
	var err error
	if err = c.ShouldBindJSON(&authority); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err = utils.Verify(authority, utils.AuthorityVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 删除角色之前需要判断是否有用户正在使用此角色
	if err = authorityService.DeleteAuthority(&authority); err != nil {
		global.GA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
		return
	}
	_ = casbinService.FreshCasbin()
	response.OkWithMessage("删除成功", c)
}

// UpdateAuthority
// 更新角色信息
func (a *AuthorityApi) UpdateAuthority(c *gin.Context) {
	var auth system.SysAuthority
	err := c.ShouldBindJSON(&auth)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(auth, utils.AuthorityVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	authority, err := authorityService.UpdateAuthority(auth)
	if err != nil {
		global.GA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(vo.SysAuthorityResponse{Authority: authority}, "更新成功", c)
}

// GetAuthorityList
// 分页获取角色列表
func (a *AuthorityApi) GetAuthorityList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := authorityService.GetAuthorityInfoList(pageInfo)
	if err != nil {
		global.GA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// SetDataAuthority
// 设置角色资源权限
func (a *AuthorityApi) SetDataAuthority(c *gin.Context) {
	var auth system.SysAuthority
	err := c.ShouldBindJSON(&auth)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(auth, utils.AuthorityIdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = authorityService.SetDataAuthority(auth)
	if err != nil {
		global.GA_LOG.Error("设置失败!", zap.Error(err))
		response.FailWithMessage("设置失败"+err.Error(), c)
		return
	}
	response.OkWithMessage("设置成功", c)
}
