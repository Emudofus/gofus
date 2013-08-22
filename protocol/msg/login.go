package msg

import (
	"fmt"
	"io"
)

type HelloConnect struct {
	Ticket string
}

func (msg *HelloConnect) Opcode() string {
	return "HC"
}

func (msg *HelloConnect) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "HC%s", msg.Ticket)
	return nil
}

func (msg *HelloConnect) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "HC%s", &msg.Ticket)
	return nil
}

type BadVersion struct {
	Required string
}

func (msg *BadVersion) Opcode() string {
	return "AlEv"
}

func (msg *BadVersion) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "AlEv%s", msg.Required)
	return nil
}

func (msg *BadVersion) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "AlEv%s", &msg.Required)
	return nil
}

type LoginError struct {}

func (msg *LoginError) Opcode() string {
	return "AlEf"
}

func (msg *LoginError) Serialize(out io.Writer) error {
	fmt.Fprint(out, "AlEf")
	return nil
}

func (msg *LoginError) Deserialize(in io.Reader) error { return nil }

type QueueStatusRequest struct {}

func (msg *QueueStatusRequest) Opcode() string {
	return "Af"
}

func (msg *QueueStatusRequest) Serialize(out io.Writer) error { return nil }
func (msg *QueueStatusRequest) Deserialize(in io.Reader) error { return nil }
