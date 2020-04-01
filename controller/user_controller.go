package controller

import (
	"fmt"
	"github.com/cihub/seelog"
	"http_guldan_server/model"
	"http_guldan_server/mw"
	"strings"
)

type UserControllerEngine struct {
}

func NewUserControllerEngine() (*UserControllerEngine, error) {
	return &UserControllerEngine{}, nil
}

func (this *UserControllerEngine) GetUserByUserHash(userHash string) (*model.User, error) {
	var user model.User
	sql := "select * from user where user_hash = ? and is_deleted = 0"
	err := mw.DB.QueryRow(sql, userHash).Scan(&user.Id, &user.Name, &user.UserHash, &user.SecretHash, &user.IsDeleted, &user.UpdatedAt, &user.CreateAt)
	if err != nil {
		seelog.Info(err.Error())
		return nil, err
	}
	return &user, nil
}

func (this *UserControllerEngine) GetUserByUserId(userId int) (*model.User, error) {
	var user model.User
	sql := "select * from user where id = ? and is_deleted = 0"
	err := mw.DB.QueryRow(sql, userId).Scan(&user.Id, &user.Name, &user.UserHash, &user.SecretHash, &user.IsDeleted, &user.UpdatedAt, &user.CreateAt)
	if err != nil {
		seelog.Info(err.Error())
		return nil, err
	}
	return &user, nil
}

func (this *UserControllerEngine) ListUserByIds(ids []int) (*[]model.User, error) {
	var users []model.User
	cnt := len(ids)
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
	err := mw.DB.Select(&users, sql, params...)
	if err != nil {
		return nil, err
	}
	return &users, nil
}
