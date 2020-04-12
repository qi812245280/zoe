package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"zoe/basic"
	"zoe/dao/db"
	"zoe/model"
	"zoe/utils"
)

func CreateProject(userHash string, req model.CreateProjectRequest) (gin.H, error) {
	conn, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	user, err := utils.GetUser(conn, userHash)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	org, err := db.QueryOrgById(conn, req.ParentId)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	flag, err := db.ValidateUserForProjectCreation(conn, user.UserHash, req.ParentId)
	if err != nil || !flag {
		_ = conn.Rollback()
		return nil, errors.New("用户无权限创建project")
	}
	var visibility int
	if req.Private == "true" {
		visibility = 0
	} else if req.Private == "false" {
		visibility = 1
	}
	id, err := db.CreateProject(conn, org.Name+"."+req.Name, visibility, req.ParentId)
	if err != nil {
		_ = conn.Rollback()
		return nil, errors.New("用户创建project失败")
	}
	err = db.AddWithCheck(conn, user.UserHash, req.Name, id,
		basic.Resource_Type_PROJECT, user.Id, basic.Privilege_Type_MODIFIER, visibility)
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
			"id":        id,
			"name":      req.Name,
			"parent_id": req.ParentId,
			"private":   visibilityStr,
		},
	}, nil
}

func UpdateProject(userHash string, projectId int, req model.UpdateProjectRequest) (gin.H, error) {
	conn, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	user, err := utils.GetUser(conn, userHash)
	if err != nil {
		_ = conn.Rollback()
		return nil, err
	}
	project, err := db.GetProjectById(conn, projectId)
	if project == nil || err != nil {
		_ = conn.Rollback()
		return nil, errors.New("目标项目不存在")
	}
	flag, err := db.ValidateForUserModifyProject(conn, user.UserHash, projectId, project.ParentId)
	if err != nil || !flag {
		_ = conn.Rollback()
		return nil, errors.New("用户无权限修改项目")
	}
	err = db.UpdateProject(conn, projectId, req.Private)
	if err != nil {
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
