package db

import (
	"database/sql"
	"fmt"
	"strings"
	"zoe/model"
)

func GetUserByUserHash(conn *sql.Tx, userHash string) (*model.User, error) {
	var user model.User
	sql := "select * from user where user_hash = ? and is_deleted = 0"
	err := conn.QueryRow(sql, userHash).Scan(&user.Id, &user.Name, &user.UserHash, &user.SecretHash, &user.IsDeleted, &user.UpdatedAt, &user.CreateAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUserId(conn *sql.Tx, userId int) (*model.User, error) {
	var user model.User
	sql := "select * from user where id = ? and is_deleted = 0"
	err := conn.QueryRow(sql, userId).Scan(&user.Id, &user.Name, &user.UserHash, &user.SecretHash, &user.IsDeleted, &user.UpdatedAt, &user.CreateAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func ListUserByIds(conn *sql.Tx, ids []int) (*[]model.User, error) {
	var users []model.User
	cnt := len(ids)
	if cnt == 0 {
		return &users, nil
	}
	sqlItems := make([]string, cnt)
	for index := range sqlItems {
		sqlItems[index] = "?"
	}
	sqlItemsStr := strings.Join(sqlItems, ",")
	sql := fmt.Sprintf("select * from user where is_deleted = 0 and id in (%s)", sqlItemsStr)
	params := make([]interface{}, cnt, cnt)
	for index := range params {
		params[index] = (ids)[index]
	}
	var user model.User
	rows, err := conn.Query(sql, params...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.UserHash, &user.SecretHash, &user.IsDeleted, &user.UpdatedAt, &user.CreateAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return &users, nil
}
