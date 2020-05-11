// reference: https://dinosaurscode.xyz/go/2016/06/19/golang-mysql-authentication/
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
//	"regexp"
)

var db *sql.DB
var err error

func signupPage(res http.ResponseWriter, req *http.Request) {
	fmt.Println("signupPage")
	if req.Method != "POST" {
		http.ServeFile(res, req, "signup.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var user string
	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		res.Write([]byte("User created!"))
		return
	case err != nil:
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	default:
		http.Redirect(res, req, "/", 301)
	}
}

func loginPage(res http.ResponseWriter, req *http.Request) {
	fmt.Println("loginPage")
	if req.Method != "POST" {
		http.ServeFile(res, req, "login.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var databaseUsername string
	var databasePassword string

	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)
	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	res.Write([]byte("Hello " + databaseUsername))
}

func homePage(res http.ResponseWriter, req *http.Request) {
	fmt.Println("homePage")
	http.ServeFile(res, req, "index.html")
}

func main() {
	fmt.Println("start application")
	dbUser := os.Getenv("MYSQL_USER")
//"root"
	dbPass := os.Getenv("MYSQL_PASSWORD")
	dbDBName := os.Getenv("MYSQL_DATABASE")
	dbContainerName := os.Getenv("MYSQL_CONTAINER_NAME")
	dbPort := os.Getenv("MYSQL_PORT")
//	re := regexp.MustCompile("tcp://(.*)")
	//access := fmt.Sprintf("%s:%s@%s/%s", dbUser, dbPass, re.ReplaceAllString(dbURL, "tcp($1)$2"), dbDBName)
	access := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbContainerName, dbPort, dbDBName)
	fmt.Println(access)
	db, err = sql.Open("mysql", access)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":8080", nil)
	fmt.Println("end application")
}
