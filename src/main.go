// reference: https://dinosaurscode.xyz/go/2016/06/19/golang-mysql-authentication/
package main

import (
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "os"
    "regexp"
    "net/http"
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
        _, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, password)
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

    if databasePassword != password {
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
    dbUser := "root"
    dbPass := os.Getenv("MYSQL_ENV_MYSQL_ROOT_PASSWORD")
    dbURL := os.Getenv("MYSQL_PORT")
    re := regexp.MustCompile("tcp://(.*)")
    access := fmt.Sprintf("%s:%s@%s/LOTRdata", dbUser, dbPass, re.ReplaceAllString(dbURL,"tcp($1)$2"))

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
