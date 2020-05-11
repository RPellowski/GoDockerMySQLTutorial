# GoDockerMySQLTutorial
A tutorial on using three technologies: Go, Docker and MySQL

---

<div style="color:darkred">

# 🚧 Work In Progress

</div>

Basic instructions to create an application container and a database service container, and deploy them together using Docker.  The application is a Golang web server with a signup page and login page, which read and write persistent data to / from the database.  Docker container deployment is on a local Linux host.  With modification, this application and its service could be enabled to run on a PaaS or IaaS.

Requires general familiarity with Linux but assumes no knowledge of Go, Docker or MySQL.

<!--
Title Simple App
!define MYFONT Dejavu Sans
skinparam titleFontName MYFONT
skinparam state {
  BackgroundColor lightblue
  BorderColor lightblue
  Arrowcolor gray
  FontName MYFONT
  AttributeFontName MYFONT
  ArrowFontName MYFONT
}
d- ->app : build and \ndeploy using \ncreate_app.sh
d- ->db : deploy using \nprovision_db.sh
app-r->db : read and write \npersistent data
state d as "Docker"
state app as "Application" : Source: main.go\nDocker container\n  Image: hobbit\n  Name: frodo
state db as "MySQL DB" : Docker container\n  Image: mysql\n  Name: mysqlshire
-->

![Application](http://www.plantuml.com/plantuml/png/VOt1QkCm48RlUeeXz-G1SbYIR2uBsq99Zu4mbgRs44cZZgG6K_hkbUr0eOMUp7p-_a-xN51B3TuyS_449mwfVVOfcNpbc50nG7CAmRi1EA2zzYerkh_YHS5pFvJELvh-YJhIdtolAZSxurvnD1zcwJd03AkZs2lfwivmPkrrpnOBIrp15avIrT8M0dBSz7AEomQinD8GwJa2_0lODsUGhkCoWKSCxNvHSNAFXpd-C9wU_iFnC9L_2OKnl_glpdpcWPMCMw__O1jtbRq3Z0xEqL7oCaBD7FjsMYKiTC6KDdcO1w4Dlab9vOqpIxouRj9mhOlNnaltfbneCaapvqAnWCE2PaVHrU0_0G00)

## Assumptions
A Linux host is available with the following attributes:
 * Ready to install Docker 
 * Provides a data directory to be used as a persistent backing store for the database
* Root privilege available for installations- root account or one with `sudo` privilege for some operations

 Operations are easier when performed as root; for a non-root user, commands may need additional settings performed with `sudo` and other operations may fail unless privileges are accounted for.

## Tested with versions
 * Ubuntu 20.04 LTS
 * Golang 1.14.2
 * Docker 19.03.8
 * MySQL 8.0.20

## Caveats
This tutorial may become stale as package location, distribution, management and UI commands change for Linux, Docker, Go and other components.

Simple password hashing is used when saving in the database.  Additional security measures used in a web application are not shown.

# Workflow
## Docker
### Install Docker
There are a number of methods that can be used to install Docker.  One way is to follow instructions from Docker.  Another method is to install a `gitlab ci-runner`, which makes use of a single-line script for Docker installation.  Either way is relatively painless.  See the links below for downloads.

**Note:** When using Docker as non-root, add to the group permissions.
```bash
rob@ubuntu:src> groups rob
rob : rob adm cdrom sudo dip plugdev lpadmin lxd sambashare
rob@ubuntu:src> sudo usermod -aG docker rob
rob@ubuntu:src> groups rob
rob : rob adm cdrom sudo dip plugdev lpadmin lxd sambashare docker
```

Once installed, perform simple tests of Docker.

#### Deploy a Hello World container

Since this is a new deployment, Docker will look for the `hello-world` container locally, proceed to download from DockerHub and then run it.

```bash
rob@ubuntu:src> docker run hello-world
Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
0e03bdcc26d7: Pull complete 
Digest: sha256:8e3114318a995a1ee497790535e7b88365222a21771ae7e53687ad76563e8e76
Status: Downloaded newer image for hello-world:latest

Hello from Docker!
This message shows that your installation appears to be working correctly.

To generate this message, Docker took the following steps:
 1. The Docker client contacted the Docker daemon.
 2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
    (amd64)
 3. The Docker daemon created a new container from that image which runs the
    executable that produces the output you are currently reading.
 4. The Docker daemon streamed that output to the Docker client, which sent it
    to your terminal.

To try something more ambitious, you can run an Ubuntu container with:
 $ docker run -it ubuntu bash

Share images, automate workflows, and more with a free Docker ID:
 https://hub.docker.com/

For more examples and ideas, visit:
 https://docs.docker.com/get-started/
```

### Deploy a test Alpine container

Run an interactive version of the Alpine container, a small Linux distro, using the `-i` option and invoking the `/bin/sh` command.  Show the filesystem and running processes.  Exit the container and return to the Linux host.  No containers are running after this operation.  

**Note:** By not specifying a version, `alpine:latest` is pulled and cached into Docker.  Subsequent references to Alpine will be faster unless the latest version has changed in the public repository, in which case the new version will be downloaded before running.

```bash
rob@ubuntu:src> docker run -i -t alpine /bin/sh
Unable to find image 'alpine:latest' locally
latest: Pulling from library/alpine
cbdbe7a5bc2a: Pull complete 
Digest: sha256:9a839e63dad54c3a6d1834e29692c8492d93f90c59c978c1ed79109ea4fb9a54
Status: Downloaded newer image for alpine:latest
/ # ls
bin    dev    etc    home   lib    media  mnt    opt    proc   root   run    sbin   srv    sys    tmp    usr    var
/ # ps
PID   USER     TIME  COMMAND
    1 root      0:00 /bin/sh
    7 root      0:00 ps
/ # exit
rob@ubuntu:src> docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
```
**Note:** The `docker ps` command shows no containers running.
### Deploy a test MySQL container
Run a detached version of the MySQL container (`mysql:latest`). Create the root password (`pass`), set the volume persistence to a local directory `foo`, name it `shire`.

```bash
rob@ubuntu:src> docker run --detach --env MYSQL_ROOT_PASSWORD=pass --name shire mysql
025a6714528268d5bdbc13524d285c21c3f44304aee284a4f4d1bf8ad2a6e2a4
```
This time a container is still running after the command returns.
```bash
rob@ubuntu:src> docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                 NAMES
025a67145282        mysql               "docker-entrypoint.s…"   45 minutes ago      Up 45 minutes       3306/tcp, 33060/tcp   shire
```

### Test MySQL instance
To test, there are two ways to attach to the newly created MySQL instance:
 * Connect directly to the container and run the MySQL CLI inside (the method used here)
 * Use the MySQL CLI instance on localhost (requires installation on Linux)

Use a Docker command to attach to the newly created MySQL instance in the container, run the MySQL CLI and perform some SQL commands.  For this container, the shell is `/bin/bash`.  

Run as the privileged root user.
```
rob@ubuntu:src> docker exec -it shire bash
root@025a67145282:/# mysql -uroot -ppass
mysql: [Warning] Using a password on the command line interface can be insecure.
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 8
Server version: 8.0.20 MySQL Community Server - GPL

Copyright (c) 2000, 2020, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| sys                |
+--------------------+
4 rows in set (0.02 sec)

mysql> show tables from mysql;
+---------------------------+
| Tables_in_mysql           |
+---------------------------+
| columns_priv              |
| component                 |
| db                        |
...
| time_zone_transition_type |
| user                      |
+---------------------------+
33 rows in set (0.00 sec)

mysql> show fields from mysql.user;
+--------------------------+-----------------------------------+------+-----+-----------------------+-------+
| Field                    | Type                              | Null | Key | Default               | Extra |
+--------------------------+-----------------------------------+------+-----+-----------------------+-------+
| Host                     | char(255)                         | NO   | PRI |                       |       |
| User                     | char(32)                          | NO   | PRI |                       |       |
...
+--------------------------+-----------------------------------+------+-----+-----------------------+-------+
51 rows in set (0.03 sec)

mysql> select host, user, account_locked from mysql.user;
+-----------+------------------+----------------+
| host      | user             | account_locked |
+-----------+------------------+----------------+
| %         | root             | N              |
| localhost | mysql.infoschema | Y              |
| localhost | mysql.session    | Y              |
| localhost | mysql.sys        | Y              |
| localhost | root             | N              |
+-----------+------------------+----------------+
5 rows in set (0.01 sec)
```
### Understand how MySQL for App is configured
Create `LOTRdata` database.  Add `users` table to LOTRdata database.  Perform some SQL operations on the newly created table.

```
mysql> CREATE DATABASE IF NOT EXISTS LOTRdata;
Query OK, 1 row affected (0.01 sec)

mysql> CREATE TABLE IF NOT EXISTS LOTRdata.users(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50),
    password VARCHAR(120)
);
Query OK, 0 rows affected (0.03 sec)

mysql> show tables from LOTRdata;
+--------------------+
| Tables_in_LOTRdata |
+--------------------+
| users              |
+--------------------+
1 row in set (0.01 sec)

mysql> show fields from LOTRdata.users;
+----------+--------------+------+-----+---------+----------------+
| Field    | Type         | Null | Key | Default | Extra          |
+----------+------`--------+------+-----+---------+----------------+
| id       | int          | NO   | PRI | NULL    | auto_increment |
| username | varchar(50)  | YES  |     | NULL    |                |
| password | varchar(120) | YES  |     | NULL    |                |
+----------+--------------+------+-----+---------+----------------+
3 rows in set (0.00 sec)

mysql> select * from LOTRdata.users;
Empty set (0.00 sec)

mysql> ^DBye
root@025a67145282:/# exit
rob@ubuntu:src>
```
### Understanding environment variables
To show the effect of the legacy `--link` parameter, two Docker commands are shown here, one without a `--link` parameter and one with the `--link` parameter.  Use of these environment variables for database access can be seen later in the Golang application.  Naming convention by Docker is to prefix the link's exposed communication and environment variables with the uppercase name (in this case, `MYSQL_` and `MYSQL_ENV_`).

Docker Compose is recommended over manually linking containers.  However, use of Compose is not needed for this tutorial.

**Note:** Even though Linux localhost access to the `mysqlshire` instance is through port `13306`, an application will access it directly using the environment (tcp) address and port because Docker manages the mapping.

Without the `--link` parameter, Alpine has a minimal environment.
```
rob@ubuntu:src> docker run --rm --name dummy alpine env
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
HOSTNAME=f382b4a6024b
HOME=/root
```
With the `--link` parameter, Alpine receives more environment context.
```
rob@ubuntu:src> docker run --rm --name dummy --link shire:mysql alpine env
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
HOSTNAME=26c3f9254473
MYSQL_PORT=tcp://172.17.0.2:3306
MYSQL_PORT_3306_TCP=tcp://172.17.0.2:3306
MYSQL_PORT_3306_TCP_ADDR=172.17.0.2
MYSQL_PORT_3306_TCP_PORT=3306
MYSQL_PORT_3306_TCP_PROTO=tcp
MYSQL_PORT_33060_TCP=tcp://172.17.0.2:33060
MYSQL_PORT_33060_TCP_ADDR=172.17.0.2
MYSQL_PORT_33060_TCP_PORT=33060
MYSQL_PORT_33060_TCP_PROTO=tcp
MYSQL_NAME=/dummy/mysql
MYSQL_ENV_MYSQL_ROOT_PASSWORD=LOTRrootpass
MYSQL_ENV_MYSQL_USER=LOTRuser
MYSQL_ENV_MYSQL_PASSWORD=LOTRpass
MYSQL_ENV_MYSQL_DATABASE=LOTRdata
MYSQL_ENV_GOSU_VERSION=1.12
MYSQL_ENV_MYSQL_MAJOR=8.0
MYSQL_ENV_MYSQL_VERSION=8.0.20-1debian10
HOME=/root
```
**Note:** The application will make use of the following environment variables when accessing the MySQL database:
 * `MYSQL_ENV_MYSQL_ROOT_PASSWORD`
 * `MYSQL_PORT`
 * `MYSQL_DATABASE`

### Remove the test MySQL container
```bash
rob@ubuntu:src> docker kill shire
shire
rob@ubuntu:src> docker rm shire
shire
```
### Deploy a test Golang container
Run an interactive version of the Golang container (`golang:latest`).  Compile and run a simple Golang application.  Use the `-v` option to map a Linux directory (`~/hello`) to a container directory (`/go/src`).  For this container, the shell is `/bin/bash`.  

**Note:** Since the container does not have an editor such as `vim` installed, create the file on the Linux host and map the directory to the container.  Editing (using `vim` or another editor) on the Linux host will be reflected inside the container. 

After building in the container, run the executable. 

The binary created in the container is also on the Linux host because of the directory mapping.  Execute it on the Linux host.

```bash
rob@ubuntu:src> mkdir hello
rob@ubuntu:src> cd hello
rob@ubuntu:hello> cat >hello.go <<'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
EOF
rob@ubuntu:hello> docker run -v $PWD:/go/src -i -t golang /bin/bash
Unable to find image 'golang:latest' locally
latest: Pulling from library/golang
90fe46dd8199: Downloading [===========>         ]  25.05MB/50.38MB
35a4f1977689: Download complete 
bbc37f14aded: Download complete 
74e27dc593d4: Downloading [================>    ]  42.08MB/51.83MB
38b1453721cb: Downloading [====                 ]  17.67MB/68.61MB
780391780e20: Waiting 
0f7fd9f8d114: Waiting 
...
Digest: sha256:b451547e2056c6369bbbaf5a306da1327cc12c074f55c311f6afe3bfc1c286b6
Status: Downloaded newer image for golang:latest
root@286532541b79:/go# 
root@721327c4bdda:/go# cd src
root@721327c4bdda:/go/src# go build hello.go
root@721327c4bdda:/go/src# ./hello
Hello, World!
root@721327c4bdda:/go/src# exit
exit
rob@ubuntu:hello> ./hello
Hello, World!
```
<div style="color:darkred">

----

# WIP HERE

----

</div>

## MySQL

### Deploy the real MySQL instance
MySQL container (`mysql:latest`) is deployed with a script (`provision_db.sh`).  The script can be improved by replacing hard-coded items to make use of environment variables and improve security.  

For this instance, port `13306` is selected to avoid collisions in the case that MySQL is already deployed on the Linux host (at default port `3306`).  Docker networking ensures that the container still sees incoming interactions at port `3306.`

```
rob@ubuntu:src> ./provision_db.sh
Starting the MySQL container as 'mysqlshire'
Unable to find image 'mysql:latest' locally
latest: Pulling from library/mysql
54fec2fa59d0: Downloading [================>      ]  24.25MB/27.1MB
bcc6c6145912: Download complete 
951c3d959c9d: Download complete 
d9185034607b: Downloading [===>                   ]  866.7kB/13.44MB
013a9c64dadc: Download complete 
42f3f7d10903: Waiting 
c4a3851d9207: Waiting 
...
bca5ce71f9ea: Pull complete 
Digest: sha256:61a2a33f4b8b4bc93b7b6b9e65e64044aaec594809f818aeffbff69a893d1944
Status: Downloaded newer image for mysql:latest
f49fdd43ddc2692372b01e5bd046561f975b3b07bb7cd69c1bbfd79159cb2316
Database 'LOTRdata' running.
  Username: LOTRuser
  Password: LOTRpass
port 13306
persisting to local directory /home/rob/mydb/mysql-datadir
rob@ubuntu:src> docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                                NAMES
f49fdd43ddc2        mysql               "docker-entrypoint.s…"   2 minutes ago       Up 2 minutes        33060/tcp, 0.0.0.0:13306->3306/tcp   mysqlshire
```

## App
### Build and deploy the app
To build and deploy the app, the following are needed (all source is available):
 * Dockerfile (used by Docker build command)
 * Application source code (main.go)
 * Three files served by the app (index.html, signup.html and login.html)
 * An open port on the Linux host (8082 is selected here)

The build components of a container are specified in the Dockerfile.  The base container is Golang.  A Golang MySQL support library is downloaded from GitHub (which Docker will cache locally).  App source file and .html files are copied into the container at the desired locations.  The executable is built using `go install`.  Since the source code specifies port `8080`, that port is exposed outwards from the container.  Finally, the app binary is specified as the container's entrypoint.

The typical Go workspace is outlined in the following visual: https://talks.golang.org/2014/organizeio.slide#11

**Build** is invoked with a single line.  The container is tagged with `-t hobbit`.  By leaving out a version from the `-t` parameter, the result is the tag  `hobbit:latest`.

**Deploy** is invoked with a single line.  A copy of the container (`hobbit`) is run as `--detached`, meaning not interactively.  The Linux host port that is to be used will be `8082` and mapped to the container's port `8080`. The deployed container is named with `--name frodo`.  Any number of `hobbit` containers could be deployed independently by giving them different names.  In addition, the legacy parameter `--link` is used to make the `frodo` container aware of the `mysql` container. 

```bash
rob@ubuntu:src> docker build -t hobbit .
Sending build context to Docker daemon  11.78kB
Step 1/7 : FROM golang
 ---> 2421885b04da
Step 2/7 : RUN go get github.com/go-sql-driver/mysql
 ---> Running in 5fc38c4a8f59
Removing intermediate container 5fc38c4a8f59
 ---> ceca9a69dd75
Step 3/7 : COPY main.go /go/src/myapp/
 ---> 3890ad46694a
Step 4/7 : COPY *.html ./
 ---> d847f8118d02
Step 5/7 : RUN go install myapp/
 ---> Running in 2505c7b7b1bd
Removing intermediate container 2505c7b7b1bd
 ---> 911af53d9cfb
Step 6/7 : EXPOSE 8080
 ---> Running in bd078437902b
Removing intermediate container bd078437902b
 ---> 429bfea1aded
Step 7/7 : ENTRYPOINT /go/bin/myapp
 ---> Running in 360270a4791e
Removing intermediate container 360270a4791e
 ---> 23cbc59306c1
Successfully built 23cbc59306c1
Successfully tagged hobbit:latest

rob@ubuntu:src> docker run --detach --publish 8082:8080 --name frodo --link mysqlshire:mysql hobbit
7ac8b4a3ee123157cc541a99496d1c3756834c428151cb56ef02b54ac547433f

rob@ubuntu:src> docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                    NAMES
7ac8b4a3ee12        hobbit              "/bin/sh -c /go/bin/m"   41 seconds ago      Up 40 seconds       0.0.0.0:8082->8080/tcp   frodo
c8248de299a6        mysql               "docker-entrypoint.sh"   14 hours ago        Up 14 hours         0.0.0.0:3306->3306/tcp   mysqlshire
```
**Note:** If builds fail with the following
```bash
# cd .; git clone -- https://github.com/go-sql-driver/mysql /go/src/github.com/go-sql-driver/mysql
Cloning into '/go/src/github.com/go-sql-driver/mysql'...
fatal: unable to access 'https://github.com/go-sql-driver/mysql/': Could not resolve host: github.com
package github.com/go-sql-driver/mysql: exit status 128
```
Change the `docker build` command to include the host as network:
```bash
docker build -t hobbit --network=host .
```
### Test the application
The application has been deployed.

Enter the Linux host's address and the exposed application port in a browser to access the homepage, `http://127.0.0.1:8082/`.

Click on the SignUp link to get the the SignUp Page.  Enter a new username and password.  

For demo purposes, the `username/password` was entered as `Bilbo/Baggins`.  Now examine the `mysql` contents.  Again, if the default port 3306 is used, -P can be omitted.

**Note:** Access is via the `LOTRuser` rather than `root`.

```
rob@ubuntu:src> docker exec -it mysqlshire /bin/bash
root@83e3291e1f4d:/# mysql -uLOTRuser -pLOTRpass
...
mysql> select * from LOTRdata.users;
+----+----------+----------+
| id | username | password |
+----+----------+----------+
|  1 | Bilbo    | Baggins  |
+----+----------+----------+
1 row in set (0.00 sec)

mysql>
```
Test the login page using the same `username/password`.

Minimal logging has been enabled in the app.  This can be seen with the `docker logs` command.  The `--follow option` is useful when watching the application in real-time.
```
rob@ubuntu:src> docker logs frodo --follow
start application
homePage
homePage
signupPage
signupPage
loginPage
^C
```
Use `docker inspect` to find out more about the configuration of the container within Docker.
```
rob@ubuntu:src> docker inspect frodo
...
        "Name": "/frodo",
...
            "Image": "hobbit",
...
            "IPAddress": "172.17.0.3",
```
Stop containers
```
rob@ubuntu:src> docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                                NAMES
94f88f802123        hobbit              "/bin/sh -c /go/bin/…"   17 minutes ago      Up 17 minutes       0.0.0.0:8082->8080/tcp               frodo
83e3291e1f4d        mysql               "docker-entrypoint.s…"   19 minutes ago      Up 19 minutes       33060/tcp, 0.0.0.0:13306->3306/tcp   mysqlshire
rob@ubuntu:src> docker stop frodo mysqlshire

frodo
mysqlshire

rob@ubuntu:src> docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
```
Remove containers
```
rob@ubuntu:src> docker rm frodo mysqlshire
frodo
mysqlshire
rob@ubuntu:src> docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
hobbit              latest              477aab4679ff        19 minutes ago      821MB
<none>              <none>              f0227f1557f1        20 minutes ago      821MB
mysql               latest              a7a67c95e831        4 days ago          541MB 
```
Remove images
```
rob@ubuntu:src> docker rmi hobbit mysql
Untagged: hobbit:latest
Deleted: sha256:477aab4679ff902a476f0eeedd3579a13a0101e50a4572cc5b8ee4053a6e42e7
Deleted: sha256:1cad9a2059c23bcf0fb44cf0e4339b9f2903261db917099048931eb148c8a599
...
Untagged: mysql:latest
Untagged: mysql@sha256:61a2a33f4b8b4bc93b7b6b9e65e64044aaec594809f818aeffbff69a893d1944
Deleted: sha256:a7a67c95e83189d60dd24cfeb13d9f235a95a7afd7749a7d09845f303fab239c
Deleted: sha256:7972c7c2b8269f6d954cae13742dea63b6b8b960adacfd2d6c4b3c9dd6f9104b
...
rob@ubuntu:src> docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
```

## Links
### Install GitLab runner
https://docs.gitlab.com/runner/install/linux-repository.html

### Installing Docker
https://docs.docker.com/engine/installation/linux/ubuntu/

### Golang and MySQL Tutorial
https://dinosaurscode.xyz/go/2016/06/19/golang-mysql-authentication/

### Deploying MySQL in a container
https://coreos.com/quay-enterprise/docs/latest/mysql-container.html

## Code Snippets
### Golang
**Hello World** (`hello.go`)
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

**App** (`main.go`)
```go
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
    access := fmt.Sprintf("%s:%s@%s/CNDPdata", dbUser, dbPass, re.ReplaceAllString(dbURL,"tcp($1)$2"))

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
```

### Files used by app (html pages)
`index.html`
```html
<!DOCTYPE html>
<html>
<head>
    <title>Home Page</title>
</head>

<body>
    <h1>Home Page</h1>
    <a href="/login">Login</a>
    <a href="/signup">Sign Up</a>
</body>
</html>
```
`login.html`
```html
<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
</head>

<body>
    <h1>Login Page</h1>
    <form method="POST" action="/login">
        <input type="text" name="username" placeholder="username">
        <input type="password" name="password" placeholder="password">
        <input type="submit" value="Login">
    </form>
</body>
</html>
```
`signup.html`
```html
<!DOCTYPE html>
<html>
<head>
    <title>Sign Up</title>
</head>

<body>
    <h1>SignUp Page</h1>
    <form method="POST" action="/signup">
        <input type="text" name="username" placeholder="username">
        <input type="password" name="password" placeholder="password">
        <input type="submit" value="SignUp">
    </form>
</body>
</html>
```

## Docker-related
### MySQL deployment script
`provision_db.sh`
```bash
#!/bin/bash
# reference: https://coreos.com/quay-enterprise/docs/latest/mysql-container.html
set -e

MYSQL_USER="CNDPuser"
MYSQL_DATABASE="CNDPdata"
MYSQL_CONTAINER_NAME="mysqlshire"
LOCAL_DB_DIR=/root/mydb/mysql-datadir
HOST_PORT=13306

# for better passwords, use
#  $(uuidgen | sed "s/-//g") or $(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | sed 1q)
MYSQL_ROOT_PASSWORD=$(echo CNDProotpass)
MYSQL_PASSWORD=$(echo CNDPpass)

mkdir -p ${LOCAL_DB_DIR}
echo "Starting the MySQL container as '${MYSQL_CONTAINER_NAME}'"

docker \
  run \
  --detach \
  --env MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} \
  --env MYSQL_USER=${MYSQL_USER} \
  --env MYSQL_PASSWORD=${MYSQL_PASSWORD} \
  --env MYSQL_DATABASE=${MYSQL_DATABASE} \
  --name ${MYSQL_CONTAINER_NAME} \
  --volume ${LOCAL_DB_DIR}:/var/lib/mysql \
  --publish ${HOST_PORT}:3306 \
  mysql;

echo "Database '${MYSQL_DATABASE}' running."
echo "  Username: ${MYSQL_USER}"
echo "  Password: ${MYSQL_PASSWORD}"
echo "port ${HOST_PORT}"
echo "persisting to local directory ${LOCAL_DB_DIR}"
```

### App dockerfile
`Dockerfile`
```
FROM golang
RUN go get github.com/go-sql-driver/mysql
COPY main.go /go/src/myapp/
COPY *.html ./
RUN go install myapp/
EXPOSE 8080
ENTRYPOINT /go/bin/myapp
```
### App build and deployment script
`create_app.sh`
```bash
#!/bin/bash
docker build -t hobbit .
docker run --detach --publish 8082:8080 --name frodo --link mysqlshire:mysql hobbit
```
## Commands used in examples
### MySQL
Summary of the mysql commands from the tutorial.

-- remote login --
```
mysql -uroot -pCNDProotpass -h0.0.0.0 -P13306
mysql -uCNDPuser -pCNDPpass -h0.0.0.0 -P13306
```
-- commands --
```
show databases;
show fields from mysql.user;
show fields from CNDPdata.users;
show tables from CNDPdata;
CREATE TABLE CNDPdata.users(....
select host, user, account_locked from mysql.user;
select * from CNDPdata.users;
```
### Docker
Summary of the docker commands from the tutorial.

-- build/run --
```
docker build -t hobbit .
docker run -i -t alpine /bin/sh
docker run -v /root/hello:/go/src -i -t golang /bin/bash
docker run --detach --publish 8082:8080 --name frodo --link mysqlshire:mysql hobbit
docker run --rm --name dummy alpine env
docker run --rm --name dummy --link mysqlshire:mysql alpine env
```
-- container and image management --
```
docker ps
docker stop frodo mysqlshire
docker rm frodo mysqlshire
docker images
docker rmi hobbit mysql
docker logs frodo --follow
docker inspect frodo
```

<!--

------------------------------------
Updates to the tutorial
https://stackoverflow.com/questions/24151129/network-calls-fail-during-image-build-on-corporate-network
edit /etc/default/docker and add the following line:
DOCKER_OPTS="--dns 8.8.8.8 --dns 8.8.4.4"
https://linuxconfig.org/temporary-failure-resolving-error-on-ubuntu-20-04-focal-fossa-linux
https://nickjanetakis.com/blog/setting-up-docker-for-windows-and-wsl-to-work-flawlessly

https://news.ycombinator.com/item?id=11935783
https://web.archive.org/web/20170711063402/http://dinosaurscode.xyz/go/2016/06/19/golang-mysql-authentication/

salting

-----
installing docker
https://docs.docker.com/engine/install/ubuntu/
sudo apt-get remove docker docker-engine docker.io containerd runc
sudo apt-get update
sudo apt-get install \
>     apt-transport-https \
>     ca-certificates \
>     curl \
>     gnupg-agent \
>     software-properties-common

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
sudo apt-get install docker-ce docker-ce-cli containerd.io
-----


## collapsible markdown?

<details><summary>CLICK ME</summary>
<p>

#### yes, even hidden code blocks!

```python
print("hello world!")
```

</p>
</details>

Cleanup
- check all prompts
- change Ubuntu to Linux
- 
-->
