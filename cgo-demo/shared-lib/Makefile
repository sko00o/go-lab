
MAKEFLAGS += --silent

.PHONY: all
all: build

.PHONY: run
run:
	# 必须指定共享库的搜索路径
	@export LD_LIBRARY_PATH="$(PWD)" && ./test test.png

.PHONY: build
build:
	go build -o libimgutil.so -buildmode=c-shared imgutil.go
	gcc -o test test.c -L. -limgutil