FROM golang:wheezy
EXPOSE 8080
RUN go get github.com/justinas/alice
