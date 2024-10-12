package source

import (
	"context"
	"github.com/pkg/errors"
	"go-admin/global"
	"go-admin/model/system"
	"go-admin/service/initdb"
	"gorm.io/gorm"
)

const initOrderDict = initOrderCasbin + 1

type initDict struct{}

// auto run
func init() {
	initdb.RegisterInit(initOrderDict, &initDict{})
}

func (id *initDict) InitializerName() string {
	return system.SysDictionary{}.TableName()
}

func (id *initDict) MigrateTable() (err error) {
	return global.GA_DB.AutoMigrate(&system.SysDictionary{})
}

func (id *initDict) InitializeData(ctx context.Context) (next context.Context, err error) {
	True := true
	entities := []system.SysDictionary{
		{Name: "性别", Type: "gender", Status: &True, Desc: "性别字典"},
		{Name: "数据库int类型", Type: "int", Status: &True, Desc: "int类型对应的数据库类型"},
		{Name: "数据库时间日期类型", Type: "time.Time", Status: &True, Desc: "数据库时间日期类型"},
		{Name: "数据库浮点型", Type: "float64", Status: &True, Desc: "数据库浮点型"},
		{Name: "数据库字符串", Type: "string", Status: &True, Desc: "数据库字符串"},
		{Name: "数据库bool类型", Type: "bool", Status: &True, Desc: "数据库bool类型"},
	}

	if err = global.GA_DB.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, system.SysDictionary{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, id.InitializerName(), entities)
	return next, nil
}

func (id *initDict) TableCreated() (created bool) {
	return global.GA_DB.Migrator().HasTable(&system.SysDictionary{})
}

func (id *initDict) DataInserted() (inserted bool) {
	if errors.Is(global.GA_DB.Where("type = ?", "bool").First(&system.SysDictionary{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
