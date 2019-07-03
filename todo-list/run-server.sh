#!/bin/sh

sh gen.sh

go run cmd/server/main.go \
    -db-host=localhost:3308 \
    -db-schema=lab \
    -db-user=root \
    -db-password=toor \
    -grpc-port=9090
