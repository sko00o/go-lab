
MAKEFLAGS += --silent

.PHONY: run cpuprofile memprofile requirements

cpuprofile: run
	go tool pprof -pdf cpuprofile > cpuprofile.pdf

memprofile: run
	go tool pprof -pdf memprofile > memprofile.pdf

run:
	go run main.go

requirements:
	sudo pacman -Syy
	sudo pacman -S graphviz
