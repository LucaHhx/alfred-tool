package service

import (
	"alfred-tool/ui"
	"fmt"

	"github.com/spf13/cobra"
)

var serviceAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加新服务",
	Long:  `添加一个新的服务到指定服务器`,
	Run: func(cmd *cobra.Command, args []string) {
		addService()
	},
}

func addService() {
	err := ui.ShowAddServiceDialog()
	if err != nil {
		fmt.Printf("打开添加服务对话框失败: %v\n", err)
		return
	}
}
