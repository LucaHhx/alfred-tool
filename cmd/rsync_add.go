package cmd

import (
	"fmt"
	"alfred-tool/ui"

	"github.com/spf13/cobra"
)

var rsyncAddCmd = &cobra.Command{
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

func init() {
	rsyncCmd.AddCommand(rsyncAddCmd)
}