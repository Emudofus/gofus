package msg

// Converts a boolean value to an integer (1 if true or 0 if false)
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
