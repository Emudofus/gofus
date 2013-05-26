package io

import (
	"bytes"
	"io"
)

func Buffer(reader io.Reader, chunkLen int) <-chan []byte {
	out := make(chan []byte)

	go func() {
		defer close(out)

		chunk := make([]byte, chunkLen)
		var buf []byte

		for {
			n, err := reader.Read(chunk)

			if n <= 0 || err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}

			if n >= chunkLen {
				if buf == nil {
					buf = make([]byte, chunkLen)
					copy(buf, chunk)
				} else {
					buf = append(buf, chunk...)
				}
			} else {
				var data []byte

				if buf != nil {
					data = make([]byte, len(buf)+n)
					copy(data, buf)
					copy(data[len(buf):], chunk[:n])

					buf = nil // clean up buffer
				} else {
					data = chunk[:n]
				}

				out <- data
			}
		}
	}()

	return out
}

func BufferLimit(reader io.Reader, chunkLen int, limit []byte) <-chan []byte {
	out := make(chan []byte)

	go func() {
		defer close(out)

		chunk := make([]byte, chunkLen)
		var buf []byte

		for {
			n, err := reader.Read(buf)

			if n <= 0 || err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}

			// add to buffer received chunk
			if buf == nil {
				if n >= chunkLen {
					buf = make([]byte, chunkLen, chunkLen*10)
					copy(buf, chunk)
				} else {
					buf = chunk
				}
			} else {
				buf = append(buf, chunk...)
			}

			// find data
			for {
				index := bytes.Index(buf, limit)

				if index < 0 { // delimiter not found
					break
				}

				data := buf[:index]
				out <- data

				buf = buf[index+len(limit):] // forget previous data

			}

			// clean buffer
			if len(buf) == 0 {
				buf = nil
			}
		}
	}()

	return out
}
