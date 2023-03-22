package main

import (
	"net/http"

	"github.com/Zepeng-Chen/taurus/api/v1/user"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 主要路由入口
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Taurus, Zepeng!")
	})

	// 探测服务器是否存活
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// user router group
	u := router.Group("/user")
	{
		u.POST("/register", user.NewUserRegister)
		u.POST("/login", user.UserLogIn)
		u.PATCH("/update", user.UpdateUserInfo)
		u.DELETE("/delete", user.UserLogIn, user.DeleteUser)
		u.GET("/search", user.SearchUser)
	}

	// apis
	api := router.Group("/api/v1")
	payment := api.Group("/payment")
	{
		payment.GET("/")
	}

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
