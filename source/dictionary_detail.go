package source

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go-admin/global"
	"go-admin/model/system"
	"go-admin/service/initdb"
	"log"
)

const initOrderDictDetail = initOrderDict + 1

type initDictDetail struct{}

// auto run
func init() {
	initdb.RegisterInit(initOrderDictDetail, &initDictDetail{})
}

func (i *initDictDetail) InitializerName() string {
	return system.SysDictionaryDetail{}.TableName()
}

func (i *initDictDetail) MigrateTable() (err error) {
	return global.GA_DB.AutoMigrate(&system.SysDictionaryDetail{})
}

func (i *initDictDetail) InitializeData(ctx context.Context) (next context.Context, err error) {
	dicts, ok := ctx.Value(system.SysDictionary{}.TableName()).([]system.SysDictionary)
	if !ok {
		return ctx, errors.Wrap(initdb.MissingDependentContextError,
			fmt.Sprintf("未找到 %s 表初始化数据", system.SysDictionary{}.TableName()))
	}
	True := true
	dicts[0].SysDictionaryDetails = []system.SysDictionaryDetail{
		{Label: "男", Value: "1", Status: &True, Sort: 1},
		{Label: "女", Value: "2", Status: &True, Sort: 2},
	}

	dicts[1].SysDictionaryDetails = []system.SysDictionaryDetail{
		{Label: "smallint", Value: "1", Status: &True, Extend: "mysql", Sort: 1},
		{Label: "mediumint", Value: "2", Status: &True, Extend: "mysql", Sort: 2},
		{Label: "int", Value: "3", Status: &True, Extend: "mysql", Sort: 3},
		{Label: "bigint", Value: "4", Status: &True, Extend: "mysql", Sort: 4},
		{Label: "int2", Value: "5", Status: &True, Extend: "pgsql", Sort: 5},
		{Label: "int4", Value: "6", Status: &True, Extend: "pgsql", Sort: 6},
		{Label: "int6", Value: "7", Status: &True, Extend: "pgsql", Sort: 7},
		{Label: "int8", Value: "8", Status: &True, Extend: "pgsql", Sort: 8},
	}

	dicts[2].SysDictionaryDetails = []system.SysDictionaryDetail{
		{Label: "date", Status: &True},
		{Label: "time", Value: "1", Status: &True, Extend: "mysql", Sort: 1},
		{Label: "year", Value: "2", Status: &True, Extend: "mysql", Sort: 2},
		{Label: "datetime", Value: "3", Status: &True, Extend: "mysql", Sort: 3},
		{Label: "timestamp", Value: "5", Status: &True, Extend: "mysql", Sort: 5},
		{Label: "timestamptz", Value: "6", Status: &True, Extend: "pgsql", Sort: 5},
	}
	dicts[3].SysDictionaryDetails = []system.SysDictionaryDetail{
		{Label: "float", Status: &True},
		{Label: "double", Value: "1", Status: &True, Extend: "mysql", Sort: 1},
		{Label: "decimal", Value: "2", Status: &True, Extend: "mysql", Sort: 2},
		{Label: "numeric", Value: "3", Status: &True, Extend: "pgsql", Sort: 3},
		{Label: "smallserial", Value: "4", Status: &True, Extend: "pgsql", Sort: 4},
	}

	dicts[4].SysDictionaryDetails = []system.SysDictionaryDetail{
		{Label: "char", Status: &True},
		{Label: "varchar", Value: "1", Status: &True, Extend: "mysql", Sort: 1},
		{Label: "tinyblob", Value: "2", Status: &True, Extend: "mysql", Sort: 2},
		{Label: "tinytext", Value: "3", Status: &True, Extend: "mysql", Sort: 3},
		{Label: "text", Value: "4", Status: &True, Extend: "mysql", Sort: 4},
		{Label: "blob", Value: "5", Status: &True, Extend: "mysql", Sort: 5},
		{Label: "mediumblob", Value: "6", Status: &True, Extend: "mysql", Sort: 6},
		{Label: "mediumtext", Value: "7", Status: &True, Extend: "mysql", Sort: 7},
		{Label: "longblob", Value: "8", Status: &True, Extend: "mysql", Sort: 8},
		{Label: "longtext", Value: "9", Status: &True, Extend: "mysql", Sort: 9},
	}

	dicts[5].SysDictionaryDetails = []system.SysDictionaryDetail{
		{Label: "tinyint", Value: "1", Extend: "mysql", Status: &True},
		{Label: "bool", Value: "2", Extend: "pgsql", Status: &True},
	}
	for _, dict := range dicts {
		if err := global.GA_DB.Model(&dict).Association("SysDictionaryDetails").
			Replace(dict.SysDictionaryDetails); err != nil {
			return ctx, errors.Wrap(err, system.SysDictionaryDetail{}.TableName()+"表数据初始化失败!")
		}
	}
	return ctx, nil
}

func (i *initDictDetail) TableCreated() (created bool) {
	return global.GA_DB.Migrator().HasTable(&system.SysDictionaryDetail{})
}

func (i *initDictDetail) DataInserted() (inserted bool) {
	var dict system.SysDictionary
	if err := global.GA_DB.Model(system.SysDictionary{}).Preload("SysDictionaryDetails").Where("name", "数据库bool类型").
		First(&dict).Error; err != nil {
		log.Println(dict)
		return false
	}
	return len(dict.SysDictionaryDetails) > 0 && dict.SysDictionaryDetails[0].Label == "tinyint"
}
