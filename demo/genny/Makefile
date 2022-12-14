.PHONY: run requirements

all: requirements run

run:
	go generate .
	cat gen_*

requirements:
	test -x genny || go get github.com/cheekybits/genny
