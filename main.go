package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type data struct {
	Title string
}

func DB(action string, name string, addr string, phone string) bool {
	pass := os.Getenv("MYPASS")
	host := os.Getenv("DBHOST")
	if host == "" {
		host = "localhost"
	}
	if pass == "" {
		pass = "password"
	}
	fmt.Println(pass, host)
	db, err := sql.Open("mysql", "root:"+pass+"@tcp("+host+":3306)/")
	if err != nil {
		panic(err.Error())
	}

	switch action {
	case "init":
		if err != nil {
			panic(err.Error())
		}
		_, err = db.Exec("USE Users")
		if err != nil {
			res := strings.Contains(err.Error(), "Unknown database")
			if res {
				db.Exec("CREATE DATABASE Users")
				db.Exec("USE Users")
				defer db.Close()
			}
		} else {
			defer db.Close()
		}
	case "UserAdd":
		db.Exec("USE Users")
		insert, err := db.Query("INSERT INTO Users (name, home_address, phone_number) VALUE(?,?,?)", name, addr, phone)
		if err != nil {
			panic(err.Error())
		} else {
			insert.Close()
			defer db.Close()
		}
	case "UserExists":
		db.Exec("USE Users")
		rows, err := db.Query("SELECT * FROM Users")
		if err != nil {
			res := strings.Contains(err.Error(), "doesn't exist")
			if res {
				db.Exec("CREATE Table Users(id int NOT NULL AUTO_INCREMENT, name varchar(50), home_address varchar(30),phone_number varchar(30), PRIMARY KEY (id));")
				defer db.Close()

			}
		}
		defer rows.Close()
		defer db.Close()
		for rows.Next() {
			var user struct {
				id           int
				name         string
				home_address string
				phone_number string
			}
			err = rows.Scan(&user.id, &user.name, &user.home_address, &user.phone_number)
			if err != nil {
				panic(err.Error())
			}
			if user.name == name && user.home_address == addr && user.phone_number == phone {
				return true
			}
		}
	default:
		return false
	}
	return false
}

func home(w http.ResponseWriter, r *http.Request) {
	d := data{Title: "Please enter yore detail:"}
	if r.Method == "POST" {
		name := r.FormValue("name")
		home_addr := r.FormValue("H_addr")
		phone := r.FormValue("phone")
		fmt.Printf("name is %v\nHome address is %v\nphone number is %v", name, home_addr, phone)
		d2 := data{Title: "Thank you"}
		tmpl, _ := template.ParseFiles("index2.html")
		check := !(DB("UserExists", name, home_addr, phone))
		if check {
			DB("UserAdd", name, home_addr, phone)
		}
		tmpl.Execute(w, d2)
		return
	}
	tmpl, _ := template.ParseFiles("index.html")
	tmpl.Execute(w, d)

}

func main() {
	DB("init", "", "", "")
	DB("UserExists", "", "", "")
	http.HandleFunc("/home", home)
	http.ListenAndServe(":8080", nil)

}
