package db

import (
	"github.com/cihub/seelog"
	"zoe/basic"
	"zoe/model"
)

func IsExistingPrivilege(userHash, resName string, resId, resType, priType int) (bool, error) {
	var privilege model.Privilege
	sql := "select id from privilege where is_deleted = 0 and user_hash = ? and resource_name = ? and resource_id = ? and resource_type = ? and privilege_type = ?"
	err := DB.QueryRow(sql, userHash, resName, resId, resType, priType).Scan(&privilege.Id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		seelog.Info(err.Error())
		return false, err
	}
	return true, nil
}

func CreatePrivilege(userHash, resName string, resId, resType, userId, priType, resVisibility int) error {
	var sql = "insert into privilege(resource_id, resource_name, resource_type, resource_visibility, user_id, user_hash, privilege_type) values (?, ?, ?, ?, ?, ?, ?)"
	_, err := DB.Exec(sql, resId, resName, resType, resVisibility, userId, userHash, priType)
	if err != nil {
		seelog.Critical(err.Error())
		return err
	}
	return nil
}

func AddWithCheck(userHash, resName string, resId, resType, userId, priType, resVisibility int) error {
	flag, err := IsExistingPrivilege(userHash, resName, resId, resType, priType)
	if err != nil {
		seelog.Critical(err.Error())
		return err
	}
	if flag {
		return nil
	}
	err = CreatePrivilege(userHash, resName, resId, resType, userId, priType, resVisibility)
	return err
}

func ValidateForUserModifyOrg(userHash string, orgId int) (bool, error) {
	var privilege model.Privilege
	sql := "select id, privilege_type from privilege where user_hash = ? and resource_id = ? and resource_type = ? and is_deleted = 0"
	err := DB.QueryRow(sql, userHash, orgId, basic.Resource_Type_ORG).Scan(&privilege.Id, &privilege.PrivilegeType)
	if err != nil {
		return false, nil
	}
	if privilege.PrivilegeType == basic.Privilege_Type_MODIFIER {
		return true, nil
	} else {
		return false, nil
	}
}

func ValidateForUserViewOrg(userHash string, orgId int) (bool, error) {
	var privilege model.Privilege
	sql := "select id, privilege_type from privilege where user_hash = ? and resource_id = ? and resource_type = ? and is_deleted = 0"
	err := DB.QueryRow(sql, userHash, orgId, basic.Resource_Type_ORG).Scan(&privilege.Id, &privilege.PrivilegeType)
	if err != nil {
		return false, err
	}
	if privilege.PrivilegeType >= basic.Privilege_Type_VIEWER {
		return true, nil
	} else {
		return false, nil
	}
}

func DeletePrivilege(resId, resType int) error {
	sql := "update privilege set is_deleted = 1 where resource_id = ? and resource_type = ? and is_deleted = 0"
	_, err := DB.Exec(sql, resId, resType)
	if err != nil {
		return err
	}
	return nil
}

func ListPrivilege(userHash string) (*[]model.Privilege, error) {
	var privileges []model.Privilege
	sql := "select * from privilege where user_hash  = ? and resource_type in (?, ?, ?) and privilege_type in (?, ?) and is_deleted = 0"
	err := DB.Select(&privileges, sql, userHash, basic.Resource_Type_ORG, basic.Resource_Type_PROJECT, basic.Resource_Type_ITEM,
		basic.Privilege_Type_MODIFIER, basic.Privilege_Type_VIEWER)
	if err != nil {
		return nil, err
	}
	return &privileges, nil
}

func ListPrivilegeByResource(resId, resType int) (*[]model.Privilege, error) {
	var privileges []model.Privilege
	sql := "select * from privilege where resource_id = ? and resource_type = ? and is_deleted = 0"
	err := DB.Select(&privileges, sql, resId, resType)
	if err != nil {
		return nil, err
	}
	return &privileges, nil
}

func ListPrivilegeByPrefixResourceName(name, userHash string) (*[]model.Privilege, error) {
	var privileges []model.Privilege
	sql := "select * from privilege where user_hash = ? and resource_name like '?%' and is_deleted = 0"
	err := DB.Select(&privileges, sql, userHash, name)
	if err != nil {
		return nil, err
	}
	return &privileges, nil
}

func QueryPrivilegeByUserHash(userHash string, resId, resType int) (*model.Privilege, error) {
	var privileges []model.Privilege
	sql := "select * from privilege where user_hash = ? and resource_id = ? and resource_type = ? and is_deleted = 0"
	err := DB.Select(&privileges, sql, userHash, resId, resType)
	if err != nil {
		return nil, err
	}
	if len(privileges) > 0 {
		return &privileges[0], nil
	}
	return nil, nil
}

func UpdatePrivilegeByUserHash(userHash string, priType, resId, resType int) error {
	sql := "update privilege set privilege_type = ? where user_hash = ? and resource_id = ? and resource_type = ?"
	_, err := DB.Exec(sql, priType, userHash, resId, resType)
	if err != nil {
		return err
	}
	return nil
}

func DeletePrivilegeByUserHash(userHash string, resId, resType int) error {
	sql := "update privilege set is_deleted = 1 where user_hash = ? and resource_id = ? and resource_type = ? and is_deleted = 0"
	_, err := DB.Exec(sql, userHash, resId, resType)
	if err != nil {
		return err
	}
	return nil
}

func ValidateUserForProjectCreation(userHash string, orgId int) (bool, error) {
	flag, err := ValidateForUserViewOrg(userHash, orgId)
	if err != nil {
		return false, err
	}
	return flag, nil
}

func ValidateForUserModifyProject(userHash string, projectId, orgId int) (bool, error) {
	var privileges []model.Privilege
	sql := "select * from privilege where user_hash = ? and resource_id = ? and resource_type = ? and is_deleted = 0"
	err := DB.Select(&privileges, sql, userHash, projectId, basic.Resource_Type_PROJECT)
	if err != nil {
		return false, err
	}
	if len(privileges) > 0 {
		return true, nil
	}
	flag, err := ValidateForUserModifyOrg(userHash, orgId)
	if err != nil {
		return false, err
	}
	if flag {
		return true, nil
	} else {
		return false, nil
	}
}
