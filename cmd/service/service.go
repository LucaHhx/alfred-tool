package service

import (
	"github.com/spf13/cobra"
)

var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "服务管理",
	Long:  `管理服务器上部署的服务，包括添加、查看、搜索、修改和删除服务信息`,
}

func init() {
	ServiceCmd.AddCommand(serviceAddCmd)
	ServiceCmd.AddCommand(serviceListCmd)
	ServiceCmd.AddCommand(serviceSearchCmd)
	ServiceCmd.AddCommand(serviceViewCmd)
	ServiceCmd.AddCommand(serviceUpdateCmd)
	ServiceCmd.AddCommand(serviceDeleteCmd)
}
