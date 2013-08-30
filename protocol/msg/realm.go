package msg

import (
	"fmt"
	"io"
)

type HelloGame struct{}

func (msg *HelloGame) Opcode() string                 { return "HG" }
func (msg *HelloGame) Serialize(out io.Writer) error  { return nil }
func (msg *HelloGame) Deserialize(in io.Reader) error { return nil }

type RealmLogin struct {
	Ticket string
}

func (msg *RealmLogin) Opcode() string { return "AT" }
func (msg *RealmLogin) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.Ticket)
	return nil
}
func (msg *RealmLogin) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "AT%s", &msg.Ticket)
	return nil
}
