package api

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/devAdhiraj/go-dev-chat/server/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JwtClaim struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	Version  string `json:"version"`
}

var secret_key, _ = hex.DecodeString("D734A99DDB1796FE1F14B68F83A1A") // TODO: switch to env var

func genJwt(username string) (string, error) {
	claims := JwtClaim{
		Username: username,
		Version:  "1.1",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(secret_key)
}

func SignupHandler(c *gin.Context) {
	var usr models.User
	if err := c.ShouldBindJSON(&usr); err != nil {
		fmt.Println("bad req error - ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error - incorrect or missing json body",
		})
		return
	}
	if err := usr.Create(); err != nil {
		parseAndSetError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func LoginHandler(c *gin.Context) {
	var usr models.User
	if err := c.ShouldBindJSON(&usr); err != nil {
		fmt.Println("bad req error - ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error - incorrect or missing json body",
		})
		return
	}
	pwd := usr.Password
	if err := usr.GetByUsername(); err != nil {
		parseAndSetError(c, err)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(pwd)); err != nil {
		fmt.Println("credentials error - ", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid credentials",
		})
		return
	}
	jwtTok, err := genJwt(usr.Username)
	if err != nil {
		fmt.Println("jwt gen error - ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error creating jwt",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"token":   jwtTok,
		"userId":  usr.ID,
	})
}

func AuthMiddleware(c *gin.Context) {
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" || len(bearerToken) <= 7 || bearerToken[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
		c.Abort()
		return
	}
	tokenString := bearerToken[7:]
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		return secret_key, nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid jwt"})
		c.Abort()
		return
	}
	claims, ok := token.Claims.(*JwtClaim)
	if !ok || !token.Valid || claims.Version != "1.1" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	user := models.User{Username: claims.Username}
	if err := user.GetByUsername(); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user's token"})
		c.Abort()
		return
	}

	c.Set("user_id", user.ID)
	c.Set("user", user)
	c.Next()
}
