package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (int, error) {
	userIDRaw, exist := c.Get("user_id")
	if !exist {
		return 0, errors.New("user_id not found")
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		return userID, errors.New("user_id not int" + strconv.Itoa(userID))
	}

	return userID, nil
}
