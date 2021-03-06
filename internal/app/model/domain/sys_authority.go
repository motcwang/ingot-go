package domain

import (
	"time"

	"github.com/ingot-cloud/ingot-go/pkg/framework/core/model/types"
)

// SysAuthority 权限
type SysAuthority struct {
	ID        types.ID `gorm:"primary_key;size:20"`
	TenantID  int
	PID       int64
	Name      string
	Code      string
	Path      string
	Method    string
	Status    string
	Remark    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// TableName 表名
func (*SysAuthority) TableName() string {
	return "sys_authority"
}
