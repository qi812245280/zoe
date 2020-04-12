package db

import (
	"database/sql"
	"github.com/cihub/seelog"
	"zoe/basic"
	"zoe/model"
)

func queryPrivilege(conn *sql.Tx, sql string, args ...interface{}) (*[]model.Privilege, error) {
	var privileges []model.Privilege
	rows, err := conn.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	var privilege model.Privilege
	for rows.Next() {
		err = rows.Scan(&privilege.Id, &privilege.ResourceId, &privilege.ResourceName, &privilege.ResourceType, &privilege.ResourceVisibility,
			&privilege.UserId, &privilege.UserHash, &privilege.PrivilegeType, &privilege.IsDeleted, &privilege.UpdatedAt, &privilege.CreateAt)
		if err != nil {
			return nil, err
		}
		privileges = append(privileges, privilege)
	}
	return &privileges, nil
}

func IsExistingPrivilege(conn *sql.Tx, userHash, resName string, resId, resType, priType int) (bool, error) {
	var privilege model.Privilege
	sql := "select id from privilege where is_deleted = 0 and user_hash = ? and resource_name = ? and resource_id = ? and resource_type = ? and privilege_type = ?"
	err := conn.QueryRow(sql, userHash, resName, resId, resType, priType).Scan(&privilege.Id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		seelog.Info(err.Error())
		return false, err
	}
	return true, nil
}

func CreatePrivilege(conn *sql.Tx, userHash, resName string, resId, resType, userId, priType, resVisibility int) error {
	var sql = "insert into privilege(resource_id, resource_name, resource_type, resource_visibility, user_id, user_hash, privilege_type) values (?, ?, ?, ?, ?, ?, ?)"
	_, err := conn.Exec(sql, resId, resName, resType, resVisibility, userId, userHash, priType)
	if err != nil {
		return err
	}
	return nil
}

func AddWithCheck(conn *sql.Tx, userHash, resName string, resId, resType, userId, priType, resVisibility int) error {
	flag, err := IsExistingPrivilege(conn, userHash, resName, resId, resType, priType)
	if err != nil {
		return err
	}
	if flag {
		return nil
	}
	err = CreatePrivilege(conn, userHash, resName, resId, resType, userId, priType, resVisibility)
	return err
}

func ValidateForUserModifyOrg(conn *sql.Tx, userHash string, orgId int) (bool, error) {
	var privilege model.Privilege
	sql := "select id, privilege_type from privilege where user_hash = ? and resource_id = ? and resource_type = ? and is_deleted = 0"
	err := conn.QueryRow(sql, userHash, orgId, basic.Resource_Type_ORG).Scan(&privilege.Id, &privilege.PrivilegeType)
	if err != nil {
		return false, nil
	}
	if privilege.PrivilegeType == basic.Privilege_Type_MODIFIER {
		return true, nil
	} else {
		return false, nil
	}
}

func ValidateForUserViewOrg(conn *sql.Tx, userHash string, orgId int) (bool, error) {
	var privilege model.Privilege
	sql := "select id, privilege_type from privilege where user_hash = ? and resource_id = ? and resource_type = ? and is_deleted = 0"
	err := conn.QueryRow(sql, userHash, orgId, basic.Resource_Type_ORG).Scan(&privilege.Id, &privilege.PrivilegeType)
	if err != nil {
		return false, err
	}
	if privilege.PrivilegeType >= basic.Privilege_Type_VIEWER {
		return true, nil
	} else {
		return false, nil
	}
}

func DeletePrivilege(conn *sql.Tx, resId, resType int) error {
	sql := "update privilege set is_deleted = 1 where resource_id = ? and resource_type = ? and is_deleted = 0"
	_, err := conn.Exec(sql, resId, resType)
	if err != nil {
		return err
	}
	return nil
}

func ListPrivilege(conn *sql.Tx, userHash string) (*[]model.Privilege, error) {
	sql := "select * from privilege where user_hash  = ? and resource_type in (?, ?, ?) and privilege_type in (?, ?) and is_deleted = 0"
	privileges, err := queryPrivilege(conn, sql, userHash, basic.Resource_Type_ORG, basic.Resource_Type_PROJECT, basic.Resource_Type_ITEM,
		basic.Privilege_Type_MODIFIER, basic.Privilege_Type_VIEWER)
	if err != nil {
		return nil, err
	}
	return privileges, nil
}

func ListPrivilegeByResource(conn *sql.Tx, resId, resType int) (*[]model.Privilege, error) {
	sql := "select * from privilege where resource_id = ? and resource_type = ? and is_deleted = 0"
	privileges, err := queryPrivilege(conn, sql, resId, resType)
	if err != nil {
		return nil, err
	}
	return privileges, nil
}

func ListPrivilegeByPrefixResourceName(conn *sql.Tx, name, userHash string) (*[]model.Privilege, error) {
	sql := "select * from privilege where user_hash = ? and resource_name like '?%' and is_deleted = 0"
	privileges, err := queryPrivilege(conn, sql, userHash, name)
	if err != nil {
		return nil, err
	}
	return privileges, nil
}

func QueryPrivilegeByUserHash(conn *sql.Tx, userHash string, resId, resType int) (*model.Privilege, error) {
	sql := "select * from privilege where user_hash = ? and resource_id = ? and resource_type = ? and is_deleted = 0"
	privileges, err := queryPrivilege(conn, sql, userHash, resId, resType)
	if err != nil {
		return nil, err
	}
	if len(*privileges) > 0 {
		return &(*privileges)[0], nil
	}
	return nil, nil
}

func UpdatePrivilegeByUserHash(conn *sql.Tx, userHash string, priType, resId, resType int) error {
	sql := "update privilege set privilege_type = ? where user_hash = ? and resource_id = ? and resource_type = ?"
	_, err := conn.Exec(sql, priType, userHash, resId, resType)
	if err != nil {
		return err
	}
	return nil
}

func DeletePrivilegeByUserHash(conn *sql.Tx, userHash string, resId, resType int) error {
	sql := "update privilege set is_deleted = 1 where user_hash = ? and resource_id = ? and resource_type = ? and is_deleted = 0"
	_, err := conn.Exec(sql, userHash, resId, resType)
	if err != nil {
		return err
	}
	return nil
}

func ValidateUserForProjectCreation(conn *sql.Tx, userHash string, orgId int) (bool, error) {
	flag, err := ValidateForUserViewOrg(conn, userHash, orgId)
	if err != nil {
		return false, err
	}
	return flag, nil
}

func ValidateForUserModifyProject(conn *sql.Tx, userHash string, projectId, orgId int) (bool, error) {
	sql := "select * from privilege where user_hash = ? and resource_id = ? and resource_type = ? and is_deleted = 0"
	privileges, err := queryPrivilege(conn, sql, userHash, projectId, basic.Resource_Type_PROJECT)
	if err != nil {
		return false, err
	}
	if len(*privileges) > 0 {
		return true, nil
	}
	flag, err := ValidateForUserModifyOrg(conn, userHash, orgId)
	if err != nil {
		return false, err
	}
	if flag {
		return true, nil
	} else {
		return false, nil
	}
}
