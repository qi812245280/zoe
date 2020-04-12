package utils

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"zoe/dao/db"
	"zoe/model"
)

func GetUser(conn *sql.Tx, userHash string) (*model.User, error) {
	user, err := db.GetUserByUserHash(conn, userHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetPrivilegeUserInfo(privileges *[]model.Privilege, users *[]model.User) []gin.H {
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

func GetProjectInfo(projects *[]model.Project) []gin.H {
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

func GetOrgInfo(projectInfo []gin.H, privilegeInfo []gin.H, org *model.Org) gin.H {
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
