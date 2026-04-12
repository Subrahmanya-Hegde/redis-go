package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Reader struct {
	br *bufio.Reader
}

func NewReader(ioReader io.Reader) *Reader {
	return &Reader{br: bufio.NewReader(ioReader)}
}

func (r *Reader) Read() (Value, error) {
	firstByte, err := r.br.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch firstByte {
	case String:
		return r.readString()
	case Error:
		return r.readError()
	case Integer:
		return r.readInteger()
	case BulkString:
		return r.readBulk()
	case Array:
		return r.readArray()
	}
	return Value{}, fmt.Errorf("unknown RESP type: %c", firstByte)
}

func (r *Reader) readString() (Value, error) {
	value := Value{}
	value.Type = String
	str, err := r.readLine()
	if err != nil {
		return value, err
	}
	value.String = str
	return value, nil
}

func (r *Reader) readError() (Value, error) {
	value := Value{}
	value.Type = Error
	str, err := r.readLine()
	if err != nil {
		return value, err
	}
	value.String = str
	return value, nil
}

func (r *Reader) readInteger() (Value, error) {
	value := Value{}
	value.Type = Integer
	str, err := r.readLine()
	if err != nil {
		return value, err
	}
	number, err := strconv.Atoi(str)
	if err != nil {
		return value, err
	}
	value.Number = number
	return value, nil
}

func (r *Reader) readBulk() (Value, error) {
	lengthStr, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return Value{}, err
	}
	buffer := make([]byte, length)
	_, err = io.ReadFull(r.br, buffer)
	if err != nil {
		return Value{}, err
	}
	_, err = r.readLine()
	if err != nil {
		return Value{}, err
	}
	return Value{Type: BulkString, String: string(buffer)}, nil
}

func (r *Reader) readArray() (Value, error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	arrayLength, err := strconv.Atoi(line)
	if err != nil {
		return Value{}, err
	}
	values := make([]Value, arrayLength)
	for i := range values {
		value, err := r.Read()
		if err != nil {
			return value, err
		}
		values[i] = value
	}
	return Value{Type: Array, Array: values}, nil
}

func (r *Reader) readLine() (string, error) {
	line, err := r.br.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading line: %s", err.Error())
	}
	lineLength := len(line)
	if lineLength < 2 {
		return "", errLineTooShort
	}
	return line[:lineLength-2], nil //Line without \r\n
}
