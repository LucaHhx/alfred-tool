package cmd

import (
	"github.com/spf13/cobra"
	"sshd/ui"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "添加新的SSH连接",
	Long:  `通过图形界面添加新的SSH连接配置。`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.ShowAddDialog()
	},
}