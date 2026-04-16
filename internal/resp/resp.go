package resp

const (
	SimpleString = '+'
	Error        = '-'
	Integer      = ':'
	BulkString   = '$'
	Array        = '*'
)

type Value struct {
	Typ   byte
	Str   []byte
	Num   int
	Array []Value
}
