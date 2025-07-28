package ssh

import (
	"alfred-tool/ui"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加SSH连接",
	Long:  `添加新的SSH连接配置。`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.ShowAddDialog()
	},
}
