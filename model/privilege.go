package model

import "time"

type Privilege struct {
	Id                 int       `db:"id"`
	ResourceId         int       `db:"resource_id"`
	ResourceName       string    `db:"resource_name"`
	ResourceType       int       `db:"resource_type"`
	ResourceVisibility int       `db:"resource_visibility"`
	UserId             int       `db:"user_id"`
	UserHash           string    `db:"user_hash"`
	PrivilegeType      int       `db:"privilege_type"`
	IsDeleted          int       `db:"is_deleted"`
	UpdatedAt          time.Time `db:"updated_at"`
	CreateAt           time.Time `db:"created_at"`
}
