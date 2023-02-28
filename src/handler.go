package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

var db sql.DB

func init() {
	query := `CREATE TABLE IF NOT EXISTS userinfo(userid int primary key auto_increment, username text,  
        password text, created_at datetime default CURRENT_TIMESTAMP, updated_at datetime default CURRENT_TIMESTAMP)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating user table", err)
		return
	}
	db, err := sql.Open("mysql", "root:root@tcp(database:3306)/usertable")
	defer db.Close()
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func registerHandler(c *gin.Context) {
	// jsonData, err := ioutil.ReadAll(c.Request.Body)
	// if err != nil {
	// 	// Handle error
	// }
	// user := &User{}
	// _ = json.Unmarshal(jsonData, &user)
	// username := user.username
	// passwd := user.password
	username := c.Request.PostForm.Get("username")
	passwd := c.Request.PostForm.Get("password")
	insert, _ := db.Query("INSERT INTO usertable VALUES (?, ?)", username, passwd)
	defer insert.Close()
}

func updateHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func deleteHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func searchHandler(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	get, _ := db.Query("SELECT password FROM usertable WHERE username = ?", username)
	defer get.Close()
	for get.Next() {
		var user User
		pwd := get.Scan(&user.Password)
		c.JSON(http.StatusOK, gin.H{
			"password": pwd,
		})
	}
}
