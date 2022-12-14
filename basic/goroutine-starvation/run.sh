#!/usr/bin/env bash

cd "$(dirname "$0")"

for image in golang:1.9.7 golang:1.8.7; do
    echo ">> ${image}"
    docker run --rm -it \
        -v "$(pwd)/main.go:/app/main.go" \
        ${image} \
        go run /app/main.go
done
