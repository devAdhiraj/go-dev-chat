package main

import (
	"fmt"
	"os"

	"github.com/devAdhiraj/go-dev-chat/server/api"
	"github.com/devAdhiraj/go-dev-chat/server/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file - %s\n", err)
		panic("error loading .env file")
	}
}

func main() {
	loadEnv()
	models.InitDB()
	r := gin.Default()

	r.POST("/signup", api.SignupHandler)
	r.POST("/login", api.LoginHandler)
	r.GET("/ws", api.AuthMiddleware, api.WsHandler)
	r.GET("/chats", api.AuthMiddleware, api.GetChats)
	r.GET("/current-user", api.AuthMiddleware, api.GetCurrentUser)
	r.GET("/msgs/:convoId", api.AuthMiddleware, api.ConvoMsgsHandler)
	r.GET("/users/:username", api.AuthMiddleware, api.GetUser)
	r.GET("/user/:id", api.AuthMiddleware, api.GetUserById)

	go api.ProducerRoutine()
	go api.ConsumerRoutine()

	r.Run(":8085")
}
