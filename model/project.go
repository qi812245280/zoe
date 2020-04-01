package model

import "time"

type Project struct {
	Id               int       `db:"id"`
	Name             string    `db:"name"`
	ParentId         string    `db:"parent_id"`
	Visibility       int       `db:"visibility"`
	CurrentVersionId int       `db:"current_version_id"`
	IsDeleted        int       `db:"is_deleted"`
	UpdatedAt        time.Time `db:"updated_at"`
	CreateAt         time.Time `db:"created_at"`
}
