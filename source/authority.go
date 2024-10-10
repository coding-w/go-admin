package source

import (
	"context"
	"github.com/pkg/errors"
	"go-admin/global"
	"go-admin/model/system"
	"go-admin/service/initdb"
	"go-admin/utils"
	"gorm.io/gorm"
)

const initOrderAuthority = initOrderCasbin + 1

type initAuthority struct{}

// auto run
func init() {
	initdb.RegisterInit(initOrderAuthority, &initAuthority{})
}

func (i *initAuthority) MigrateTable() error {
	return global.GA_DB.AutoMigrate(&system.SysAuthority{})
}

func (i *initAuthority) TableCreated() bool {
	return global.GA_DB.Migrator().HasTable(&system.SysAuthority{})
}

func (i *initAuthority) InitializerName() string {
	return system.SysAuthority{}.TableName()
}

func (i *initAuthority) InitializeData(ctx context.Context) (context.Context, error) {
	entities := []system.SysAuthority{
		{AuthorityId: 888, AuthorityName: "普通用户", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
		{AuthorityId: 9528, AuthorityName: "测试角色", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
		{AuthorityId: 8881, AuthorityName: "普通用户子角色", ParentId: utils.Pointer[uint](888), DefaultRouter: "dashboard"},
	}

	if err := global.GA_DB.Create(&entities).Error; err != nil {
		return ctx, errors.Wrapf(err, "%s表数据初始化失败!", system.SysAuthority{}.TableName())
	}

	// data authority
	if err := global.GA_DB.Model(&entities[0]).Association("DataAuthorityId").Replace(
		[]*system.SysAuthority{
			{AuthorityId: 888},
			{AuthorityId: 9528},
			{AuthorityId: 8881},
		}); err != nil {
		return ctx, errors.Wrapf(err, "%s表数据初始化失败!",
			global.GA_DB.Model(&entities[0]).Association("DataAuthorityId").Relationship.JoinTable.Name)
	}
	if err := global.GA_DB.Model(&entities[1]).Association("DataAuthorityId").Replace(
		[]*system.SysAuthority{
			{AuthorityId: 9528},
			{AuthorityId: 8881},
		}); err != nil {
		return ctx, errors.Wrapf(err, "%s表数据初始化失败!",
			global.GA_DB.Model(&entities[1]).Association("DataAuthorityId").Relationship.JoinTable.Name)
	}

	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initAuthority) DataInserted() bool {
	if errors.Is(global.GA_DB.Where("authority_id = ?", "8881").
		First(&system.SysAuthority{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
