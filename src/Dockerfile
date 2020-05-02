FROM golang
#RUN go get -d -v github.com/go-sql-driver/mysql
RUN go get github.com/go-sql-driver/mysql
#ADD . /go/src/github.com/go-sql-driver/mysql
COPY main.go /go/src/myapp/
COPY *.html ./
RUN go install myapp/
EXPOSE 8080
ENTRYPOINT /go/bin/myapp
