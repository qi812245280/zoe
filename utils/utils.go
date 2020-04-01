package utils

import (
	"github.com/gin-gonic/gin"
	"http_guldan_server/controller"
	"http_guldan_server/model"
)

func GetUser(c *gin.Context) (*model.User, error) {
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
