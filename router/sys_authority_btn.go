package router

import (
	"github.com/gin-gonic/gin"
)

type AuthorityBtnRouter struct{}

func (a *AuthorityBtnRouter) InitAuthorityBtnRouter(Router *gin.RouterGroup) {
	authorityRouterWithoutRecord := Router.Group("authorityBtn")
	{
		authorityRouterWithoutRecord.POST("getAuthorityBtn", authorityBtnApi.GetAuthorityBtn)
		authorityRouterWithoutRecord.POST("setAuthorityBtn", authorityBtnApi.SetAuthorityBtn)
		authorityRouterWithoutRecord.POST("canRemoveAuthorityBtn", authorityBtnApi.CanRemoveAuthorityBtn)
	}

}
