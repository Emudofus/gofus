package shared

import (
	"math/rand"
)

// A Starter implements only one method 'Start'
type Starter interface {
	// Start should start a service without any error, only panic is acceptable
	Start()
}

// A Stopper implements only one method 'Stop'
type Stopper interface {
	// Stop should stop a service without any error, only panic is acceptable
	Stop()
}

// A StartStopper implements the two underlying interfaces Starter and Stopper
type StartStopper interface {
	Starter
	Stopper
}

// Generates a random string given the random source and the fixed length
func NextString(src rand.Source, length int) string {
	const alphanum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, length)
	offset := 0

	for {
		val := src.Int63()
		for i := 0; i < 8; i++ {
			result[offset] = alphanum[val % int64(len(alphanum))]
			length--
			if length == 0 {
				return string(result)
			}
			offset++
			val >>= 8
		}
	}

	panic("unreachable")
}
