package model

import "time"

type User struct {
	Id         int       `db:"id"`
	Name       string    `db:"name"`
	UserHash   string    `db:"user_hash"`
	SecretHash int       `db:"secret_hash"`
	IsDeleted  int       `db:"is_deleted"`
	UpdatedAt  time.Time `db:"updated_at"`
	CreateAt   time.Time `db:"created_at"`
}
