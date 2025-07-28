package service

import (
	"alfred-tool/services"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有服务",
	Long:  `显示所有已添加的服务列表`,
	Run: func(cmd *cobra.Command, args []string) {
		listServices()
	},
}

func listServices() {
	serviceService := services.NewServiceService()

	serviceList, err := serviceService.GetAllServices()
	if err != nil {
		fmt.Printf("获取服务列表失败: %v\n", err)
		return
	}

	if len(serviceList) == 0 {
		fmt.Println("没有找到任何服务")
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
