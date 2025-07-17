package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"sshd/services"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "搜索SSH连接",
	Long:  `根据名称或地址搜索SSH连接。`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.Join(args, " ")
		connections, err := services.SearchConnections(query)
		if err != nil {
			fmt.Printf("搜索失败: %v\n", err)
			return
		}

		if len(connections) == 0 {
			fmt.Println("未找到匹配的连接")
			return
		}

		fmt.Printf("找到 %d 个匹配的连接:\n\n", len(connections))
		displayConnections(connections)
	},
}