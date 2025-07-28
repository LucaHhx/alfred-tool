package rsync

import (
	"alfred-tool/models"
	"alfred-tool/services"
	"encoding/json"
	"fmt"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有rsync配置",
	Long:  `显示所有已保存的rsync配置`,
	Run: func(cmd *cobra.Command, args []string) {
		configs, err := services.GetAllRsyncConfigs()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}
		displayRsyncConfigs(configs)
	},
}

func displayRsyncConfigs(configs []models.RsyncConfig) {
	alfredData := models.AlfredData{
		Items: lo.Map(configs, func(item models.RsyncConfig, index int) models.AlfredItem {
			direction := "↑"
			directionText := "上传"
			if item.Direction == models.RsyncDirectionDownload {
				direction = "↓"
				directionText = "下载"
			}

			title := fmt.Sprintf("%s %s [%s]", direction, item.Name, item.SSHName)
			subtitle := fmt.Sprintf("%s: %s ↔ %s", directionText, truncateString(item.LocalPath, 25), truncateString(item.RemotePath, 25))

			if item.Description != "" {
				subtitle += fmt.Sprintf(" - %s", truncateString(item.Description, 30))
			}

			return models.AlfredItem{
				Uid:       item.Name,
				Title:     title,
				Subtitle:  subtitle,
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
