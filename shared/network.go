package shared

import (
	"bytes"
	"io"
)

func delimited_buffer(input io.Reader, output chan<- []byte, delimiter []byte, chunkLen int) {
	defer close(output)

	var buffer []byte
	chunk := make([]byte, chunkLen)

	for {
		n, err := input.Read(chunk)

		if n <= 0 || err == io.EOF {
			break
		}
		if err != nil {
			panic(err.Error())
		}

		received := chunk[:n]
		for len(received) > 0 {
			index := bytes.Index(received, delimiter)
			if index < 0 {
				buffer = append(buffer, received...)
				break
			}

			var data []byte
			if len(buffer) > 0 {
				data = make([]byte, index+len(buffer))
				copy(data, buffer)
				copy(data[len(buffer):], received)
			} else {
				data = make([]byte, index)
				copy(data, received)
			}

			output <- data

			received = received[index+len(delimiter):]
		}
	}
}

func Bufferize(input io.Reader, delimiter []byte, chunkLen int) <-chan []byte {
	output := make(chan []byte)
	go delimited_buffer(input, output, delimiter, chunkLen)
	return output
}
