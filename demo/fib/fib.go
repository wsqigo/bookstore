package fib

// Fibonacci returns n-th fibonacci number.
func Fibonacci(n uint) (uint64, error) {
	if n <= 1 {
		return uint64(n), nil
	}

	var n1, n2 uint64 = 1, 0
	for i := uint(2); i < n; i++ {
		n1, n2 = n1+n2, n1
	}
	return n2 + n1, nil
}
