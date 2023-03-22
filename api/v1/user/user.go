package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Age      *int   `json:"age,omitempty"`
	Phone    *int64 `json:"phone,omitempty"`
	Address  string `json:"address,omitempty"`
	Account  Account
}

type Account struct {
	AccountID    string
	Balance      credit
	LastModified time.Time
}

type credit float64

var userMap = make(map[string]User)

// 注册新用户，如果存在相同用户名就会提示换一个用户名
func NewUserRegister(c *gin.Context) {
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	new_user := User{}
	json.Unmarshal(jsonData, &new_user)
	if _, ok := userMap[new_user.Username]; !ok {
		userMap[new_user.Username] = new_user
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": fmt.Sprintf("User %s has just been created. Welcome!", new_user.Username),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "Same username has already existed, please change to another one.",
		})
	}
}

// 用户登录，登录后才可以做写操作
func UserLogIn(c *gin.Context) {
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	user := User{}
	json.Unmarshal(jsonData, &user)

	if u, ok := userMap[user.Username]; !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "username not found",
		})
		return
	} else if u.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "Authentication failed, password mismatch. Please try again.",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "log in succeeded",
		})
		c.Next()
	}
}

// 用户退出登录，退出后将不可以做写操作
func UserLogOut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "log out succeeded",
	})
}

// 修改用户基本信息
func UpdateUserInfo(c *gin.Context) {
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	user := User{}
	json.Unmarshal(jsonData, &user)

	if u, ok := userMap[user.Username]; !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "username not found",
		})
	} else {
		*u.Age = *user.Age
		*u.Phone = *user.Phone
		fmt.Println("Age is", *u.Age)
		fmt.Println("Phone is", *u.Phone)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "update succeeded",
		})
	}
}

// 删除用户
func DeleteUser(c *gin.Context) {
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	user := User{}
	json.Unmarshal(jsonData, &user)
	delete(userMap, user.Username)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "user deleted",
	})
}

// 查询用户信息
func SearchUser(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	item := c.Query("item")

	if user, ok := userMap[username]; !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "user not found",
		})
	} else {
		switch item {
		case "phone":
			c.JSON(http.StatusOK, gin.H{
				"code":     0,
				"username": user.Username,
				"phone":    user.Phone,
			})
		case "age":
			c.JSON(http.StatusOK, gin.H{
				"code":     0,
				"username": user.Username,
				"age":      user.Age,
			})
		default:
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "invalid item",
			})
		}
	}
}
