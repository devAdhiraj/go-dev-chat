package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/devAdhiraj/go-dev-chat/server/models"
	"github.com/gin-gonic/gin"
)

const DefaultPaginationLimit = 25

func GetChats(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if _, ok := uid.(uint); !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthenticated",
		})
		return
	}
	userId := uint(uid.(uint))
	msgs, err := models.GetUserChats(userId)
	if err != nil {
		fmt.Println("error fetching msgs - ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "unable to fetch messages"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"chats": msgs})
}

func ConvoMsgsHandler(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if _, ok := uid.(uint); !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthenticated",
		})
		return
	}
	convoId, err := strconv.ParseUint(c.Param("convoId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid convoId",
		})
		return
	}
	limit, err := strconv.ParseUint(c.Query("limit"), 10, 64)
	if err != nil {
		limit = DefaultPaginationLimit
	}
	offset, err := strconv.ParseUint(c.Query("offset"), 10, 64)
	if err != nil {
		offset = 0
	}
	userId := uint(uid.(uint))
	msgs, err := models.GetUserConvo(userId, uint(convoId), int(limit), uint(offset))
	if err != nil {
		fmt.Println("error fetching msgs - ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "unable to fetch messages"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"messages": msgs})
}

func MarkAsSeen(c *gin.Context) {
	var lastSeenMsg models.Msg
	if err := c.ShouldBindJSON(&lastSeenMsg); err != nil ||
		lastSeenMsg.ID == 0 || lastSeenMsg.Timestamp.IsZero() {
		parseAndSetError(c, errors.New("ParamError:invalid msg id or timestamp"))
		return
	}
	if err := models.MarkAsSeen(lastSeenMsg.ID, lastSeenMsg.Timestamp); err != nil {
		fmt.Println("error updating seen timestamp", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error updating seen timestamp",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
