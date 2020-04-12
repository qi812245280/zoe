package controller

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"zoe/model"
	"zoe/service"
)

func CreateProjectHandler(c *gin.Context) {
	userHash, err := c.Cookie("user_hash")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	var req model.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = seelog.Critical(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "参数错误"})
		return
	}
	result, err := service.CreateProject(userHash, req)
	if err != nil {
		_ = seelog.Critical(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func UpdateProjectHandler(c *gin.Context) {
	userHash, err := c.Cookie("user_hash")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	var req model.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = seelog.Critical(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "参数错误"})
		return
	}
	projectId, _ := strconv.Atoi(c.Param("project_id"))
	result, err := service.UpdateProject(userHash, projectId, req)
	if err != nil {
		_ = seelog.Critical(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
