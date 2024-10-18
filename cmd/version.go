package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	ServeCmd.AddCommand(version)
}

var version = &cobra.Command{
	Use:   "version",
	Short: "打印程序版本",
	Long:  `打印程序版本，打印程序版本，打印程序版本`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("1.0.0")
	},
}
