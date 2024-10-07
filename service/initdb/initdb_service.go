package initdb

import (
	"context"
	"errors"
	"go-admin/global"
	"sort"
)

type InitDBService struct{}

var (
	initializers initSlice
	cache        map[string]*orderedInitializer
)

// InitDB 创建数据库并初始化 总入口
func (initDBService *InitDBService) InitDB() (err error) {
	ctx := context.TODO()
	if len(initializers) == 0 {
		return errors.New("无可用初始化过程，请检查初始化是否已执行完成")
	}
	sort.Sort(&initializers) // 保证有依赖的 initializer 排在后面执行
	// Note: 若 initializer 只有单一依赖，可以写为 B=A+1, C=A+1; 由于 BC 之间没有依赖关系，所以谁先谁后并不影响初始化
	// 若存在多个依赖，可以写为 C=A+B, D=A+B+C, E=A+1;
	// C必然>A|B，因此在AB之后执行，D必然>A|B|C，因此在ABC后执行，而E只依赖A，顺序与CD无关，因此E与CD哪个先执行并不影响
	var initHandler TypedDBInitHandler
	switch global.GA_CONFIG.System.DbType {
	case "mysql":
		initHandler = NewMysqlInitHandler()
		ctx = context.WithValue(ctx, "dbtype", "mysql")
	case "pgsql":
		initHandler = NewPgsqlInitHandler()
		ctx = context.WithValue(ctx, "dbtype", "pgsql")
	default:
		return errors.New("必须明确数据库类型，pgsql 或者 mysql")
	}
	ctx, err = initHandler.EnsureDB(ctx)
	if err != nil {
		return err
	}

	if err = initHandler.InitTables(ctx, initializers); err != nil {
		return err
	}
	if err = initHandler.InitData(ctx, initializers); err != nil {
		return err
	}

	initializers = initSlice{}
	cache = map[string]*orderedInitializer{}
	return nil
}
