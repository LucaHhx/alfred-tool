package cmd

import (
	"fmt"
	"alfred-tool/ui"

	"github.com/spf13/cobra"
)

var rsyncUpdateCmd = &cobra.Command{
	Use:   "update [配置名称]",
	Short: "修改rsync配置",
	Long:  `打开图形界面修改指定的rsync配置`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configName := args[0]
		
		err := ui.ShowUpdateRsyncDialog(configName)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
		}
	},
}

func init() {
	rsyncCmd.AddCommand(rsyncUpdateCmd)
}