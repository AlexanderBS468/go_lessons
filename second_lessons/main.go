package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	fmt.Println("Working with MySQL in Go")

	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golangdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("Connected to MySQL database")

	// Insert a new user into the database
	// insert, err := db.Query("INSERT INTO `users` (`name`, `age`) VALUES ('John Doe', 30)")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer insert.Close()
	// fmt.Println("User inserted successfully")

	res, err := db.Query("SELECT `name`, `age` FROM `users`")
	if err != nil {
		panic(err.Error())
	}

	for res.Next() {
		var user User
		err = res.Scan(&user.Name, &user.Age)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("User: %s, Age: %d\n", user.Name, user.Age)
	}

	defer res.Close()
}
