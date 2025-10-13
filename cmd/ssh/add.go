package ssh

import (
	"fmt"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加SSH连接",
	Long:  `添加新的SSH连接配置。`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ShowAddDialogV2(); err != nil {
			fmt.Printf("添加连接失败: %v\n", err)
			return
		}
	},
}
