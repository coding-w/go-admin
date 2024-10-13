package source

import (
	"context"
	"github.com/gofrs/uuid/v5"
	"github.com/pkg/errors"
	"go-admin/global"
	"go-admin/model/system"
	"go-admin/service/initdb"
	"go-admin/utils"
	"gorm.io/gorm"
)

const initOrderUser = initOrderAuthority + 1

type initUser struct{}

// auto run
func init() {
	initdb.RegisterInit(initOrderUser, &initUser{})
}

func (i *initUser) InitializerName() string {
	return system.SysUser{}.TableName()
}

func (i *initUser) MigrateTable() (err error) {
	return global.GA_DB.AutoMigrate(&system.SysUser{})
}

func (i *initUser) InitializeData(ctx context.Context) (next context.Context, err error) {
	ap := ctx.Value("adminPassword")
	apStr, ok := ap.(string)
	if !ok {
		apStr = "123456"
	}
	password := utils.BcryptHash(apStr)
	adminPassword := utils.BcryptHash(apStr)

	entities := []system.SysUser{
		{
			UUID:        uuid.Must(uuid.NewV4()),
			Username:    "admin",
			Password:    adminPassword,
			NickName:    "Mr.奇淼",
			HeaderImg:   "https://qmplusimg.henrongyi.top/gva_header.jpg",
			AuthorityId: 888,
			Phone:       "17611111111",
			Email:       "333333333@qq.com",
		},
		{
			UUID:        uuid.Must(uuid.NewV4()),
			Username:    "a303176530",
			Password:    password,
			NickName:    "用户1",
			HeaderImg:   "https:///qmplusimg.henrongyi.top/1572075907logo.png",
			AuthorityId: 9528,
			Phone:       "17611111111",
			Email:       "333333333@qq.com"},
	}
	if err = global.GA_DB.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, system.SysUser{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)
	authorityEntities, ok := ctx.Value(system.SysAuthority{}.TableName()).([]system.SysAuthority)
	if !ok {
		return next, errors.Wrap(initdb.MissingDependentContextError, "创建 [用户-权限] 关联失败, 未找到权限表初始化数据")
	}
	if err = global.GA_DB.Model(&entities[0]).Association("Authorities").Replace(authorityEntities); err != nil {
		return next, err
	}
	if err = global.GA_DB.Model(&entities[1]).Association("Authorities").Replace(authorityEntities[:1]); err != nil {
		return next, err
	}
	return next, nil
}

func (i *initUser) TableCreated() (created bool) {
	return global.GA_DB.Migrator().HasTable(&system.SysUser{})
}

func (i *initUser) DataInserted() (inserted bool) {
	var record system.SysUser
	if errors.Is(global.GA_DB.Where("username = ?", "a303176530").
		Preload("Authorities").First(&record).Error,
		gorm.ErrRecordNotFound) {
		return false
	}
	return len(record.Authorities) > 0 && record.Authorities[0].AuthorityId == 888
}
