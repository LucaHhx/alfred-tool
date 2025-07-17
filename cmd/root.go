package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sshd/database"
)

var rootCmd = &cobra.Command{
	Use:   "sshd",
	Short: "SSH连接管理工具",
	Long:  `一个用于管理SSH连接的命令行工具，支持添加、搜索、列出、修改和删除SSH连接配置。`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		database.InitDB()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(sshsyncCmd)
}