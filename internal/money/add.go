package money

// AddCents returns a + b. If the sum does not fit in int64, ok is false (sum is the wrapped value and must be ignored).
func AddCents(a, b int64) (sum int64, ok bool) {
	sum = a + b
	if (b > 0 && sum < a) || (b < 0 && sum > a) {
		return sum, false
	}
	return sum, true
}
