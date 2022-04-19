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

var DBinitCheck bool

func DB(action string, name string, addr string, phone string) bool {
	pass := os.Getenv("MYPASS")
	host := os.Getenv("DBHOST")
	if host == "" {
		host = "localhost"
	}
	if pass == "" {
		pass = "password"
	}
	db, err := sql.Open("mysql", "root:"+pass+"@tcp("+host+":3306)/")
	if err != nil {
		fmt.Println("error with connect to DB")
		panic(err.Error())
	}

	switch action {
	case "init":
		_, err = db.Exec("USE Users")
		for err != nil {
			check := strings.Contains(err.Error(), "Unknown database")
			if check {
				db.Exec("CREATE DATABASE Users")
				_, err = db.Exec("USE Users")
			} else {
				fmt.Print("waiting for database...")
				defer db.Close()
				return false
			}

		}
		_, err = db.Query("SELECT * FROM Users")
		for err != nil {
			check := strings.Contains(err.Error(), "doesn't exist")
			if check {
				db.Exec("CREATE Table Users(id int NOT NULL AUTO_INCREMENT, name varchar(50), home_address varchar(30),phone_number varchar(30), PRIMARY KEY (id));")
				_, err = db.Query("SELECT * FROM Users")
			} else {
				fmt.Print("error with create table")
				panic(err.Error())
				defer db.Close()
				return false
			}
		}
		defer db.Close()
		return true

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
			panic(err.Error())
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
		return false
	default:
		return false
	}
	return false
}

func home(w http.ResponseWriter, r *http.Request) {
	if !DBinitCheck {
		DBinitCheck = DB("init", "", "", "")
		fmt.Println(DBinitCheck)
		wait := data{Title: "loading..."}
		tmpl, _ := template.ParseFiles("index2.html")
		tmpl.Execute(w, wait)
	} else {
		d := data{Title: "Please enter yore detail:"}
		if r.Method == "POST" {
			name := r.FormValue("name")
			home_addr := r.FormValue("H_addr")
			phone := r.FormValue("phone")
			fmt.Printf("name is %v\nHome address is %v\nphone number is %v\n", name, home_addr, phone)
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
}

func main() {
	DBinitCheck = DB("init", "", "", "")
	http.HandleFunc("/home", home)
	http.ListenAndServe(":8080", nil)

}
