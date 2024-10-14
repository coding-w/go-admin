package service

import (
	"errors"
	"go-admin/global"
	"go-admin/model/common/request"
	"go-admin/model/dto"
	"go-admin/model/system"
	"go-admin/model/vo"
	"gorm.io/gorm"
	"strconv"
)

var ErrRoleExistence = errors.New("存在相同角色id")

type AuthorityService struct{}

var authServiceApp = new(AuthorityService)

// CreateAuthority
// 创建角色
func (as *AuthorityService) CreateAuthority(auth system.SysAuthority) (authority system.SysAuthority, err error) {
	if err = global.GA_DB.Where("authority_id = ?", auth.AuthorityId).First(&system.SysAuthority{}).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return auth, ErrRoleExistence
	}
	e := global.GA_DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(&auth).Error; err != nil {
			return err
		}
		auth.SysBaseMenus = dto.DefaultMenu()
		if err = tx.Model(&auth).Association("SysBaseMenus").Replace(&auth.SysBaseMenus); err != nil {
			return err
		}
		casbinInfos := dto.DefaultCasbin()
		authorityId := strconv.Itoa(int(auth.AuthorityId))
		rules := [][]string{}
		for _, v := range casbinInfos {
			rules = append(rules, []string{authorityId, v.Path, v.Method})
		}
		return casbinServiceApp.AddPolicies(tx, rules)
	})
	return auth, e
}

// CopyAuthority
// 拷贝角色
func (as *AuthorityService) CopyAuthority(copyInfo vo.SysAuthorityCopyResponse) (authority system.SysAuthority, err error) {
	var authorityBox system.SysAuthority
	// 判断 角色id是否存在
	if !errors.Is(global.GA_DB.Where("authority_id = ? ", copyInfo.Authority.AuthorityId).First(&authorityBox).Error, gorm.ErrRecordNotFound) {
		return authority, ErrRoleExistence
	}
	copyInfo.Authority.Children = []system.SysAuthority{}
	// 获取拷贝角色的菜单
	menus, err := menuServiceApp.getMenuAuthority(&request.GetAuthorityId{AuthorityId: copyInfo.OldAuthorityId})
	if err != nil {
		return
	}
	var baseMenu []system.SysBaseMenu
	for _, v := range menus {
		intNum := v.MenuId
		v.SysBaseMenu.ID = uint(intNum)
		baseMenu = append(baseMenu, v.SysBaseMenu)
	}
	copyInfo.Authority.SysBaseMenus = baseMenu
	// 创建角色
	err = global.GA_DB.Create(&copyInfo.Authority).Error
	if err != nil {
		return
	}
	// 拷贝按钮
	var btns []system.SysAuthorityBtn
	err = global.GA_DB.Find(&btns, "authority_id = ?", copyInfo.OldAuthorityId).Error
	if err != nil {
		return
	}
	if len(btns) > 0 {
		for i := range btns {
			btns[i].AuthorityId = copyInfo.Authority.AuthorityId
		}
		err = global.GA_DB.Create(&btns).Error
		if err != nil {
			return
		}
	}
	paths := casbinServiceApp.GetPolicyPathByAuthorityId(copyInfo.OldAuthorityId)
	err = casbinServiceApp.UpdateCasbin(copyInfo.Authority.AuthorityId, paths)
	if err != nil {
		_ = as.DeleteAuthority(&copyInfo.Authority)
	}
	return copyInfo.Authority, err

}

func (as *AuthorityService) DeleteAuthority(auth *system.SysAuthority) error {
	if errors.Is(global.GA_DB.Debug().Preload("Users").First(&auth).Error, gorm.ErrRecordNotFound) {
		return errors.New("该角色不存在")
	}
	if len(auth.Users) != 0 {
		return errors.New("此角色有用户正在使用禁止删除")
	}
	if !errors.Is(global.GA_DB.Where("authority_id = ?", auth.AuthorityId).First(&system.SysUser{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色有用户正在使用禁止删除")
	}
	if !errors.Is(global.GA_DB.Where("parent_id = ?", auth.AuthorityId).First(&system.SysAuthority{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色存在子角色不允许删除")
	}
	return global.GA_DB.Transaction(func(tx *gorm.DB) error {
		var err error
		if err = tx.Preload("SysBaseMenus").Preload("DataAuthorityId").Where("authority_id = ?", auth.AuthorityId).First(auth).Unscoped().Delete(auth).Error; err != nil {
			return err
		}

		if len(auth.SysBaseMenus) > 0 {
			if err = tx.Model(auth).Association("SysBaseMenus").Delete(auth.SysBaseMenus); err != nil {
				return err
			}
			// err = db.Association("SysBaseMenus").Delete(&auth)
		}
		if len(auth.DataAuthorityId) > 0 {
			if err = tx.Model(auth).Association("DataAuthorityId").Delete(auth.DataAuthorityId); err != nil {
				return err
			}
		}

		if err = tx.Delete(&system.SysUserAuthority{}, "sys_authority_authority_id = ?", auth.AuthorityId).Error; err != nil {
			return err
		}
		if err = tx.Where("authority_id = ?", auth.AuthorityId).Delete(&[]system.SysAuthorityBtn{}).Error; err != nil {
			return err
		}

		authorityId := strconv.Itoa(int(auth.AuthorityId))

		if err = casbinServiceApp.RemoveFilteredPolicy(tx, authorityId); err != nil {
			return err
		}

		return nil
	})
}

func (as *AuthorityService) UpdateAuthority(auth system.SysAuthority) (authority system.SysAuthority, err error) {
	var oldAuthority system.SysAuthority
	err = global.GA_DB.Where("authority_id = ?", auth.AuthorityId).First(&oldAuthority).Error
	if err != nil {
		global.GA_LOG.Debug(err.Error())
		return system.SysAuthority{}, errors.New("查询角色数据失败")
	}
	err = global.GA_DB.Model(&oldAuthority).Updates(&auth).Error
	return auth, err
}

// GetAuthorityInfoList
// 分页获取数据
func (as *AuthorityService) GetAuthorityInfoList(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GA_DB.Model(&system.SysAuthority{})
	if err = db.Where("parent_id = ?", "0").Count(&total).Error; total == 0 || err != nil {
		return
	}
	var authority []system.SysAuthority
	err = db.Limit(limit).Offset(offset).Preload("DataAuthorityId").Where("parent_id = ?", "0").Find(&authority).Error
	for k := range authority {
		err = as.findChildrenAuthority(&authority[k])
	}
	return authority, total, err
}

// findChildrenAuthority
// 查询子角色
func (as *AuthorityService) findChildrenAuthority(authority *system.SysAuthority) (err error) {
	err = global.GA_DB.Preload("DataAuthorityId").Where("parent_id = ?", authority.AuthorityId).Find(&authority.Children).Error
	if len(authority.Children) > 0 {
		for k := range authority.Children {
			err = as.findChildrenAuthority(&authority.Children[k])
		}
	}
	return err
}

// SetDataAuthority
// 设置角色资源权限
func (as *AuthorityService) SetDataAuthority(auth system.SysAuthority) error {
	var s system.SysAuthority
	global.GA_DB.Preload("DataAuthorityId").First(&s, "authority_id = ?", auth.AuthorityId)
	err := global.GA_DB.Model(&s).Association("DataAuthorityId").Replace(&auth.DataAuthorityId)
	return err
}

// SetMenuAuthority
// 菜单与角色绑定
func (as *AuthorityService) SetMenuAuthority(auth *system.SysAuthority) error {
	var s system.SysAuthority
	global.GA_DB.Preload("SysBaseMenus").First(&s, "authority_id = ?", auth.AuthorityId)
	err := global.GA_DB.Model(&s).Association("SysBaseMenus").Replace(&auth.SysBaseMenus)
	return err
}
