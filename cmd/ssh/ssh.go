package ssh

import (
	"github.com/spf13/cobra"
)

var SshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH连接管理",
	Long:  `SSH连接管理命令，支持添加、列出、搜索、更新、删除和使用SSH连接。`,
}

func init() {
	SshCmd.AddCommand(AddCmd)
	SshCmd.AddCommand(ListCmd)
	SshCmd.AddCommand(SearchCmd)
	SshCmd.AddCommand(UpdateCmd)
	SshCmd.AddCommand(DeleteCmd)
	SshCmd.AddCommand(UseCmd)
	SshCmd.AddCommand(SyncCmd)
}
