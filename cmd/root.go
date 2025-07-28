package cmd

import (
	"fmt"
	"os"

	"alfred-tool/cmd/rsync"
	"alfred-tool/cmd/service"
	"alfred-tool/cmd/ssh"
	"alfred-tool/database"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "alfred-tool",
	Short: "Alfred效率工具箱",
	Long:  `Alfred效率工具箱 - 一个多功能的命令行工具，支持SSH连接管理、Rsync同步和服务管理。`,
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
	rootCmd.AddCommand(ssh.SshCmd)
	rootCmd.AddCommand(rsync.RsyncCmd)
	rootCmd.AddCommand(service.ServiceCmd)
}
