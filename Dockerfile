FROM golang:wheezy
EXPOSE 3000
RUN go get github.com/justinas/alice
RUN go get github.com/go-sql-driver/mysql