package db

import (
	"database/sql"
	"fmt"
	"github.com/cihub/seelog"
	"strings"
	"zoe/model"
)

func getOrgById(conn *sql.Tx, id int) (*model.Org, error) {
	var org model.Org
	sql := "select * from org where id = ? and is_deleted = 0"
	err := conn.QueryRow(sql, id).Scan(&org.Id, &org.Name, &org.Visibility, &org.CurrentVersionId, &org.IsDeleted, &org.UpdatedAt, &org.CreateAt)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		seelog.Info(err.Error())
		return nil, err
	}
	return &org, nil
}

func getOrgByName(conn *sql.Tx, name string) (*model.Org, error) {
	var org model.Org
	sql := "select * from org where name = ? and is_deleted = 0"
	err := conn.QueryRow(sql, name).Scan(&org.Id, &org.Name, &org.Visibility, &org.CurrentVersionId, &org.IsDeleted, &org.UpdatedAt, &org.CreateAt)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		seelog.Info(err.Error())
		return nil, err
	}
	return &org, nil
}

func QueryOrgById(conn *sql.Tx, id int) (*model.Org, error) {
	org, err := getOrgById(conn, id)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func IsExistingOrgByName(conn *sql.Tx, name string) (bool, error) {
	org, err := getOrgByName(conn, name)
	if err != nil {
		return false, err
	}
	if org != nil {
		return true, nil
	}
	return false, nil
}

func IsExistingOrgById(conn *sql.Tx, id int) (bool, error) {
	org, err := getOrgById(conn, id)
	if err != nil {
		return false, err
	}
	if org != nil {
		return true, nil
	}
	return false, nil
}

func CreateOrg(conn *sql.Tx, name string, visibility int) (int, error) {
	sql := "insert into org(name, visibility) values(?, ?)"
	r, err := conn.Exec(sql, name, visibility)
	if err != nil {
		return 0, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func UpdateOrg(conn *sql.Tx, orgId int, private bool) error {
	visibility := 0
	if private {
		visibility = 1
	}
	sql := "update org set visibility = ? where id = ? and is_deleted = 0"
	_, err := conn.Exec(sql, visibility, orgId)
	if err != nil {
		return err
	}
	return nil
}

func DeleteOrg(conn *sql.Tx, orgId int) error {
	// todo 删除组织的所有project和item
	sql := "update org set is_deleted = 1 where id = ?"
	_, err := conn.Exec(sql, orgId)
	if err != nil {
		return err
	}
	return nil
}
func ListOrg(conn *sql.Tx, orgNames *[]string) (*[]model.Org, error) {
	var orgs []model.Org
	cnt := len(*orgNames)
	if cnt == 0 {
		return &orgs, nil
	}
	sqlItems := make([]string, cnt)
	for index := range sqlItems {
		sqlItems[index] = "?"
	}
	sqlItemsStr := strings.Join(sqlItems, ",")
	sql := fmt.Sprintf("select * from org where name in (%s) and is_deleted = 0", sqlItemsStr)
	params := make([]interface{}, cnt, cnt)
	for index := range params {
		params[index] = (*orgNames)[index]
	}
	var org model.Org
	rows, err := conn.Query(sql, params...)
	if err != nil {
		_ = seelog.Critical(err)
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&org.Id, &org.Name, &org.Visibility, &org.CurrentVersionId, &org.IsDeleted, &org.UpdatedAt, &org.CreateAt)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, org)
	}
	return &orgs, nil
}
