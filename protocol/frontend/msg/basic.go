package msg

import (
	"fmt"
	"io"
	"time"
)

type CurrentDateReq struct{}

func (msg *CurrentDateReq) Opcode() string                 { return "BD" }
func (msg *CurrentDateReq) Serialize(out io.Writer) error  { return nil }
func (msg *CurrentDateReq) Deserialize(in io.Reader) error { return nil }

const current_date_format = "2006|01|02"

type CurrentDateResp struct {
	Current time.Time
}

func (msg *CurrentDateResp) Opcode() string { return "BD" }
func (msg *CurrentDateResp) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.Current.Format(current_date_format))

	return nil
}
func (msg *CurrentDateResp) Deserialize(in io.Reader) error {
	var date string
	fmt.Fscanf(in, "BD%s", &date)

	msg.Current, _ = time.Parse(current_date_format, date)

	return nil
}
