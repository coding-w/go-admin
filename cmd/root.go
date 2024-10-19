package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-admin/global"
)

var (
	cfgFile string
	port    int
)

func init() {
	cobra.OnInitialize(initConfig)
	ServeCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file path (default is config.yaml)")
	ServeCmd.PersistentFlags().IntVarP(&port, "port", "p", 8888, "serve port")
}

var ServeCmd = &cobra.Command{
	Use:   "app",
	Short: "go-admin 是基于 Gin 和 Gorm 开发的管理系统，学习和实践 Gin 的使用。",
	Long:  `go-admin 是基于 Gin 和 Gorm 开发的管理系统，主要具备 JWT + Casbin 鉴权、动态路由、动态菜单以及文件上传下载等功能。该项目旨在学习和实践 Gin 的使用，源自于某开源项目的改进与扩展。`,
	Run: func(cmd *cobra.Command, args []string) {
		global.GA_CONFIG.System.Port = port
		startApp()
	},
}

func initConfig() {
	v := viper.New()
	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigFile("config.yaml")
	}
	if len(cfgFile) > 0 {
		v.SetConfigFile(cfgFile)
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
}
