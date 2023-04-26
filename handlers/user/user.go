package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Zepeng-Chen/taurus/internal"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	p "github.com/Zepeng-Chen/taurus/handlers/payment"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Age      *int   `form:"age,omitempty"`
	Phone    *int64 `form:"phone,omitempty"`
	Address  string `form:"address,omitempty"`
	Account  p.Account
}

var userMap = make(map[string]User)

func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "registerForm.html", gin.H{
		"message": "Please fill out the registration form to sign up.",
	})
}

// 注册新用户，如果存在相同用户名就会提示换一个用户名
func NewUserRegister(c *gin.Context) {
	var new_user User
	if err := c.ShouldBind(&new_user); err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("Request not in correct structure"))
		return
	}

	if _, ok := userMap[new_user.Username]; !ok {
		new_user.Password, _ = internal.HashPasswd(new_user.Password)
		userMap[new_user.Username] = new_user
		c.HTML(http.StatusCreated, "loginForm.html", gin.H{
			"message": fmt.Sprintf("User %s has been successfully created. Welcome! Please log in.", new_user.Username),
		})
	} else {
		c.HTML(http.StatusConflict, "registerForm.html", gin.H{
			"message": "Same username has already existed, please change to another one.",
		})
	}
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "loginForm.html", gin.H{
		"message": "Please log in.",
	})
}

// 用户登录，登录后才可以做写操作
func UserLogIn(c *gin.Context) {
	var user User
	if err := c.ShouldBind(&user); err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("Request not in correct structure"))
		return
	}

	if u, ok := userMap[user.Username]; !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "username not found",
		})
		return
	} else if !internal.CheckPasswdHash(user.Password, u.Password) {
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
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "User sign out successfully",
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
