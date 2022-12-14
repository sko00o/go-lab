#!/usr/bin/env bash

cd "$(dirname "$0")"

docker run --rm -v "$(pwd)/.:/app" golang:1.14.15-buster \
    sh -c "cd /app && make"
