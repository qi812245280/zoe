package project

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"http_guldan_server/basic"
	"http_guldan_server/controller"
	"http_guldan_server/utils"
	"net/http"
)

type ProjectCreateRequest struct {
	ParentId int    `json:"parent_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Private  string `json:"private"`
}

func CreateProjectHandler(c *gin.Context) {
	user, err := utils.GetUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "无效的用户"})
		return
	}
	var req ProjectCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = seelog.Critical(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "参数错误"})
		return
	}
	org, err := controller.OrgController.QueryOrgById(req.ParentId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	flag, err := controller.PrivilegeController.ValidateUserForProjectCreation(user.UserHash, req.ParentId)
	if err != nil || !flag {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "用户无权限创建project"})
		return
	}
	project, err := controller.ProjectController.CreateProject(org.Name+"."+req.Name, req.Private, req.ParentId)
	if err != nil || project == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "用户创建project失败"})
		return
	}
	err = controller.PrivilegeController.AddWithCheck(user.UserHash, project.Name, project.Id,
		basic.Resource_Type_PROJECT, user.Id, basic.Privilege_Type_MODIFIER, project.Visibility)
	if err != nil {
		_ = seelog.Critical(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	var visibility string
	if project.Visibility == 0 {
		visibility = "private"
	} else {
		visibility = "public"
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "OK",
		"data": gin.H{
			"id":        project.Id,
			"name":      project.Name,
			"parent_id": project.ParentId,
			"private":   visibility,
		},
	})
}
