package router

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"zoe/basic"
	"zoe/controller"
	"zoe/model"
)

type OrgCreateRequest struct {
	Name    string `json:"name" binding:"required"`
	Private bool   `json:"private"`
}

type OrgUpdateRequest struct {
	Private bool `json:"private"`
}

type AuthorizeOrgRequest struct {
	Type   string `json:"type"`
	UserId int    `json:"user_id"`
}

func getUser(c *gin.Context) (*model.User, error) {
	userHash, err := c.Cookie("user_hash")
	if err != nil {
		return nil, err
	}
	user, err := controller.UserController.GetUserByUserHash(userHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func getPrivilegeUserInfo(privileges *[]model.Privilege, users *[]model.User) []gin.H {
	if privileges == nil || len(*privileges) == 0 {
		return nil
	}
	resData := make([]gin.H, len(*privileges))
	for index := range resData {
		var pType string
		typeId := (*privileges)[index].PrivilegeType
		if typeId == 0 {
			pType = "puller"
		} else if typeId == 1 {
			pType = "viewer"
		} else if typeId == 2 {
			pType = "modifier"
		} else {
			pType = "unknown_privilege_type"
		}
		resData[index] = gin.H{
			"id":        (*privileges)[index].Id,
			"type":      pType,
			"user_id":   (*users)[index].Id,
			"user_name": (*users)[index].Name,
		}
	}
	return resData
}

func getProjectInfo(projects *[]model.Project) []gin.H {
	if projects == nil || len(*projects) == 0 {
		return nil
	}
	resData := make([]gin.H, len(*projects))
	for index := range resData {
		var visibility string
		if (*projects)[index].Visibility == 0 {
			visibility = "private"
		} else if (*projects)[index].Visibility == 1 {
			visibility = "public"
		} else {
			visibility = "unknown_visibility_type"
		}
		resData[index] = gin.H{
			"id":         (*projects)[index].Id,
			"name":       (*projects)[index].Name,
			"parent_id":  (*projects)[index].ParentId,
			"visibility": visibility,
		}
	}
	return resData
}

func getOrgInfo(projectInfo []gin.H, privilegeInfo []gin.H, org *model.Org) gin.H {
	var visibility string
	if (*org).Visibility == 0 {
		visibility = "private"
	} else if (*org).Visibility == 1 {
		visibility = "public"
	} else {
		visibility = "unknown_visibility_type"
	}
	orgInfo := gin.H{
		"id":          (*org).Id,
		"name":        (*org).Name,
		"visibility":  visibility,
		"access_mode": "modifier",
	}
	if privilegeInfo != nil {
		orgInfo["privileges"] = privilegeInfo
	}
	if projectInfo != nil {
		orgInfo["projects"] = projectInfo
	}
	return orgInfo
}

func CreateOrgHandler(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "无效的用户"})
		return
	}
	var req OrgCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = seelog.Critical(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "参数错误"})
		return
	}
	if len(req.Name) >= basic.MAX_RESOURCE_NAME_LENGTH {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "组织名长度过长"})
		return
	}
	flag, err := controller.OrgController.IsExistingOrgByName(req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	if flag {
		c.JSON(http.StatusForbidden, gin.H{"code": -1, "msg": "组织已经存在"})
		return
	}
	org, err := controller.OrgController.CreateOrg(req.Name, req.Private)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	err = controller.PrivilegeController.AddWithCheck(user.UserHash, org.Name, org.Id,
		basic.Resource_Type_ORG, user.Id, basic.Privilege_Type_MODIFIER, org.Visibility)
	if err != nil {
		_ = seelog.Critical(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	var visibility string
	if org.Visibility == 0 {
		visibility = "private"
	} else {
		visibility = "public"
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "OK",
		"data": gin.H{
			"id":         org.Id,
			"name":       org.Name,
			"visibility": visibility,
		},
	})
}

func UpdateOrgHandler(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "无效的用户"})
		return
	}
	var req OrgUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = seelog.Critical(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "参数错误"})
		return
	}
	orgId, _ := strconv.Atoi(c.Param("org_id"))
	flag, err := controller.OrgController.IsExistingOrgById(orgId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	if !flag {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "不存在的组织"})
		return
	}

	flag, err = controller.PrivilegeController.ValidateForUserModifyOrg(user.UserHash, orgId)
	if err != nil || !flag {
		c.JSON(http.StatusForbidden, gin.H{"code": -1, "msg": "用户无效权限修改该组织"})
		return
	}
	if err := controller.OrgController.UpdateOrg(orgId, req.Private); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK"})
}

func DeleteOrgHandler(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "无效的用户"})
		return
	}
	orgId, _ := strconv.Atoi(c.Param("org_id"))
	flag, err := controller.PrivilegeController.ValidateForUserModifyOrg(user.UserHash, orgId)
	if err != nil || !flag {
		c.JSON(http.StatusForbidden, gin.H{"code": -1, "msg": "用户无权限修改该组织"})
		return
	}
	if err = controller.OrgController.DeleteOrg(orgId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err})
		return
	}
	if err = controller.PrivilegeController.DeletePrivilege(orgId, basic.Resource_Type_ORG); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "OK",
	})
}

func ListOrgHandler(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "无效的用户"})
		return
	}
	privileges, err := controller.PrivilegeController.ListPrivilege(user.UserHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	orgNameList := make([]string, len(*privileges))
	for index, privilege := range *privileges {
		orgNameList[index] = strings.Split(privilege.ResourceName, ".")[0]
	}
	seelog.Info(orgNameList)
	orgs, err := controller.OrgController.ListOrg(&orgNameList)
	if err != nil {
		_ = seelog.Critical(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	data := make([]gin.H, len(*orgs))
	for index := range data {
		data[index] = gin.H{"id": (*orgs)[index].Id, "name": (*orgs)[index].Name}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "OK",
		"data": data,
	})
}

func SingleOrgHandler(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "无效的用户"})
		return
	}
	orgId, _ := strconv.Atoi(c.Param("org_id"))
	flag, err := controller.OrgController.IsExistingOrgById(orgId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	if !flag {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "不存在的组织"})
		return
	}
	org, err := controller.OrgController.QueryOrgById(orgId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	flag, err = controller.PrivilegeController.ValidateForUserModifyOrg(user.UserHash, orgId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	if flag {
		projects, err := controller.OrgController.ListAllProject(orgId)
		if err != nil {
			_ = seelog.Critical(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
			return
		}
		privileges, err := controller.PrivilegeController.ListPrivilegeByResource(orgId, basic.Resource_Type_ORG)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
			return
		}
		userIds := make([]int, len(*privileges))
		for index, privilege := range *privileges {
			userIds[index] = privilege.UserId
		}
		users, err := controller.UserController.ListUserByIds(userIds)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
			return
		}
		privilegeInfo := getPrivilegeUserInfo(privileges, users)
		projectInfo := getProjectInfo(projects)
		orgInfo := getOrgInfo(projectInfo, privilegeInfo, org)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "OK",
			"data": orgInfo,
		})
	} else {
		flag, err = controller.PrivilegeController.ValidateForUserViewOrg(user.UserHash, orgId)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"code": -1, "msg": err.Error()})
			return
		}
		if flag {
			projects, err := controller.OrgController.ListAllProject(orgId)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
				return
			}
			projectInfo := getProjectInfo(projects)
			orgInfo := getOrgInfo(projectInfo, nil, org)
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "OK",
				"data": orgInfo,
			})
		} else {
			publicProjects, privateProject, err := controller.ProjectController.ListProjectByVisibility(orgId)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
				return
			}
			privileges, err := controller.PrivilegeController.ListPrivilegeByPrefixResourceName(org.Name+".", user.UserHash)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
				return
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
			projectInfo := getProjectInfo(&projects)
			orgInfo := getOrgInfo(projectInfo, nil, org)
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "OK",
				"data": orgInfo,
			})
		}
	}
}

func AuthorizeOrgHandler(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "无效的用户"})
		return
	}
	orgId, _ := strconv.Atoi(c.Param("org_id"))
	flag, err := controller.PrivilegeController.ValidateForUserModifyOrg(user.UserHash, orgId)
	if err != nil || !flag {
		c.JSON(http.StatusForbidden, gin.H{"code": -1, "msg": "用户无权限修改该组织"})
		return
	}
	var req AuthorizeOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = seelog.Critical(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "参数错误"})
		return
	}
	targetUser, err := controller.UserController.GetUserByUserId(req.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "目标用户不存在"})
		return
	}
	if targetUser.Id == user.Id {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "你不能为自己授权"})
		return
	}
	org, err := controller.OrgController.QueryOrgById(orgId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "不存在灯组织"})
		return
	}
	privilege, err := controller.PrivilegeController.QueryPrivilegeByUserHash(targetUser.UserHash, orgId, basic.Resource_Type_ORG)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	var priType int
	if req.Type == "modifier" {
		priType = basic.Privilege_Type_MODIFIER
	} else if req.Type == "viewer" {
		priType = basic.Privilege_Type_VIEWER
	} else if req.Type == "puller" {
		priType = basic.Privilege_Type_PULLER
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "非法的授权类型"})
		return
	}
	if privilege != nil {
		err = controller.PrivilegeController.UpdatePrivilegeByUserHash(targetUser.UserHash, priType, orgId, basic.Resource_Type_ORG)
	} else {
		err = controller.PrivilegeController.CreatePrivilege(targetUser.UserHash, org.Name, orgId, basic.Resource_Type_ORG, targetUser.Id, priType, org.Visibility)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	// todo 删除拉取的缓存

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK"})
}

func DeleteAuthorizeOrgHandler(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "无效的用户"})
		return
	}
	orgId, _ := strconv.Atoi(c.Param("org_id"))
	userId, _ := strconv.Atoi(c.Param("user_id"))
	targetUser, err := controller.UserController.GetUserByUserId(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "目标用户不存在"})
		return
	}
	flag, err := controller.PrivilegeController.ValidateForUserModifyOrg(user.UserHash, orgId)
	if err != nil || !flag {
		c.JSON(http.StatusForbidden, gin.H{"code": -1, "msg": "用户无权限修改该组织"})
		return
	}
	privilege, err := controller.PrivilegeController.QueryPrivilegeByUserHash(targetUser.UserHash, orgId, basic.Resource_Type_ORG)
	if err != nil || privilege == nil {
		c.JSON(http.StatusForbidden, gin.H{"code": -1, "msg": "目标用户无该组织的权限"})
		return
	}
	err = controller.PrivilegeController.DeletePrivilegeByUserHash(targetUser.UserHash, orgId, basic.Resource_Type_ORG)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK"})
}
