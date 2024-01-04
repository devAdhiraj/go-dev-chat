package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/devAdhiraj/go-dev-chat/server/models"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	var usr models.User
	usr.Username = c.Param("username")
	if err := usr.GetByUsername(); err != nil {
		parseAndSetError(c, err)
		return
	}
	c.JSON(http.StatusOK, usr)
}

func GetUserById(c *gin.Context) {
	var usr models.User

	userId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid userId",
		})
		return
	}
	usr.ID = uint(userId)

	if err := usr.GetById(); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}
	usr.Password = ""
	c.JSON(http.StatusOK, usr)
}

func GetCurrentUser(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if _, ok := uid.(uint); !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthenticated",
		})
		return
	}
	u, _ := c.Get("user")
	usr, ok := u.(models.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid user",
		})
		return
	}
	fmt.Println("user = ", usr)

	c.JSON(http.StatusOK, usr)
}

func parseAndSetError(c *gin.Context, err error) {
	fmt.Println("error - ", err)
	errSlice := strings.Split(err.Error(), ":")
	errType := errSlice[0]
	errMsg := "error creating user"
	switch errType {
	case "ParamError":
		errMsg = errSlice[1]
	case "DBError":
		errMsg = "username not available"
	case "PasswordError":
		errMsg = "invalid password"
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": errMsg,
	})
}
