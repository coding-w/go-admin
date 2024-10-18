package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-admin/global"
)

var (
	cfgFile string
)

func init() {
	cobra.OnInitialize(initConfig)
	ServeCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is config.yaml)")
	ServeCmd.PersistentFlags().StringP("author", "n", "wangx", "author info")
	ServeCmd.PersistentFlags().IntP("port", "p", 8888, "serve port")
}

var ServeCmd = &cobra.Command{
	Use:   "app",
	Short: "Serve the application Short Description",
	Long:  `Serve the application Long Long Long Long Long Long Description`,
	Run: func(cmd *cobra.Command, args []string) {
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
