package service

import (
	"alfred-tool/ui"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var serviceUpdateCmd = &cobra.Command{
	Use:   "update [服务ID]",
	Short: "更新服务信息",
	Long:  `更新指定ID的服务信息`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			fmt.Println("无效的服务ID")
			return
		}
		updateService(uint(id))
	},
}

func updateService(id uint) {
	err := ui.ShowUpdateServiceDialog(id)
	if err != nil {
		fmt.Printf("打开修改服务对话框失败: %v\n", err)
		return
	}
}
