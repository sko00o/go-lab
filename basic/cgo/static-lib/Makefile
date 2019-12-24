
MAKEFLAGS += --silent

.PHONY: all
all: build

.PHONY: run
run:
	@go run number.go

.PHONY: build
build:
	cd number && gcc -c -o number.o number.c
	cd number && ar rcs libnumber.a number.o