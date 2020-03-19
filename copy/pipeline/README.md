# learning

## fan-int & fan-out

- fan-out: multiple functions can read from the same channel until that channel
    is closed.
- fan-in: a function can read from multiple inputs and proceed until all are
    closed by multiplexing the input channels onto a single channel that's 
    closed when all the inputs are closed.

# ref

- [Go Concurrency Patterns: Pipelines and cancellation](https://blog.golang.org/pipelines)
- [Source Code](https://blog.golang.org/pipelines/)
