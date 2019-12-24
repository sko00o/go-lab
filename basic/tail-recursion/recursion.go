package tail_recursion

func fact(n int) int {
	if n == 1 {
		return n
	}
	return n * fact(n-1)
}

func tailFact(n int) int {
	return factT(n-1, n)
}

func factT(n, curr int) int {
	if n == 1 {
		return curr
	}
	return factT(n-1, n*curr)
}
