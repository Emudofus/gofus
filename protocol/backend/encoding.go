package backend

import (
	"fmt"
	"github.com/Blackrush/gofus/protocol"
	"io"
	"math"
	"time"
)

func Put(out io.Writer, arg interface{}) (n int, err error) {
	switch value := arg.(type) {
	case []byte:
		return out.Write(value)
	case byte:
		return out.Write([]byte{value})
	case int8:
		return Put(out, byte(value))
	case int16:
		return out.Write([]byte{
			byte(value >> 8),
			byte(value & 0xff),
		})
	case uint16:
		return out.Write([]byte{
			byte(value >> 8),
			byte(value & 0xff),
		})
	case int32:
		return out.Write([]byte{
			byte(value >> 24),
			byte(value >> 16 & 0xff),
			byte(value >> 8 & 0xff),
			byte(value & 0xff),
		})
	case uint32:
		return out.Write([]byte{
			byte(value >> 24),
			byte(value >> 16 & 0xff),
			byte(value >> 8 & 0xff),
			byte(value & 0xff),
		})
	case int64:
		return out.Write([]byte{
			byte(value >> 56),
			byte(value >> 48 & 0xff),
			byte(value >> 40 & 0xff),
			byte(value >> 32 & 0xff),
			byte(value >> 24 & 0xff),
			byte(value >> 16 & 0xff),
			byte(value >> 8 & 0xff),
			byte(value & 0xff),
		})
	case uint64:
		return out.Write([]byte{
			byte(value >> 56),
			byte(value >> 48 & 0xff),
			byte(value >> 40 & 0xff),
			byte(value >> 32 & 0xff),
			byte(value >> 24 & 0xff),
			byte(value >> 16 & 0xff),
			byte(value >> 8 & 0xff),
			byte(value & 0xff),
		})
	case float32:
		return Put(out, math.Float32bits(value))
	case float64:
		return Put(out, math.Float64bits(value))
	case string:
		Put(out, uint32(len(value)))
		return Put(out, []byte(value))
	case bool:
		if value {
			return Put(out, byte(1))
		} else {
			return Put(out, byte(0))
		}
	case protocol.Serializer:
		err = value.Serialize(out)
	case time.Time:
		return Put(out, value.UnixNano())

	default:
		panic(fmt.Sprintf("can't convert %T to bytes"))
	}
	return
}

func Read(in io.Reader, arg interface{}) (n int, err error) {
	switch value := arg.(type) {
	case []byte:
		return in.Read(value)
	case *byte:
		tmp := []byte{0}
		n, err = in.Read(tmp)
		*value = tmp[0]
	case *int8:
		tmp := []byte{0}
		n, err = in.Read(tmp)
		*value = int8(tmp[0])
	case *int16:
		tmp := []byte{0, 0}
		n, err = in.Read(tmp)
		*value = int16(tmp[0])<<8 | int16(tmp[1])
	case *uint16:
		tmp := []byte{0, 0}
		n, err = in.Read(tmp)
		*value = uint16(tmp[0])<<8 | uint16(tmp[1])
	case *int32:
		tmp := []byte{0, 0, 0, 0}
		n, err = in.Read(tmp)
		*value = int32(tmp[0])<<24 | int32(tmp[1])<<16 | int32(tmp[2])<<8 | int32(tmp[3])
	case *uint32:
		tmp := []byte{0, 0, 0, 0}
		n, err = in.Read(tmp)
		*value = uint32(tmp[0])<<24 | uint32(tmp[1])<<16 | uint32(tmp[2])<<8 | uint32(tmp[3])
	case *int64:
		tmp := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		n, err = in.Read(tmp)
		*value = int64(tmp[0])<<56 | int64(tmp[1])<<48 | int64(tmp[2])<<40 | int64(tmp[3])<<32 | int64(tmp[4])<<24 | int64(tmp[5])<<16 | int64(tmp[6])<<8 | int64(tmp[7])
	case *uint64:
		tmp := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		n, err = in.Read(tmp)
		*value = uint64(tmp[0])<<56 | uint64(tmp[1])<<48 | uint64(tmp[2])<<40 | uint64(tmp[3])<<32 | uint64(tmp[4])<<24 | uint64(tmp[5])<<16 | uint64(tmp[6])<<8 | uint64(tmp[7])
	case *float32:
		var tmp uint32
		n, err = Read(in, &tmp)
		*value = math.Float32frombits(tmp)
	case *float64:
		var tmp uint64
		n, err = Read(in, &tmp)
		*value = math.Float64frombits(tmp)
	case *string:
		var tmp uint32
		n, err = Read(in, &tmp)
		tmp2 := make([]byte, tmp)
		in.Read(tmp2)
		*value = string(tmp2)
	case *bool:
		tmp := []byte{0}
		n, err = in.Read(tmp)
		*value = tmp[0] == 1
	case protocol.Deserializer:
		err = value.Deserialize(in)
	case *time.Time:
		var tmp int64
		n, err = Read(in, &tmp)
		*value = time.Unix(0, tmp)

	default:
		panic(fmt.Sprintf("can't convert %T to bytes"))
	}
	return
}
