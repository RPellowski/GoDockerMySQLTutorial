FROM golang
RUN go get github.com/go-sql-driver/mysql
RUN go get golang.org/x/crypto/bcrypt
COPY main.go /go/src/myapp/
COPY *.html ./
RUN go install myapp/
EXPOSE 8080
ENTRYPOINT /go/bin/myapp
