package service

import (
	"alfred-tool/services"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var serviceViewCmd = &cobra.Command{
	Use:   "view [服务ID]",
	Short: "查看服务详情",
	Long:  `查看指定ID服务的详细信息`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			fmt.Println("无效的服务ID")
			return
		}
		viewService(uint(id))
	},
}

func viewService(id uint) {
	serviceService := services.NewServiceService()

	service, err := serviceService.GetServiceByID(id)
	if err != nil {
		fmt.Printf("获取服务信息失败: %v\n", err)
		return
	}

	fmt.Printf("# 服务详情\n\n")
	fmt.Printf("## 基本信息\n\n")
	fmt.Printf("| 字段 | 值 |\n")
	fmt.Printf("|------|----|\n")
	fmt.Printf("| ID | %d |\n", service.ID)
	fmt.Printf("| 服务名称 | **%s** |\n", service.Name)
	fmt.Printf("| 使用次数 | %d |\n", service.UsageCount)

	fmt.Printf("\n## 服务描述\n\n")
	if service.Description != "" {
		fmt.Printf("**简介:** %s\n\n", service.Description)
	}
	if service.Details != "" {
		fmt.Printf("**详情:**\n\n```\n%s\n```\n\n", service.Details)
	}

	if service.SSHConnectionID > 0 {
		fmt.Printf("## 关联SSH连接\n\n")
		fmt.Printf("| 字段 | 值 |\n")
		fmt.Printf("|------|----|\n")
		fmt.Printf("| 连接名称 | **%s** |\n", service.SSHConnection.Name)
		fmt.Printf("| 服务器地址 | `%s:%d` |\n", service.SSHConnection.Address, service.SSHConnection.Port)
		fmt.Printf("| 用户名 | `%s` |\n", service.SSHConnection.Username)
		fmt.Printf("| 连接描述 | %s |\n", service.SSHConnection.Description)
	}

	fmt.Printf("\n## 时间信息\n\n")
	fmt.Printf("| 字段 | 值 |\n")
	fmt.Printf("|------|----|\n")
	fmt.Printf("| 创建时间 | %s |\n", service.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("| 更新时间 | %s |\n", service.UpdatedAt.Format("2006-01-02 15:04:05"))
}
