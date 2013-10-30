package msg

import (
	"fmt"
	"io"
)

type SetNumberOfFights struct {
	Fights int
}

func (msg *SetNumberOfFights) Opcode() string { return "fC" }
func (msg *SetNumberOfFights) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "%d", msg.Fights)
	return nil
}
func (msg *SetNumberOfFights) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "fC%d", &msg.Fights)
	return nil
}
