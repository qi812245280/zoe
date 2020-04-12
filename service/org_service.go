package service

import (
	"errors"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"strings"
	"zoe/basic"
	"zoe/dao/db"
	"zoe/model"
	"zoe/utils"
)

func CreateOrg(userHash string, req model.OrgCreateRequest) (gin.H, error) {
	user, err := utils.GetUser(userHash)
	if err != nil {
		return nil, err
	}
	if len(req.Name) >= basic.MAX_RESOURCE_NAME_LENGTH {
		return nil, errors.New("组织名长度过长")
	}
	flag, err := db.IsExistingOrgByName(req.Name)
	if err != nil {
		return nil, err
	}
	if flag {
		return nil, errors.New("组织已经存在")
	}
	org, err := db.CreateOrg(req.Name, req.Private)
	if err != nil {
		return nil, err
	}
	err = db.AddWithCheck(user.UserHash, org.Name, org.Id,
		basic.Resource_Type_ORG, user.Id, basic.Privilege_Type_MODIFIER, org.Visibility)
	if err != nil {
		_ = seelog.Critical(err.Error())
		return nil, err
	}
	var visibility string
	if org.Visibility == 0 {
		visibility = "private"
	} else {
		visibility = "public"
	}
	return gin.H{
		"code": 0,
		"msg":  "OK",
		"data": gin.H{
			"id":         org.Id,
			"name":       org.Name,
			"visibility": visibility,
		},
	}, nil
}

func UpdateOrg(userHash string, orgId int, req model.OrgUpdateRequest) (gin.H, error) {
	user, err := utils.GetUser(userHash)
	if err != nil {
		return nil, err
	}
	flag, err := db.IsExistingOrgById(orgId)
	if err != nil {
		return nil, err
	}
	if !flag {
		return nil, errors.New("不存在的组织")
	}

	flag, err = db.ValidateForUserModifyOrg(user.UserHash, orgId)
	if err != nil || !flag {
		return nil, errors.New("用户无效权限修改该组织")
	}
	if err := db.UpdateOrg(orgId, req.Private); err != nil {
		return nil, err
	}
	return gin.H{"code": 0, "msg": "OK"}, nil
}

func DeleteOrg(userHash string, orgId int) (gin.H, error) {
	user, err := utils.GetUser(userHash)
	if err != nil {
		return nil, err
	}
	flag, err := db.ValidateForUserModifyOrg(user.UserHash, orgId)
	if err != nil || !flag {
		return nil, errors.New("用户无权限修改该组织")
	}
	if err = db.DeleteOrg(orgId); err != nil {
		return nil, err
	}
	if err = db.DeletePrivilege(orgId, basic.Resource_Type_ORG); err != nil {
		return nil, err
	}
	return gin.H{
		"code": 0,
		"msg":  "OK",
	}, nil
}

func ListOrg(userHash string) (gin.H, error) {
	user, err := utils.GetUser(userHash)
	if err != nil {
		return nil, err
	}
	privileges, err := db.ListPrivilege(user.UserHash)
	if err != nil {
		return nil, err
	}
	orgNameList := make([]string, len(*privileges))
	for index, privilege := range *privileges {
		orgNameList[index] = strings.Split(privilege.ResourceName, ".")[0]
	}
	seelog.Info(orgNameList)
	orgs, err := db.ListOrg(&orgNameList)
	if err != nil {
		_ = seelog.Critical(err.Error())
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
	user, err := utils.GetUser(userHash)
	if err != nil {
		return nil, err
	}
	flag, err := db.IsExistingOrgById(orgId)
	if err != nil {
		return nil, err
	}
	if !flag {
		return nil, errors.New("不存在的组织")
	}
	org, err := db.QueryOrgById(orgId)
	if err != nil {
		return nil, err
	}
	flag, err = db.ValidateForUserModifyOrg(user.UserHash, orgId)
	if err != nil {
		return nil, err
	}
	if flag {
		projects, err := db.ListAllProject(orgId)
		if err != nil {
			_ = seelog.Critical(err.Error())
			return nil, err
		}
		privileges, err := db.ListPrivilegeByResource(orgId, basic.Resource_Type_ORG)
		if err != nil {
			return nil, err
		}
		userIds := make([]int, len(*privileges))
		for index, privilege := range *privileges {
			userIds[index] = privilege.UserId
		}
		users, err := db.ListUserByIds(userIds)
		if err != nil {
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
		flag, err = db.ValidateForUserViewOrg(user.UserHash, orgId)
		if err != nil {
			return nil, err
		}
		if flag {
			projects, err := db.ListAllProject(orgId)
			if err != nil {
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
			publicProjects, privateProject, err := db.ListProjectByVisibility(orgId)
			if err != nil {
				return nil, err
			}
			privileges, err := db.ListPrivilegeByPrefixResourceName(org.Name+".", user.UserHash)
			if err != nil {
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
	user, err := utils.GetUser(userHash)
	if err != nil {
		return nil, err
	}
	flag, err := db.ValidateForUserModifyOrg(user.UserHash, orgId)
	if err != nil || !flag {
		return nil, errors.New("用户无权限修改该组织")
	}
	targetUser, err := db.GetUserByUserId(req.UserId)
	if err != nil {
		return nil, errors.New("目标用户不存在")
	}
	if targetUser.Id == user.Id {
		return nil, errors.New("你不能为自己授权")
	}
	org, err := db.QueryOrgById(orgId)
	if err != nil {
		return nil, errors.New("不存在灯组织")
	}
	privilege, err := db.QueryPrivilegeByUserHash(targetUser.UserHash, orgId, basic.Resource_Type_ORG)
	if err != nil {
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
		err = db.UpdatePrivilegeByUserHash(targetUser.UserHash, priType, orgId, basic.Resource_Type_ORG)
	} else {
		err = db.CreatePrivilege(targetUser.UserHash, org.Name, orgId, basic.Resource_Type_ORG, targetUser.Id, priType, org.Visibility)
	}
	if err != nil {
		return nil, err
	}
	// todo 删除拉取的缓存

	return gin.H{"code": 0, "msg": "OK"}, nil
}

func DeleteAuthorizeOrg(userHash string, orgId, userId int) (gin.H, error) {
	user, err := utils.GetUser(userHash)
	if err != nil {
		return nil, err
	}
	targetUser, err := db.GetUserByUserId(userId)
	if err != nil {
		return nil, errors.New("目标用户不存在")
	}
	flag, err := db.ValidateForUserModifyOrg(user.UserHash, orgId)
	if err != nil || !flag {
		return nil, errors.New("用户无权限修改该组织")
	}
	privilege, err := db.QueryPrivilegeByUserHash(targetUser.UserHash, orgId, basic.Resource_Type_ORG)
	if err != nil || privilege == nil {
		return nil, errors.New("目标用户无该组织的权限")
	}
	err = db.DeletePrivilegeByUserHash(targetUser.UserHash, orgId, basic.Resource_Type_ORG)
	if err != nil {
		return nil, err
	}
	return gin.H{"code": 0, "msg": "OK"}, nil
}
