package models

import (
	"fmt"
	"gorm.io/gorm"
)

type PasswordType string

const (
	PasswordTypePassword PasswordType = "password"
	PasswordTypeKeyPath  PasswordType = "keypath"
)

type SSHConnection struct {
	gorm.Model
	Name         string       `gorm:"uniqueIndex;not null" json:"name"`
	Address      string       `gorm:"not null" json:"address"`
	Port         int          `gorm:"default:22" json:"port"`
	Username     string       `gorm:"not null" json:"username"`
	PasswordType PasswordType `gorm:"not null" json:"password_type"`
	Password     string       `json:"password,omitempty"`
	KeyPath      string       `json:"key_path,omitempty"`
	LocalIP      string       `json:"local_ip"`
	Description  string       `json:"description"`
	UsageCount   int          `gorm:"default:0" json:"usage_count"`
}

func (s *SSHConnection) GetConnectionString() string {
	return s.Name + ":" + s.Username + "@" + s.Address
}

func (s *SSHConnection) GetArg() []string {
	return []string{s.Name, s.Address, fmt.Sprintf("%d", s.Port), s.LocalIP, s.Username, s.KeyPath}
}

func (s *SSHConnection) GetVariables() map[string]string {
	return map[string]string{
		"ssh_name":     s.Name,
		"ssh_address":  s.Address,
		"ssh_port":     fmt.Sprintf("%d", s.Port),
		"ssh_username": s.Username,
		"ssh_key_path": s.KeyPath,
		"ssh_local_ip": s.LocalIP,
		"ssh_desc":     s.Description,
	}
}
