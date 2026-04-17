package resp

const (
	SimpleString = '+'
	Error        = '-'
	Integer      = ':'
	BulkString   = '$'
	Array        = '*'
	Null         = '0'
)

type Value struct {
	Typ   byte
	Str   []byte
	Num   int
	Array []Value
}
