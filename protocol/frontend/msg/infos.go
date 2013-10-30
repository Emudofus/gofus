package msg

import (
	"fmt"
	"io"
)

type ScreenType int

const (
	InvalidScreenType ScreenType = iota
	NormalScreenType
	FullScreenType
	UnknownScreenType
)

type SetScreenInfos struct {
	Width  int
	Height int
	Type   ScreenType
}

func (msg *SetScreenInfos) Opcode() string { return "Ir" }
func (msg *SetScreenInfos) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "%d;%d;%d", msg.Width, msg.Height, int(msg.Type))
	return nil
}
func (msg *SetScreenInfos) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "Ir%d;%d;%d", &msg.Width, &msg.Height, &msg.Type)
	return nil
}
