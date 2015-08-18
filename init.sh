#!/bin/sh
docker run --rm \
--name goServe \
-p 3000:3000 \
-v `pwd`:/usr/src/myapp \
-w /usr/src/myapp \
drpain:goServe \
go run main.go -port 3000

