package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Zepeng-Chen/taurus/handlers/user"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

const (
	dbUser   = "root"
	host     = "localhost"
	port     = 3306
	database = "taurus"
)

var password = os.Getenv("DB_PASSWD")

func connectDB() *sql.DB {
	// build the DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, password, host, port, database)
	// Open the connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to create the connection: %s", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}

func createTable() {
	db := connectDB()
	defer db.Close()
	query := `CREATE TABLE IF NOT EXISTS Users(
		id INT primary key auto_increment, 
		userid VARCHAR(100), 
		username VARCHAR(100) NOT NULL,
        password VARCHAR(100) NOT NULL, 
		age INT, 
		phone BIGINT,
		address TEXT,
		created_at datetime default CURRENT_TIMESTAMP, 
		updated_at datetime default CURRENT_TIMESTAMP)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating Users table", err)
		return
	}
}

func insert(db *sql.DB, user user.User) {
	userid := uuid.New()
	insert := "INSERT INTO Users(userid, username, password, age, phone, address) VALUES(?, ?, ?, ?, ?, ?)"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, insert)
	if err != nil {
		log.Fatalf("Error %s when inserting the user's record.", err)
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, userid.String(), user.Username, user.Password, user.Age, user.Phone, user.Address)
}

func query(db *sql.DB, username string) string {
	var user user.User
	if err := db.QueryRow("SELECT password FROM Users WHERE username = ?", username).Scan(&user.Password); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("username %s not found in db", username)
		} else {
			log.Println(err)
		}
	}
	return user.Password
}

func main() {
	db := connectDB()
	createTable()
	user := user.User{Username: "chenzepeng", Password: "123123"}
	fmt.Println(user)
	insert(db, user)
	fmt.Println(query(db, "chenzepeng"))
}
