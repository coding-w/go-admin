package main

import (
	"fmt"
	"go-admin/cmd"
	"os"
)

func main() {
	if err := cmd.ServeCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
