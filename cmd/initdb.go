package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-admin/service"
)

func init() {
	ServeCmd.AddCommand(run)
}

var run = &cobra.Command{
	Use:   "init",
	Short: "初始化数据库",
	Long:  "初始化数据库需要完善配置信息",
	Run: func(cmd *cobra.Command, args []string) {
		err := service.ServiceGroup.InitDB()
		if err != nil {
			fmt.Println(err)
		}
	},
}
