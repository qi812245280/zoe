package controller

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"zoe/model"
	"zoe/service"
)

func CreateOrgHandler(c *gin.Context) {
	userHash, err := c.Cookie("user_hash")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	var req model.OrgCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = seelog.Critical(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "参数错误"})
		return
	}
	result, err := service.CreateOrg(userHash, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func UpdateOrgHandler(c *gin.Context) {
	userHash, err := c.Cookie("user_hash")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	var req model.OrgUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = seelog.Critical(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "参数错误"})
		return
	}
	orgId, _ := strconv.Atoi(c.Param("org_id"))
	result, err := service.UpdateOrg(userHash, orgId, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func DeleteOrgHandler(c *gin.Context) {
	userHash, err := c.Cookie("user_hash")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	orgId, _ := strconv.Atoi(c.Param("org_id"))
	result, err := service.DeleteOrg(userHash, orgId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err})
		return
	}
	c.JSON(http.StatusOK, result)
}

func ListOrgHandler(c *gin.Context) {
	userHash, err := c.Cookie("user_hash")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	result, err := service.ListOrg(userHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err})
		return
	}
	c.JSON(http.StatusOK, result)
}

func SingleOrgHandler(c *gin.Context) {
	userHash, err := c.Cookie("user_hash")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	orgId, _ := strconv.Atoi(c.Param("org_id"))
	result, err := service.SingleOrg(userHash, orgId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err})
		return
	}
	c.JSON(http.StatusOK, result)
}

func AuthorizeOrgHandler(c *gin.Context) {
	userHash, err := c.Cookie("user_hash")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	orgId, _ := strconv.Atoi(c.Param("org_id"))
	var req model.AuthorizeOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = seelog.Critical(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "参数错误"})
		return
	}
	result, err := service.AuthorizeOrg(userHash, orgId, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func DeleteAuthorizeOrgHandler(c *gin.Context) {
	userHash, err := c.Cookie("user_hash")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	orgId, _ := strconv.Atoi(c.Param("org_id"))
	userId, _ := strconv.Atoi(c.Param("user_id"))
	result, err := service.DeleteAuthorizeOrg(userHash, orgId, userId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
