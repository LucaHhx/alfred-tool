package rsync

import (
	"alfred-tool/ui"
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "添加rsync配置",
	Long:  `打开图形界面添加新的rsync配置`,
	Run: func(cmd *cobra.Command, args []string) {
		err := ui.ShowAddRsyncDialog()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
		}
	},
}
