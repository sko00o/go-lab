
MAKEFLAGS += --silent

.PHONY: bench1 bench2 test1 test2

bench1:
	# Benchmark with time limit
	go test -run=^$ -bench "^(BenchmarkFib40)$" -v -benchtime=1ms

bench2:
	# Benchmarking various inputs
	go test -run=BenchmarkFib -bench=.

test1:
	# Test one func
	go test -run=^TestSplit$ -v

test2:
	# Individual sub test cases can be executed directly
	go test -run=.*/trailing -v