package service

import (
	"errors"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"zoe/basic"
	"zoe/dao/db"
	"zoe/model"
	"zoe/utils"
)

func CreateProject(userHash string, req model.CreateProjectRequest) (gin.H, error) {
	user, err := utils.GetUser(userHash)
	if err != nil {
		return nil, err
	}
	org, err := db.QueryOrgById(req.ParentId)
	if err != nil {
		return nil, err
	}
	flag, err := db.ValidateUserForProjectCreation(user.UserHash, req.ParentId)
	if err != nil || !flag {
		return nil, errors.New("用户无权限创建project")
	}
	project, err := db.CreateProject(org.Name+"."+req.Name, req.Private, req.ParentId)
	if err != nil || project == nil {
		return nil, errors.New("用户创建project失败")
	}
	err = db.AddWithCheck(user.UserHash, project.Name, project.Id,
		basic.Resource_Type_PROJECT, user.Id, basic.Privilege_Type_MODIFIER, project.Visibility)
	if err != nil {
		_ = seelog.Critical(err.Error())
		return nil, err
	}
	var visibility string
	if project.Visibility == 0 {
		visibility = "private"
	} else {
		visibility = "public"
	}
	return gin.H{
		"code": 0,
		"msg":  "OK",
		"data": gin.H{
			"id":        project.Id,
			"name":      project.Name,
			"parent_id": project.ParentId,
			"private":   visibility,
		},
	}, nil
}

func UpdateProject(userHash string, projectId int, req model.UpdateProjectRequest) (gin.H, error) {
	user, err := utils.GetUser(userHash)
	if err != nil {
		return nil, err
	}
	project, err := db.GetProjectById(projectId)
	if project == nil || err != nil {
		return nil, errors.New("目标项目不存在")
	}
	flag, err := db.ValidateForUserModifyProject(user.UserHash, projectId, project.ParentId)
	if err != nil || !flag {
		return nil, errors.New("用户无权限修改项目")
	}
	err = db.UpdateProject(projectId, req.Private)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"code": 0,
		"msg":  "OK",
	}, nil
}
