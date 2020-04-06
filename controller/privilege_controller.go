package controller

import (
	"github.com/cihub/seelog"
	"zoe/basic"
	"zoe/model"
	"zoe/mw"
)

type PrivilegeControllerEngine struct {
}

func NewPrivilegeControllerEngine() (*PrivilegeControllerEngine, error) {
	return &PrivilegeControllerEngine{}, nil
}

func (this *PrivilegeControllerEngine) IsExistingPrivilege(userHash, resName string, resId, resType, priType int) (bool, error) {
	var privilege model.Privilege
	sql := "select id from privilege where is_deleted = 0 and user_hash = ? and resource_name = ? and resource_id = ? and resource_type = ? and privilege_type = ?"
	err := mw.DB.QueryRow(sql, userHash, resName, resId, resType, priType).Scan(&privilege.Id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		seelog.Info(err.Error())
		return false, err
	}
	return true, nil
}

func (this *PrivilegeControllerEngine) CreatePrivilege(userHash, resName string, resId, resType, userId, priType, resVisibility int) error {
	var sql = "insert into privilege(resource_id, resource_name, resource_type, resource_visibility, user_id, user_hash, privilege_type) values (?, ?, ?, ?, ?, ?, ?)"
	_, err := mw.DB.Exec(sql, resId, resName, resType, resVisibility, userId, userHash, priType)
	if err != nil {
		seelog.Critical(err.Error())
		return err
	}
	return nil
}

func (this *PrivilegeControllerEngine) AddWithCheck(userHash, resName string, resId, resType, userId, priType, resVisibility int) error {
	flag, err := this.IsExistingPrivilege(userHash, resName, resId, resType, priType)
	if err != nil {
		seelog.Critical(err.Error())
		return err
	}
	if flag {
		return nil
	}
	err = this.CreatePrivilege(userHash, resName, resId, resType, userId, priType, resVisibility)
	return err
}

func (this *PrivilegeControllerEngine) ValidateForUserModifyOrg(userHash string, orgId int) (bool, error) {
	var privilege model.Privilege
	sql := "select id, privilege_type from privilege where user_hash = ? and resource_id = ? and resource_type = ? and is_deleted = 0"
	err := mw.DB.QueryRow(sql, userHash, orgId, basic.Resource_Type_ORG).Scan(&privilege.Id, &privilege.PrivilegeType)
	if err != nil {
		return false, nil
	}
	if privilege.PrivilegeType == basic.Privilege_Type_MODIFIER {
		return true, nil
	} else {
		return false, nil
	}
}

func (this *PrivilegeControllerEngine) ValidateForUserViewOrg(userHash string, orgId int) (bool, error) {
	var privilege model.Privilege
	sql := "select id, privilege_type from privilege where user_hash = ? and resource_id = ? and resource_type = ? and is_deleted = 0"
	err := mw.DB.QueryRow(sql, userHash, orgId, basic.Resource_Type_ORG).Scan(&privilege.Id, &privilege.PrivilegeType)
	if err != nil {
		return false, err
	}
	if privilege.PrivilegeType >= basic.Privilege_Type_VIEWER {
		return true, nil
	} else {
		return false, nil
	}
}

func (this *PrivilegeControllerEngine) DeletePrivilege(resId, resType int) error {
	sql := "update privilege set is_deleted = 1 where resource_id = ? and resource_type = ? and is_deleted = 0"
	_, err := mw.DB.Exec(sql, resId, resType)
	if err != nil {
		return err
	}
	return nil
}

func (this *PrivilegeControllerEngine) ListPrivilege(userHash string) (*[]model.Privilege, error) {
	var privileges []model.Privilege
	sql := "select * from privilege where user_hash  = ? and resource_type in (?, ?, ?) and privilege_type in (?, ?) and is_deleted = 0"
	err := mw.DB.Select(&privileges, sql, userHash, basic.Resource_Type_ORG, basic.Resource_Type_PROJECT, basic.Resource_Type_ITEM,
		basic.Privilege_Type_MODIFIER, basic.Privilege_Type_VIEWER)
	if err != nil {
		return nil, err
	}
	return &privileges, nil
}

func (this *PrivilegeControllerEngine) ListPrivilegeByResource(resId, resType int) (*[]model.Privilege, error) {
	var privileges []model.Privilege
	sql := "select * from privilege where resource_id = ? and resource_type = ? and is_deleted = 0"
	err := mw.DB.Select(&privileges, sql, resId, resType)
	if err != nil {
		return nil, err
	}
	return &privileges, nil
}

func (this *PrivilegeControllerEngine) ListPrivilegeByPrefixResourceName(name, userHash string) (*[]model.Privilege, error) {
	var privileges []model.Privilege
	sql := "select * from privilege where user_hash = ? and resource_name like '?%' and is_deleted = 0"
	err := mw.DB.Select(&privileges, sql, userHash, name)
	if err != nil {
		return nil, err
	}
	return &privileges, nil
}

func (this *PrivilegeControllerEngine) QueryPrivilegeByUserHash(userHash string, resId, resType int) (*model.Privilege, error) {
	var privileges []model.Privilege
	sql := "select * from privilege where user_hash = ? and resource_id = ? and resource_type = ? and is_deleted = 0"
	err := mw.DB.Select(&privileges, sql, userHash, resId, resType)
	if err != nil {
		return nil, err
	}
	if len(privileges) > 0 {
		return &privileges[0], nil
	}
	return nil, nil
}

func (this *PrivilegeControllerEngine) UpdatePrivilegeByUserHash(userHash string, priType, resId, resType int) error {
	sql := "update privilege set privilege_type = ? where user_hash = ? and resource_id = ? and resource_type = ?"
	_, err := mw.DB.Exec(sql, priType, userHash, resId, resType)
	if err != nil {
		return err
	}
	return nil
}

func (this *PrivilegeControllerEngine) DeletePrivilegeByUserHash(userHash string, resId, resType int) error {
	sql := "update privilege set is_deleted = 1 where user_hash = ? and resource_id = ? and resource_type = ? and is_deleted = 0"
	_, err := mw.DB.Exec(sql, userHash, resId, resType)
	if err != nil {
		return err
	}
	return nil
}

func (this *PrivilegeControllerEngine) ValidateUserForProjectCreation(userHash string, orgId int) (bool, error) {
	flag, err := this.ValidateForUserViewOrg(userHash, orgId)
	if err != nil {
		return false, err
	}
	return flag, nil
}

func (this *PrivilegeControllerEngine) ValidateForUserModifyProject(userHash string, projectId, orgId int) (bool, error) {
	var privileges []model.Privilege
	sql := "select * from privilege where user_hash = ? and resource_id = ? and resource_type = ? and is_deleted = 0"
	err := mw.DB.Select(&privileges, sql, userHash, projectId, basic.Resource_Type_PROJECT)
	if err != nil {
		return false, err
	}
	if len(privileges) > 0 {
		return true, nil
	}
	flag, err := this.ValidateForUserModifyOrg(userHash, orgId)
	if err != nil {
		return false, err
	}
	if flag {
		return true, nil
	} else {
		return false, nil
	}
}
