package database

// func init() {
// 	query := `CREATE TABLE IF NOT EXISTS userinfo(userid int primary key auto_increment, username text,
//         password text, created_at datetime default CURRENT_TIMESTAMP, updated_at datetime default CURRENT_TIMESTAMP)`
// 	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancelfunc()
// 	_, err := db.ExecContext(ctx, query)
// 	if err != nil {
// 		log.Printf("Error %s when creating user table", err)
// 		return
// 	}
// 	db, err := sql.Open("mysql", "root:root@tcp(database:3306)/usertable")
// 	defer db.Close()
// 	if err != nil {
// 		panic(err)
// 	}
// 	db.SetConnMaxLifetime(time.Minute * 3)
// 	db.SetMaxOpenConns(10)
// 	db.SetMaxIdleConns(10)
// }

// insert, _ := db.Query("INSERT INTO usertable VALUES (?, ?)", username, passwd)
// defer insert.Close()

// get, _ := db.Query("SELECT password FROM usertable WHERE username = ?", username)
// defer get.Close()
// for get.Next() {
// 	var user User
// 	pwd := get.Scan(&user.Password)
// 	c.JSON(http.StatusOK, gin.H{
// 		"password": pwd,
// 	})
// }
