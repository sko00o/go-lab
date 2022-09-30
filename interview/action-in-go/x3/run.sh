#!/bin/bash

# first test
go test ./... -bench=BenchmarkMemoryStack -benchmem -run=^$ -count=10 > stack.txt && tail -f stack.txt

go tool trace *.out