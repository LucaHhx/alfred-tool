package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"alfred-tool/services"
)

var useCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "增加SSH连接的使用次数",
	Long:  `根据名称增加SSH连接的使用次数。如果名称存在则使用次数+1，否则不进行任何操作。`,
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

func init() {
	rootCmd.AddCommand(useCmd)
}