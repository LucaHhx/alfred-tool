package service

import (
	"alfred-tool/services"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var serviceSearchCmd = &cobra.Command{
	Use:   "search [关键词]",
	Short: "搜索服务",
	Long:  `根据关键词搜索服务名称、服务器名称、描述或类型`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		searchServices(args[0])
	},
}

func searchServices(keyword string) {
	serviceService := services.NewServiceService()

	serviceList, err := serviceService.SearchServices(keyword)
	if err != nil {
		fmt.Printf("搜索服务失败: %v\n", err)
		return
	}

	if len(serviceList) == 0 {
		fmt.Printf("没有找到包含 '%s' 的服务\n", keyword)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\t服务名称\t关联SSH连接\t描述")
	fmt.Fprintln(w, "---\t--------\t----------\t----")

	for _, service := range serviceList {
		sshConnection := "无"
		if service.SSHConnectionID > 0 && service.SSHConnection.Name != "" {
			sshConnection = service.SSHConnection.Name
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
			service.ID,
			service.Name,
			sshConnection,
			service.Description,
		)
	}

	w.Flush()
}
