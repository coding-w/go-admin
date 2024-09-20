package core

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go-admin/global"
)

func Viper() *viper.Viper {
	var configPath string
	pflag.StringVarP(&configPath, "config", "c", "", "config file path")
	// 解析命令行参数
	pflag.Parse()
	v := viper.New()
	if len(configPath) > 0 {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigFile("config.yaml")
	}
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	if err := v.Unmarshal(&global.GA_CONFIG); err != nil {
		panic(err)
	}
	return v
}
