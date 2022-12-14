#!/usr/bin/env bash

cd "$(dirname "$0")"

go test -fuzz=Fuzz -fuzztime 30s
