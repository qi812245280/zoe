package model

import "time"

type Org struct {
	Id               int       `db:"id"`
	Name             string    `db:"name"`
	Visibility       int       `db:"visibility"`
	CurrentVersionId int       `db:"current_version_id"`
	IsDeleted        int       `db:"is_deleted"`
	UpdatedAt        time.Time `db:"updated_at"`
	CreateAt         time.Time `db:"created_at"`
}
