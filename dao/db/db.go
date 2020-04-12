package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"zoe/config"
)

var DB *sqlx.DB

func InitMysql(c *config.Config) error {
	database, err := sqlx.Open("mysql", c.Database.ConnectionString)
	if err != nil {
		fmt.Println("init mysql error: ", err)
		return err
	}
	DB = database
	return nil
}

func Destroy() {
	if DB != nil {
		DB.Close()
	}
}
