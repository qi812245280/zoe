package mw

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"zoe/config"
)

var DB *sqlx.DB

func InitMysql(c *config.Config) error {
	database, err := sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/guldandb?charset=utf8&parseTime=true")
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
