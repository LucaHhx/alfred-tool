package models

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type RsyncDirection string

const (
	RsyncDirectionUpload   RsyncDirection = "upload"   // 本地 -> 服务器
	RsyncDirectionDownload RsyncDirection = "download" // 服务器 -> 本地
)

type RsyncConfig struct {
	gorm.Model
	Name         string         `gorm:"uniqueIndex;not null" json:"name"`
	SSHName      string         `gorm:"not null" json:"ssh_name"` // 关联的SSH连接名称
	Direction    RsyncDirection `gorm:"not null" json:"direction"`
	LocalPath    string         `gorm:"not null" json:"local_path"`
	RemotePath   string         `gorm:"not null" json:"remote_path"`
	ExcludeRules string         `json:"exclude_rules"` // 排除规则，换行分隔
	Options      string         `json:"options"`       // 额外的rsync选项
	Description  string         `json:"description"`
	UsageCount   int            `gorm:"default:0" json:"usage_count"`

	// 常用rsync选项
	Verbose   bool `json:"verbose"`   // -v 详细输出
	Recursive bool `json:"recursive"` // -r 递归
	Archive   bool `json:"archive"`   // -a 归档模式
	Compress  bool `json:"compress"`  // -z 压缩
	Times     bool `json:"times"`     // -t 保持时间戳
	Progress  bool `json:"progress"`  // --progress 显示进度
	Delete    bool `json:"delete"`    // --delete 删除目标中多余的文件
	DryRun    bool `json:"dry_run"`   // --dry-run 预览模式
	Checksum  bool `json:"checksum"`  // -c 使用校验和
	Links     bool `json:"links"`     // -l 复制符号链接
	Perms     bool `json:"perms"`     // -p 保持权限
	Owner     bool `json:"owner"`     // -o 保持所有者
	Group     bool `json:"group"`     // -g 保持组
}

func (r *RsyncConfig) GetExcludeRulesSlice() []string {

	if r.ExcludeRules == "" {
		return []string{}
	}
	lines := strings.Split(r.ExcludeRules, "\n")
	var rules []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			rules = append(rules, line)
		}
	}
	return rules
}

func (r *RsyncConfig) GetOptionsSlice() []string {
	if r.Options == "" {
		return []string{}
	}
	return strings.Fields(r.Options)
}

func (r *RsyncConfig) BuildRsyncCommand(sshConnection *SSHConnection) []string {
	var cmd []string
	cmd = append(cmd, "rsync")

	// 构建选项字符串
	var options []string

	// 基本选项组合
	if r.Archive {
		options = append(options, "-a")
	} else {
		// 如果没有使用归档模式，检查单独的选项
		if r.Recursive {
			options = append(options, "-r")
		}
		if r.Links {
			options = append(options, "-l")
		}
		if r.Perms {
			options = append(options, "-p")
		}
		if r.Times {
			options = append(options, "-t")
		}
		if r.Owner {
			options = append(options, "-o")
		}
		if r.Group {
			options = append(options, "-g")
		}
	}

	// 其他选项
	if r.Verbose {
		options = append(options, "-v")
	}
	if r.Compress {
		options = append(options, "-z")
	}
	if r.Checksum {
		options = append(options, "-c")
	}

	// 组合短选项
	if len(options) > 0 {
		cmd = append(cmd, options...)
	}

	// 长选项
	if r.Progress {
		cmd = append(cmd, "--progress")
	}
	if r.Delete {
		cmd = append(cmd, "--delete")
	}
	if r.DryRun {
		cmd = append(cmd, "--dry-run")
	}

	// 添加用户自定义选项
	if r.Options != "" {
		cmd = append(cmd, r.GetOptionsSlice()...)
	}

	// 添加排除规则
	excludeRules := r.GetExcludeRulesSlice()
	for _, rule := range excludeRules {
		cmd = append(cmd, "--exclude", rule)
	}

	// SSH 选项
	sshOptions := fmt.Sprintf("-p %d", sshConnection.Port)
	if sshConnection.PasswordType == PasswordTypeKeyPath && sshConnection.KeyPath != "" {
		sshOptions += fmt.Sprintf(" -i %s", sshConnection.KeyPath)
	}
	cmd = append(cmd, "-e", fmt.Sprintf("ssh %s", sshOptions))

	// 源路径和目标路径
	if r.Direction == RsyncDirectionUpload {
		// 本地 -> 服务器
		cmd = append(cmd, r.LocalPath)
		cmd = append(cmd, fmt.Sprintf("%s@%s:%s", sshConnection.Username, sshConnection.Address, r.RemotePath))
	} else {
		// 服务器 -> 本地
		cmd = append(cmd, fmt.Sprintf("%s@%s:%s", sshConnection.Username, sshConnection.Address, r.RemotePath))
		cmd = append(cmd, r.LocalPath)
	}

	return cmd
}

func (r *RsyncConfig) GetDisplayInfo() string {
	direction := "↑"
	if r.Direction == RsyncDirectionDownload {
		direction = "↓"
	}
	return fmt.Sprintf("%s %s [%s] %s <-> %s", direction, r.Name, r.SSHName, r.LocalPath, r.RemotePath)
}

func (r *RsyncConfig) GetArg() []string {
	return []string{r.Name}
}

func (r *RsyncConfig) GetVariables() map[string]string {
	direction := "upload"
	if r.Direction == RsyncDirectionDownload {
		direction = "download"
	}
	
	return map[string]string{
		"rsync_name":         r.Name,
		"rsync_ssh_name":     r.SSHName,
		"rsync_direction":    direction,
		"rsync_local_path":   r.LocalPath,
		"rsync_remote_path":  r.RemotePath,
		"rsync_exclude":      r.ExcludeRules,
		"rsync_options":      r.Options,
		"rsync_description":  r.Description,
		"rsync_usage_count":  fmt.Sprintf("%d", r.UsageCount),
	}
}
