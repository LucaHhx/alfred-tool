package cmd

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"sshd/models"
	"sshd/services"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "显示所有SSH连接",
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
				//Mods: map[models.ModName]models.AlfredMod{
				//	models.Mod_Cmd: models.NewAlfredMod(fmt.Sprintf("复制IP:%s", item.Address), item.Address),
				//	models.Mod_Alt: models.NewAlfredMod(fmt.Sprintf("复制内网IP:%s", item.LocalIP), item.LocalIP),
				//	models.Mod_Fn:  models.NewAlfredMod("连接服务器", item.GetArg()...),
				//},
			}
		}),
	}
	Echo_Success(alfredData)
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}
