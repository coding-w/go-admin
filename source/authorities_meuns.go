package source

import (
	"context"
	"github.com/pkg/errors"
	"go-admin/global"
	"go-admin/model/system"
	"go-admin/service/initdb"
)

const initOrderMenuAuthority = initOrderMenu + initOrderAuthority

type initMenuAuthority struct{}

// auto run
func init() {
	initdb.RegisterInit(initOrderMenuAuthority, &initMenuAuthority{})
}

func (i *initMenuAuthority) InitializerName() string {
	return "sys_authority_menus"
}

func (i *initMenuAuthority) MigrateTable() (err error) {
	return nil
}

func (i *initMenuAuthority) InitializeData(ctx context.Context) (next context.Context, err error) {
	authorities, ok := ctx.Value(system.SysAuthority{}.TableName()).([]system.SysAuthority)
	if !ok {
		return ctx, errors.Wrap(initdb.MissingDependentContextError, "创建 [菜单-权限] 关联失败, 未找到权限表初始化数据")
	}
	menus, ok := ctx.Value(system.SysMenu{}.TableName()).([]system.SysBaseMenu)
	if !ok {
		return next, errors.Wrap(errors.New(""), "创建 [菜单-权限] 关联失败, 未找到菜单表初始化数据")
	}
	next = ctx
	// 888
	if err = global.GA_DB.Model(&authorities[0]).Association("SysBaseMenus").Replace(menus); err != nil {
		return next, err
	}

	// 8881
	menu8881 := menus[:2]
	menu8881 = append(menu8881, menus[7])
	if err = global.GA_DB.Model(&authorities[1]).Association("SysBaseMenus").Replace(menu8881); err != nil {
		return next, err
	}

	// 9528
	if err = global.GA_DB.Model(&authorities[2]).Association("SysBaseMenus").Replace(menus[:11]); err != nil {
		return next, err
	}
	if err = global.GA_DB.Model(&authorities[2]).Association("SysBaseMenus").Append(menus[12:17]); err != nil {
		return next, err
	}
	return next, nil
}

func (i *initMenuAuthority) TableCreated() (created bool) {
	return false
}

func (i *initMenuAuthority) DataInserted() (inserted bool) {
	var auth system.SysAuthority
	ret := global.GA_DB.Model(auth).Where("authority_id = ?", 9528).First(&auth)
	if ret != nil {
		if ret.Error != nil {
			return false

		}
		return len(auth.SysBaseMenus) > 0
	}
	return false
}
