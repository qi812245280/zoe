package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"zoe/basic"
	"zoe/dao/db"
	"zoe/model"
	"zoe/utils"
)

func CreateOrg(userHash string, req model.OrgCreateRequest) (gin.H, error) {
	conn, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	user, err := utils.GetUser(conn, userHash)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	if len(req.Name) >= basic.MAX_RESOURCE_NAME_LENGTH {
		_ = conn.Rollback()
		return nil, errors.New("组织名长度过长")
	}
	flag, err := db.IsExistingOrgByName(conn, req.Name)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	if flag {
		return nil, errors.New("组织已经存在")
	}
	visibility := 0
	if req.Private {
		visibility = 1
	}
	id, err := db.CreateOrg(conn, req.Name, visibility)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	err = db.AddWithCheck(conn, user.UserHash, req.Name, id,
		basic.Resource_Type_ORG, user.Id, basic.Privilege_Type_MODIFIER, visibility)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	err = conn.Commit()
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	var visibilityStr string
	if visibility == 0 {
		visibilityStr = "private"
	} else {
		visibilityStr = "public"
	}
	return gin.H{
		"code": 0,
		"msg":  "OK",
		"data": gin.H{
			"id":         id,
			"name":       req.Name,
			"visibility": visibilityStr,
		},
	}, nil
}

func UpdateOrg(userHash string, orgId int, req model.OrgUpdateRequest) (gin.H, error) {
	conn, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	user, err := utils.GetUser(conn, userHash)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	flag, err := db.IsExistingOrgById(conn, orgId)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	if !flag {
		return nil, errors.New("不存在的组织")
	}

	flag, err = db.ValidateForUserModifyOrg(conn, user.UserHash, orgId)
	if err != nil || !flag {
		_ = conn.Rollback()
		return nil, errors.New("用户无效权限修改该组织")
	}
	if err := db.UpdateOrg(conn, orgId, req.Private); err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	err = conn.Commit()
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	return gin.H{"code": 0, "msg": "OK"}, nil
}

func DeleteOrg(userHash string, orgId int) (gin.H, error) {
	conn, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	user, err := utils.GetUser(conn, userHash)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	flag, err := db.ValidateForUserModifyOrg(conn, user.UserHash, orgId)
	if err != nil || !flag {
		_ = conn.Rollback()
		return nil, errors.New("用户无权限修改该组织")
	}
	if err = db.DeleteOrg(conn, orgId); err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	if err = db.DeletePrivilege(conn, orgId, basic.Resource_Type_ORG); err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	err = conn.Commit()
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	return gin.H{
		"code": 0,
		"msg":  "OK",
	}, nil
}

func ListOrg(userHash string) (gin.H, error) {
	conn, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	user, err := utils.GetUser(conn, userHash)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	privileges, err := db.ListPrivilege(conn, user.UserHash)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	orgNameList := make([]string, len(*privileges))
	for index, privilege := range *privileges {
		orgNameList[index] = strings.Split(privilege.ResourceName, ".")[0]
	}
	orgs, err := db.ListOrg(conn, &orgNameList)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	data := make([]gin.H, len(*orgs))
	for index := range data {
		data[index] = gin.H{"id": (*orgs)[index].Id, "name": (*orgs)[index].Name}
	}
	return gin.H{
		"code": 0,
		"msg":  "OK",
		"data": data,
	}, nil
}

func SingleOrg(userHash string, orgId int) (gin.H, error) {
	conn, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	user, err := utils.GetUser(conn, userHash)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	flag, err := db.IsExistingOrgById(conn, orgId)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	if !flag {
		return nil, errors.New("不存在的组织")
	}
	org, err := db.QueryOrgById(conn, orgId)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	flag, err = db.ValidateForUserModifyOrg(conn, user.UserHash, orgId)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	if flag {
		projects, err := db.ListProjectByParentId(conn, orgId)
		if err != nil {
			_ = conn.Rollback()
			return nil, err
		}
		privileges, err := db.ListPrivilegeByResource(conn, orgId, basic.Resource_Type_ORG)
		if err != nil {
			_ = conn.Rollback()
			return nil, err
		}
		userIds := make([]int, len(*privileges))
		for index, privilege := range *privileges {
			userIds[index] = privilege.UserId
		}
		users, err := db.ListUserByIds(conn, userIds)
		if err != nil {
			_ = conn.Rollback()
			return nil, err
		}
		privilegeInfo := utils.GetPrivilegeUserInfo(privileges, users)
		projectInfo := utils.GetProjectInfo(projects)
		orgInfo := utils.GetOrgInfo(projectInfo, privilegeInfo, org)
		return gin.H{
			"code": 0,
			"msg":  "OK",
			"data": orgInfo,
		}, nil
	} else {
		flag, err = db.ValidateForUserViewOrg(conn, user.UserHash, orgId)
		if err != nil {
			_ = conn.Rollback()
			return nil, err
		}
		if flag {
			projects, err := db.ListProjectByParentId(conn, orgId)
			if err != nil {
				_ = conn.Rollback()
				return nil, err
			}
			projectInfo := utils.GetProjectInfo(projects)
			orgInfo := utils.GetOrgInfo(projectInfo, nil, org)
			return gin.H{
				"code": 0,
				"msg":  "OK",
				"data": orgInfo,
			}, nil
		} else {
			publicProjects, privateProject, err := db.ListProjectByVisibility(conn, orgId)
			if err != nil {
				_ = conn.Rollback()
				return nil, err
			}
			privileges, err := db.ListPrivilegeByPrefixResourceName(conn, org.Name+".", user.UserHash)
			if err != nil {
				_ = conn.Rollback()
				return nil, err
			}
			var projectNames []string
			var projects []model.Project
			for _, item := range *privileges {
				arr := strings.Split(item.ResourceName, ".")
				projectNames = append(projectNames, arr[1])
			}
			for _, item := range *privateProject {
				for _, name := range projectNames {
					if name == item.Name {
						projects = append(projects, item)
						break
					}
				}
			}
			if len(*publicProjects) > 0 {
				projects = append(projects, *publicProjects...)
			}
			projectInfo := utils.GetProjectInfo(&projects)
			orgInfo := utils.GetOrgInfo(projectInfo, nil, org)
			return gin.H{
				"code": 0,
				"msg":  "OK",
				"data": orgInfo,
			}, nil
		}
	}
}

func AuthorizeOrg(userHash string, orgId int, req model.AuthorizeOrgRequest) (gin.H, error) {
	conn, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	user, err := utils.GetUser(conn, userHash)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	flag, err := db.ValidateForUserModifyOrg(conn, user.UserHash, orgId)
	if err != nil || !flag {
		_ = conn.Rollback()
		return nil, errors.New("用户无权限修改该组织")
	}
	targetUser, err := db.GetUserByUserId(conn, req.UserId)
	if err != nil {
		_ = conn.Rollback()
		return nil, errors.New("目标用户不存在")
	}
	if targetUser.Id == user.Id {
		_ = conn.Rollback()
		return nil, errors.New("你不能为自己授权")
	}
	org, err := db.QueryOrgById(conn, orgId)
	if err != nil {
		_ = conn.Rollback()
		return nil, errors.New("不存在灯组织")
	}
	privilege, err := db.QueryPrivilegeByUserHash(conn, targetUser.UserHash, orgId, basic.Resource_Type_ORG)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	var priType int
	if req.Type == "modifier" {
		priType = basic.Privilege_Type_MODIFIER
	} else if req.Type == "viewer" {
		priType = basic.Privilege_Type_VIEWER
	} else if req.Type == "puller" {
		priType = basic.Privilege_Type_PULLER
	} else {
		return nil, errors.New("非法的授权类型")
	}
	if privilege != nil {
		err = db.UpdatePrivilegeByUserHash(conn, targetUser.UserHash, priType, orgId, basic.Resource_Type_ORG)
	} else {
		err = db.CreatePrivilege(conn, targetUser.UserHash, org.Name, orgId, basic.Resource_Type_ORG, targetUser.Id, priType, org.Visibility)
	}
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	err = conn.Commit()
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	// todo 删除拉取的缓存

	return gin.H{"code": 0, "msg": "OK"}, nil
}

func DeleteAuthorizeOrg(userHash string, orgId, userId int) (gin.H, error) {
	conn, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	user, err := utils.GetUser(conn, userHash)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	targetUser, err := db.GetUserByUserId(conn, userId)
	if err != nil {
		_ = conn.Rollback()
		return nil, errors.New("目标用户不存在")
	}
	flag, err := db.ValidateForUserModifyOrg(conn, user.UserHash, orgId)
	if err != nil || !flag {
		_ = conn.Rollback()
		return nil, errors.New("用户无权限修改该组织")
	}
	privilege, err := db.QueryPrivilegeByUserHash(conn, targetUser.UserHash, orgId, basic.Resource_Type_ORG)
	if err != nil || privilege == nil {
		_ = conn.Rollback()
		return nil, errors.New("目标用户无该组织的权限")
	}
	err = db.DeletePrivilegeByUserHash(conn, targetUser.UserHash, orgId, basic.Resource_Type_ORG)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	err = conn.Commit()
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	return gin.H{"code": 0, "msg": "OK"}, nil
}
