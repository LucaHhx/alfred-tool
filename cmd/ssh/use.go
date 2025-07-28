package ssh

import (
	"fmt"

	"alfred-tool/services"
	"github.com/spf13/cobra"
)

var UseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "使用SSH连接",
	Long:  `使用指定的SSH连接，并增加使用次数。`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		err := services.IncrementUsageCount(name)
		if err != nil {
			fmt.Printf("增加使用次数失败: %v\n", err)
			return
		}
		fmt.Printf("已增加连接 %s 的使用次数\n", name)
	},
}
