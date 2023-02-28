package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 探测服务器是否存活
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	user := r.Group("/user")
	{
		user.POST("/register", registerHandler)
		user.PATCH("/update", updateHandler)
		user.POST("/delete", deleteHandler)
		user.GET("/search", searchHandler)
	}

	api := r.Group("/api")
	{
		api.GET("")
	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
