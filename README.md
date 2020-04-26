# GoDockerMySQLTutorial
A tutorial on using three technologies: Go, Docker and MySQL

---

# Work In Progress

---

Basic instructions to create an application container and a database service container, and deploy them together using Docker.  The application is a Golang web server with a signup page and login page, which read and write persistent data to / from the database.  Docker container deployment is on a local Ubuntu host.  With modification, this application and its service could be enabled to run on a PaaS or IaaS.

<!--
![Lack of integration across legacy systems](http://plantuml.com/plantuml/svg/3Sp13S8m3030LM20ndysGEB12gvh4ak97JasgDlJr_tkBaez3qxljnOnrmF0yLUgHCiz5pkP1ciKiW6VR3vGCSo1B7tnXXJobJYtsL6L7GQkk3330rvSaSxd5LJ74DEtszvvb9cZ_m40)

Application
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
![Application]
(http://www.plantuml.com/plantuml/svg/VOt1QkCm48RlUeeXT-S1SWX9DfU5tOKbFGR2M9hQGoIDEf8QJEcxLuv33otqPERpzt_QvO9QQl3cYOidE758xRDFoUGhnuIA0PfJ2DuCm07jTL2fqVqIBmgUXx7qljByJzIHVkTxLRdPEnuK9_Dkbfu3pB0wYhsIhXuCxwozxjbYOOahsC19gbhQG42Ewq7ESTc0bLWQ8Zr7WDy1X-QCqlTSPl0FGxkVLPmyuk4U_pkT_l-1us4k_n0AKtxndvtbp2Ch6TUvRejjtLVs3Z0wE4T7oSi4DNCSsccLiD05KrhdiIswRMY3Br9IUUNC4Y-kdpNiTF7QUEdUb0lD9cdcN2WMS5ZGx2Yw6lm7)

### Assumptions
An Ubuntu host is available with the following attributes:
 * Ready to install Docker (64-bit Ubuntu 18 or 20)
 * Provides a data directory to be used as a persistent backing store for the database
 * Operations are performed as root; if done as non-root user, commands may need `sudo` and other operations may fail unless privileges are accounted for

### Caveats
This tutorial may become stale as package location, distribution, management and UI commands change for Ubuntu, Docker and other components.

## Workflow
### Docker
#### Install Docker
There are a number of methods that can be used to install Docker.  One way is to install a `gitlab ci-runner`, which makes use of a single-line script for Docker installation.  Another method is to follow instructions from Docker.  Either way is relatively painless.  See the links below for downloads.

#### Test the deployment of an Alpine container
Run an interactive version of the Alpine container, a small distro of Linux, using the `-i` option (and invoking the `/bin/sh` command).  Show the filesystem and running processes.  Exit the container and return to Ubuntu.  No containers are running after this operation.  

**Note:** By not specifying a version, `alpine:latest` is pulled and cached into Docker.  Subsequent references to Alpine will be faster unless the latest version has changed in the public repository.

```
root@ubuntu1604:~# docker run -i -t alpine /bin/sh
Unable to find image 'alpine:latest' locally
latest: Pulling from library/alpine
0a8490d0dfd3: Pull complete 
Digest: sha256:dfbd4a3a8ebca874ebd2474f044a0b33600d4523d03b0df76e5c5986cb02d7e8
Status: Downloaded newer image for alpine:latest
/ # ls
bin    dev    etc    home   lib    media  mnt    proc   root   run    sbin   srv    sys    tmp    usr    var
/ # ps
PID   USER     TIME   COMMAND
    1 root       0:00 /bin/sh
    7 root       0:00 ps
/ # exit
root@ubuntu1604:~# docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
<no containers>
```

### Test the deployment of a Golang container
Run an interactive version of the Golang container.  Compile and run a simple Golang application.  Use the `-v` option to map an Ubuntu directory (`~/hello`) to a container directory (`/go/src`).  For this container, the shell is `/bin/bash`.  

**Note:** Since the container does not have an editor such as `vim` installed, create the file on the Ubuntu host and map the directory to the container.  Editing (using `vim` or another editor) on the Ubuntu host will be reflected inside the container.  The binary created during the container build is also on the Ubuntu host.

```
root@ubuntu1604:~# mkdir hello
root@ubuntu1604:~# cd hello
root@ubuntu1604:~/hello# 
root@ubuntu1604:~/hello# cat >hello.go <<'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
EOF
root@ubuntu1604:~/hello# docker run -v /root/hello:/go/src -i -t golang /bin/bash
root@a12a96685191:/go# cd src
root@a12a96685191:/go/src# go build hello.go
root@a12a96685191:/go/src# ./hello
Hello, World!
root@a12a96685191:/go/src# exit
root@ubuntu1604:~/hello# ./hello
Hello, World!
```

## MySQL
### Deploy MySQL instance
MySQL is deployed with a script (`provision_db.sh`, below).  The script can be improved by replacing hard-coded items to make use of environment variables and improve security.  

For this instance, port 13306 is selected to avoid collisions in the case that MySQL is already deployed on the Ubuntu host (at default port 3306).  The container still sees incoming interactions at port 3306.

```
root@ubuntu1604:~# ./provision_db.sh 
Starting the MySQL container as 'mysqlshire'
Unable to find image 'mysql:latest' locally
latest: Pulling from library/mysql
5040bd298390: Already exists 
55370df68315: Pull complete 
...
Digest: sha256:5e2ec5964847dd78c83410f228325a462a3bfd796b6133b2bdd590b71721fea6
Status: Downloaded newer image for mysql:latest
76cac61dc5d09a3aeb4d1849ad070a6c8acd4c035c76456fe16c6922a0bc4227
Database 'CNDPdata' running.
  Username: CNDPuser
  Password: CNDPpass
port 13306
persisting to local directory /root/mydb/mysql-datadir

root@ubuntu1604:~# docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                     NAMES
76cac61dc5d0        mysql               "docker-entrypoint..."   2 minutes ago       Up 2 minutes        0.0.0.0:13306->3306/tcp   mysqlshire
```

### Test MySQL instance
Use the MySQL CLI to attach (as root) to the newly created MySQL instance on localhost and inspect some of the MySQL content.  If the default port 3306 is used, `-P` can be omitted.

```
root@ubuntu1604:~# mysql -uroot -pCNDProotpass -h0.0.0.0 -P13306
mysql: [Warning] Using a password on the command line interface can be insecure.
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 4
Server version: 5.7.17 MySQL Community Server (GPL)

Copyright (c) 2000, 2016, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| CNDPdata           |
| mysql              |
| performance_schema |
| sys                |
+--------------------+
5 rows in set (0.00 sec)

mysql> show fields from mysql.user;
+------------------------+-----------------------------------+------+-----+-----------------------+-------+
| Field                  | Type                              | Null | Key | Default               | Extra |
+------------------------+-----------------------------------+------+-----+-----------------------+-------+
| Host                   | char(60)                          | NO   | PRI |                       |       |
| User                   | char(32)                          | NO   | PRI |                       |       |
...

mysql> select host, user, account_locked from mysql.user;
+-----------+-----------+----------------+
| host      | user      | account_locked |
+-----------+-----------+----------------+
| localhost | root      | N              |
| localhost | mysql.sys | Y              |
| %         | root      | N              |
| %         | CNDPuser  | N              |
+-----------+-----------+----------------+
4 rows in set (0.00 sec)

mysql> show tables from CNDPdata;
Empty set (0.00 sec)
Configure MySQL for app
Add table to CNDPdata database.  Perform some SQL operations on the database and the newly created table.

mysql> CREATE TABLE CNDPdata.users(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50),
    password VARCHAR(120)
);
 
mysql> show tables from CNDPdata;
+--------------------+
| Tables_in_CNDPdata |
+--------------------+
| users              |
+--------------------+
1 row in set (0.00 sec)

mysql> show fields from CNDPdata.users;
+----------+--------------+------+-----+---------+----------------+
| Field    | Type         | Null | Key | Default | Extra          |
+----------+--------------+------+-----+---------+----------------+
| id       | int(11)      | NO   | PRI | NULL    | auto_increment |
| username | varchar(50)  | YES  |     | NULL    |                |
| password | varchar(120) | YES  |     | NULL    |                |
+----------+--------------+------+-----+---------+----------------+
3 rows in set (0.00 sec)

mysql> select * from CNDPdata.users;
Empty set (0.01 sec)
```

## App
### Build and deploy app
To build and deploy the app, the following are needed (all source is contained in the code section below):
 * Dockerfile (used by Docker build command)
 * main.go (application source code)
 * three .html files served by the app (index, signup and login)
 * an open port on the Ubuntu host (8082 is selected here)

The build components of a container are specified in the Dockerfile.  The base container is Golang.  A Golang MySQL support library is downloaded from GitHub (which Docker will cache locally).  App source file and .html files are copied into the container at the desired locations.  The executable is built using `go install`.  Since the source code specifies port 8080, that port is exposed outwards from the container.  Finally, the app binary is specified as the container's entrypoint.

Build is invoked with a single line.  The container is tagged with `-t hobbit`.  By leaving out a version from the `-t` parameter, it is actually tagged with `hobbit:latest`.

Deploy (run) is invoked with a single line.  A copy of the container (`hobbit`) is run as `--detached`, meaning not interactively.  The Ubuntu port that is to be used will be 8082 and mapped to the container's port 8080. The deployed container is named with `--name frodo`.  Any number of `hobbit` containers could be deployed independently by giving them different names.  In addition, the legacy parameter `--link` is used to make the `frodo` container aware of the `mysql` container. 

The typical Go workspace is outlined in the following visual: https://talks.golang.org/2014/organizeio.slide#11

**Note:** The following `docker build` command may send a large amount build context to Docker daemon, which might be dependent on the docker images you have on your system.  For example, on one test instance, 61.43 GB of build context was sent to the Docker daemon before the `hobbit` build was complete.

<!--

        root@ubuntu1604:~# docker build -t hobbit .

    Sending build context to Docker daemon 61.43 GB

         ...

 

root@ubuntu1604:~# docker build -t hobbit .
Sending build context to Docker daemon 405.4 MB
Step 1 : FROM golang
 - - -> 9752d71739d2
Step 2 : RUN go get github.com/go-sql-driver/mysql
 - - -> Running in c196b5f131fb
 - - -> 94273fecf75a
Removing intermediate container c196b5f131fb
Step 3 : COPY main.go /go/src/myapp/
 - - -> e24a61023112
Removing intermediate container aabadb00ccc0
Step 4 : COPY *.html ./
 - - -> e66a3c648b29
Removing intermediate container 13f480b1fd93
Step 5 : RUN go install myapp/
 - - -> Running in b6544cc6015b
 - - -> 7029199037cd
Removing intermediate container b6544cc6015b
Step 6 : EXPOSE 8080
 -- -> Running in 693339d2302f
 -- -> 80dc41ad479d
Removing intermediate container 693339d2302f
Step 7 : ENTRYPOINT /go/bin/myapp
 -- -> Running in 9b8b0bf9184d
 -- -> 41c564d723d4
Removing intermediate container 9b8b0bf9184d
Successfully built 41c564d723d4

root@ubuntu1604:~# docker images
REPOSITORY                                     TAG                 IMAGE ID            CREATED              SIZE
hobbit                                         latest              41c564d723d4        About a minute ago   681.9 MB
mysql                                          latest              7666f75adb6b        3 weeks ago          405.6 MB
 
root@ubuntu1604:~# docker run --detach --publish 8082:8080 --name frodo --link mysqlshire:mysql hobbit
7ac8b4a3ee123157cc541a99496d1c3756834c428151cb56ef02b54ac547433f

root@ubuntu1604:~# docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                    NAMES
7ac8b4a3ee12        hobbit              "/bin/sh -c /go/bin/m"   41 seconds ago      Up 40 seconds       0.0.0.0:8082->8080/tcp   frodo
c8248de299a6        mysql               "docker-entrypoint.sh"   14 hours ago        Up 14 hours         0.0.0.0:3306->3306/tcp   mysqlshire
-->

### Examining environment variables
To show the effect of the legacy `--link` parameter, two Docker commands are shown here, one without a `--link` parameter and one with the `--link` parameter.  Use of these environment variables for database access can be seen in the Golang application.  Naming convention by Docker is to prefix the link's exposed communication and environment variables with the uppercase name (in this case, `MYSQL_` and `MYSQL_ENV_`). 

The use of Docker Compose is now recommended over manually linking containers.  However, use of Compose is not needed for this tutorial.

**Note:** Even though Ubuntu localhost access to the `mysqlshire` instance is through port 13306, an application will access it directly using the environment (tcp) address and port.

```
root@ubuntu1604:~# docker run --rm --name dummy alpine env
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
HOSTNAME=7e32ca8ba459
HOME=/root

root@ubuntu1604:~# docker run --rm --name dummy --link mysqlshire:mysql alpine env
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
HOSTNAME=ca5ff99f8fce
MYSQL_PORT=tcp://172.17.0.2:3306
MYSQL_PORT_3306_TCP=tcp://172.17.0.2:3306
MYSQL_PORT_3306_TCP_ADDR=172.17.0.2
MYSQL_PORT_3306_TCP_PORT=3306
MYSQL_PORT_3306_TCP_PROTO=tcp
MYSQL_NAME=/dummy/mysql
MYSQL_ENV_MYSQL_ROOT_PASSWORD=CNDProotpass
MYSQL_ENV_MYSQL_USER=CNDPuser
MYSQL_ENV_MYSQL_PASSWORD=CNDPpass
MYSQL_ENV_MYSQL_DATABASE=CNDPdata
MYSQL_ENV_GOSU_VERSION=1.7
MYSQL_ENV_MYSQL_MAJOR=5.7
MYSQL_ENV_MYSQL_VERSION=5.7.17-1debian8
HOME=/root
```

### Test app
The application has been deployed.

Enter the Ubuntu host's address and the exposed application port in a browser to access the homepage.

Click on the SignUp link to get the the SignUp Page.  Enter a new username and password.  

For demo purposes, the `username/password` was entered as `Bilbo/Baggins`.  Now examine the `mysql` contents.  Again, if the default port 3306 is used, -P can be omitted.

**Note:** This time, access is via the `CNDPuser` rather than `root`.

```
root@ubuntu1604:~# mysql -uCNDPuser -pCNDPpass -h0.0.0.0 -P13306
mysql: [Warning] Using a password on the command line interface can be insecure.
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 11
Server version: 5.7.17 MySQL Community Server (GPL)

Copyright (c) 2000, 2016, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> select * from CNDPdata.users;
+----+----------+----------+
| id | username | password |
+----+----------+----------+
|  1 | Bilbo    | Baggins  |
+----+----------+----------+
1 row in set (0.00 sec)
```

Test the login page using the same `username/password`.

Minimal logging has been enabled in the app.  This can be seen with the `docker logs` command.  The `--follow option` is useful when watching the application in real-time.

Use `docker inspect` to find out more about the configuration of the container within Docker.

```
root@ubuntu1604:~# docker logs frodo --follow
start application
signupPage
signupPage
homePage
loginPage
^C
 
root@ubuntu1604:~# docker inspect frodo
...
        "Name": "/frodo",
...
            "Image": "hobbit",
...
            "IPAddress": "172.17.0.3",
```
Stop containers, remove containers, remove images
```
root@ubuntu1604:~# docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                    NAMES
7ac8b4a3ee12        hobbit              "/bin/sh -c /go/bin/m"   About an hour ago   Up About an hour    0.0.0.0:8082->8080/tcp   frodo
c8248de299a6        mysql               "docker-entrypoint.sh"   16 hours ago        Up 16 hours         0.0.0.0:3306->3306/tcp   mysqlshire

root@ubuntu1604:~# docker stop frodo mysqlshire
frodo
mysqlshire

root@ubuntu1604:~# docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES

root@ubuntu1604:~# docker rm frodo mysqlshire
frodo
mysqlshire

root@ubuntu1604:~# docker images
REPOSITORY                                     TAG                 IMAGE ID            CREATED             SIZE
hobbit                                         latest              cca46f65dec4        About an hour ago   681.9 MB
<none>                                         <none>              41c564d723d4        About an hour ago   681.9 MB
mysql                                          latest              7666f75adb6b        3 weeks ago         405.6 MB
...
 
root@ubuntu1604:~# docker rmi hobbit mysql
Untagged: hobbit:latest
Deleted: sha256:cca46f65dec43fcf1a74a89a6b1a1f571f081f53c6c5b83493eafabcb4d5ad35
Deleted: sha256:7cf253249c61fc2f1987eb13a1d1ed58544e248dc27ff3c49175a54c7290d1c6
Deleted: sha256:3d9626d3cecdd3b6ec015e080c3715605e0359e15a69953a620e0068f73265a1
Deleted: sha256:f766e9b2106de36bbdcaf81c0668c0203fceed0e44db97fe1f91f601f9368806
Deleted: sha256:398d43513416532543a628b2c420a2c1fdd5160e3898e2a94f5a6c8433d7c380
Deleted: sha256:94273fecf75afbc4fe121794f8e19894b3a21e57c159d377a76330c8b577ed10
Untagged: mysql:latest
Untagged: mysql@sha256:5e2ec5964847dd78c83410f228325a462a3bfd796b6133b2bdd590b71721fea6
Deleted: sha256:7666f75adb6b50676a366c6fd7a3916cb41f6e8eaf336c3d3ab7d35317fed0b9

root@ubuntu1604:~# docker images
REPOSITORY                                     TAG                 IMAGE ID            CREATED             SIZE
...
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
**Hello World**
`hello.go`
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}```

### App
`main.go`
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
login.html

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
`dockerfile` or `Dockerfile`
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
-->
