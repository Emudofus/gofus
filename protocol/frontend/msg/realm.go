package msg

import (
	"fmt"
	"io"
)

type HelloGame struct{}

func (msg *HelloGame) Opcode() string                 { return "HG" }
func (msg *HelloGame) Serialize(out io.Writer) error  { return nil }
func (msg *HelloGame) Deserialize(in io.Reader) error { return nil }

type RealmLoginReq struct {
	Ticket string
}

func (msg *RealmLoginReq) Opcode() string { return "AT" }
func (msg *RealmLoginReq) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.Ticket)
	return nil
}
func (msg *RealmLoginReq) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "AT%s", &msg.Ticket)
	return nil
}
