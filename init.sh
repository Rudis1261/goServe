#!/bin/sh
docker run --rm \
--name goServe \
-p 8080:8080 \
-v `pwd`:/usr/src/myapp \
-w /usr/src/myapp \
drpain:goServe \
go run main.go

