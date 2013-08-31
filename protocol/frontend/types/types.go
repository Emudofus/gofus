package types

import (
	"fmt"
)

// returns 1 if true or 0 if false
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// return true if different of zero, false otherwise
func itob(i int) bool {
	return i != 0
}

type multipleFormatter struct {
	delimiter string
	values    []interface{}
}

func (m multipleFormatter) Format(state fmt.State, verb rune) {
	for i, value := range m.values {
		if i != 0 {
			fmt.Fprint(state, m.delimiter)
		}
		fmt.Fprint(state, value)
	}
}

func NewMultipleFormatter(delimiter string, values ...interface{}) fmt.Formatter {
	return multipleFormatter{delimiter, values}
}
