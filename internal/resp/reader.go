package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Reader struct {
	reader *bufio.Reader
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{reader: bufio.NewReader(rd)}
}

func (r *Reader) readLine() ([]byte, int, error) {
	line, err := r.reader.ReadBytes('\n')
	if err != nil {
		return nil, 0, err
	}
	return line[:len(line)-2], len(line), nil
}

func (r *Reader) Read() (v Value, err error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch _type {
	case Array:
		return r.readArray()
	case BulkString:
		return r.readBulkString()
	default:
		return Value{}, fmt.Errorf("unknown type: %v", string(_type))
	}
}

func (r *Reader) readArray() (Value, error) {
	v := Value{Typ: Array}

	line, _, err := r.readLine()
	if err != nil {
		return v, err
	}

	length, _ := strconv.Atoi(string(line))
	v.Array = make([]Value, length)

	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		v.Array[i] = val
	}

	return v, nil
}
func (r *Reader) readBulkString() (Value, error) {
	v := Value{Typ: BulkString}

	line, _, err := r.readLine()
	if err != nil {
		return v, err
	}

	length, _ := strconv.Atoi(string(line))

	v.Str = make([]byte, length)
	_, err = io.ReadFull(r.reader, v.Str)
	if err != nil {
		return v, err
	}

	r.readLine()

	return v, nil
}
