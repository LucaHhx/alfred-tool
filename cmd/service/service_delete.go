package service

import (
	"alfred-tool/services"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var serviceDeleteCmd = &cobra.Command{
	Use:   "delete [服务ID]",
	Short: "删除服务",
	Long:  `删除指定ID的服务`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			fmt.Println("无效的服务ID")
			return
		}
		deleteService(uint(id))
	},
}

func deleteService(id uint) {
	serviceService := services.NewServiceService()

	service, err := serviceService.GetServiceByID(id)
	if err != nil {
		fmt.Printf("获取服务信息失败: %v\n", err)
		return
	}

	fmt.Printf("确认删除服务: %s? (y/N): ", service.Name)

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	if response != "y" && response != "yes" {
		fmt.Println("取消删除")
		return
	}

	if err := serviceService.DeleteService(id); err != nil {
		fmt.Printf("删除服务失败: %v\n", err)
		return
	}

	fmt.Println("服务删除成功!")
}
