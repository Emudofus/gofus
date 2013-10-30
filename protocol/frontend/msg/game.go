package msg

import (
	"fmt"
	"io"
)

type SetCurrentMapData struct {
	Id   int
	Date string
	Key  string
}

func (msg *SetCurrentMapData) Opcode() string { return "GDM" }

func (msg *SetCurrentMapData) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "|%d|%s|%s", msg.Id, msg.Date, msg.Key)
	return nil
}

func (msg *SetCurrentMapData) Deserialize(in io.Reader) error { return nil }

type ContextInfosReq struct{}

func (msg *ContextInfosReq) Opcode() string                 { return "GI" }
func (msg *ContextInfosReq) Serialize(out io.Writer) error  { return nil }
func (msg *ContextInfosReq) Deserialize(in io.Reader) error { return nil }

type SetMapLoaded struct{}

func (msg *SetMapLoaded) Opcode() string                 { return "GDK" }
func (msg *SetMapLoaded) Serialize(out io.Writer) error  { return nil }
func (msg *SetMapLoaded) Deserialize(in io.Reader) error { return nil }
