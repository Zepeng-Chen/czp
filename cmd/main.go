package main

import (
	"net/http"

	"github.com/Zepeng-Chen/taurus/handlers/user"
	"github.com/Zepeng-Chen/taurus/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	qps = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_qps",
		Help: "The number of HTTP requests on / served in the last second",
	})
	statusCode = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_status_code",
		Help: "The number of HTTP status code in the last second",
	})
)

func init() {
	prometheus.MustRegister(qps)
	prometheus.MustRegister(statusCode)
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	// 主要路由入口
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	authorized := router.Group("/")
	authorized.Use(middleware.Authentication())

	// 探测服务器是否存活
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// prometheus metrics collections
	router.GET("/metrics", prometheusHandler())

	// frontend
	router.Static("/assets", "../assets")

	// cookie
	router.Use(sessions.Sessions("session", cookie.NewStore()))

	// user router group
	u := router.Group("/user")
	{
		u.POST("/register", user.NewUserRegister)
		u.POST("/login", user.UserLogIn)
		u.GET("/logout", user.UserLogOut)
		u.PATCH("/update", user.UpdateUserInfo)
		u.DELETE("/delete", user.UserLogIn, user.DeleteUser)
		u.GET("/search", user.SearchUser)
		u.GET("/:userid/profile", user.SearchUser)
	}

	// apis
	api := router.Group("/api/v1")
	payment := api.Group("/payment")
	{
		payment.GET("/")
	}

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
