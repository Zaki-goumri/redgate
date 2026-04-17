package resp

import "fmt"

func (v Value) Marshal() []byte {
	switch v.Typ {
	case Array:
		var buf []byte
		buf = append(buf, []byte(fmt.Sprintf("*%d\r\n", len(v.Array)))...)
		for _, val := range v.Array {
			buf = append(buf, val.Marshal()...)
		}
		return buf
	case BulkString:
		if v.Str == nil {
			return []byte("$-1\r\n")
		}
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v.Str), v.Str))
	case SimpleString:
		return []byte(fmt.Sprintf("+%s\r\n", v.Str))
	case Error:
		return []byte(fmt.Sprintf("-%s\r\n", v.Str))
	case Integer:
		return []byte(fmt.Sprintf(":%s\r\n", v.Str))
	case Null:
		return []byte("$-1\r\n")
	default:
		return []byte{}
	}
}
