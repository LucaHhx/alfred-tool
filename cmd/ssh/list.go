package ssh

import (
	"encoding/json"
	"fmt"

	"alfred-tool/models"
	"alfred-tool/services"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有SSH连接",
	Long:  `列出所有已保存的SSH连接配置。`,
	Run: func(cmd *cobra.Command, args []string) {
		connections, err := services.ListAllConnections()
		if err != nil {
			fmt.Printf("获取连接列表失败: %v\n", err)
			return
		}
		displayConnections(connections)
	},
}

func displayConnections(connections []models.SSHConnection) {
	alfredData := models.AlfredData{

		Items: lo.Map(connections, func(item models.SSHConnection, index int) models.AlfredItem {
			return models.AlfredItem{
				Uid:       item.Name,
				Title:     fmt.Sprintf("%s (%s)", item.Name, item.Description),
				Subtitle:  "\U00100A80 复制ip \U00100195 复制内网IP \U0010094C 连接服务器(%s)",
				Arg:       item.GetArg(),
				Variables: item.GetVariables(),
			}
		}),
	}
	marshal, err := json.Marshal(alfredData)
	if err != nil {
		fmt.Printf("JSON序列化失败: %v\n", err)
		return
	}
	fmt.Println(string(marshal))
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}
