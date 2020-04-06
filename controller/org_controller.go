package controller

import (
	"fmt"
	"github.com/cihub/seelog"
	"zoe/model"
	"zoe/mw"
	"strings"
)

type OrgControllerEngine struct {
}

func NewOrgControllerEngine() (*OrgControllerEngine, error) {
	return &OrgControllerEngine{}, nil
}

func (this *OrgControllerEngine) getOrgById(id int) (*model.Org, error) {
	var org model.Org
	sql := "select * from org where id = ? and is_deleted = 0"
	err := mw.DB.QueryRow(sql, id).Scan(&org.Id, &org.Name, &org.Visibility, &org.CurrentVersionId, &org.IsDeleted, &org.UpdatedAt, &org.CreateAt)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		seelog.Info(err.Error())
		return nil, err
	}
	return &org, nil
}

func (this *OrgControllerEngine) getOrgByName(name string) (*model.Org, error) {
	var org model.Org
	sql := "select * from org where name = ? and is_deleted = 0"
	err := mw.DB.QueryRow(sql, name).Scan(&org.Id, &org.Name, &org.Visibility, &org.CurrentVersionId, &org.IsDeleted, &org.UpdatedAt, &org.CreateAt)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		seelog.Info(err.Error())
		return nil, err
	}
	return &org, nil
}

func (this *OrgControllerEngine) QueryOrgById(id int) (*model.Org, error) {
	org, err := this.getOrgById(id)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (this *OrgControllerEngine) IsExistingOrgByName(name string) (bool, error) {
	org, err := this.getOrgByName(name)
	if err != nil {
		return false, err
	}
	if org != nil {
		return true, nil
	}
	return false, nil
}

func (this *OrgControllerEngine) IsExistingOrgById(id int) (bool, error) {
	org, err := this.getOrgById(id)
	if err != nil {
		return false, err
	}
	if org != nil {
		return true, nil
	}
	return false, nil
}

func (this *OrgControllerEngine) CreateOrg(name string, private bool) (*model.Org, error) {
	visibility := 0
	if private {
		visibility = 1
	}
	sql := "insert into org(name, visibility) values(?, ?)"
	r, err := mw.DB.Exec(sql, name, visibility)
	if err != nil {
		return nil, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}
	fmt.Println("insert success: ", id)
	org, err := this.getOrgById(int(id))
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (this *OrgControllerEngine) UpdateOrg(orgId int, private bool) error {
	visibility := 0
	if private {
		visibility = 1
	}
	sql := "update org set visibility = ? where id = ? and is_deleted = 0"
	_, err := mw.DB.Exec(sql, visibility, orgId)
	if err != nil {
		return err
	}
	return nil
}

func (this *OrgControllerEngine) DeleteOrg(orgId int) error {
	// todo 删除组织的所有project和item
	sql := "update org set is_deleted = 1 where id = ?"
	_, err := mw.DB.Exec(sql, orgId)
	if err != nil {
		return err
	}
	return nil
}
func (this *OrgControllerEngine) ListOrg(orgNames *[]string) (*[]model.Org, error) {
	var orgs []model.Org
	cnt := len(*orgNames)
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
	err := mw.DB.Select(&orgs, sql, params...)
	if err != nil {
		return nil, err
	}
	return &orgs, nil
}

func (this *OrgControllerEngine) ListAllProject(orgId int) (*[]model.Project, error) {
	var projects []model.Project
	sql := "select * from project where parent_id = ? and is_deleted = 0"
	err := mw.DB.Select(&projects, sql, orgId)
	if err != nil {
		return nil, err
	}
	return &projects, nil
}
