package net

import (
	"fmt"
)

func handle_connect(msg connectMsg) error {
	if msg.new {
		fmt.Fprintf(msg, "HC%s", "abcdefg") // TODO
	}

	return nil
}

func handle_rcv(msg rcvMsg) error {
	return nil // todo
}
