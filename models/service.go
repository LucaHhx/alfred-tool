package models

import (
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	Name            string        `gorm:"not null" json:"name"`
	Description     string        `json:"description"`
	Details         string        `json:"details"`
	SSHConnectionID uint          `json:"ssh_connection_id"`
	SSHConnection   SSHConnection `gorm:"foreignKey:SSHConnectionID" json:"ssh_connection"`
	UsageCount      int           `gorm:"default:0" json:"usage_count"`
}

func (s *Service) GetDisplayName() string {
	return s.Name + "@" + s.Name
}
