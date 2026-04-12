package resp

import "errors"

const (
	String     = '+'
	Error      = '-'
	Integer    = ':'
	BulkString = '$'
	Array      = '*'
)

type Value struct {
	Type   byte
	String string
	Number int
	Array  []Value
}

var stringResponse = "+%s\r\n"
var errorResponse = "-ERR %s\r\n"
var integerResponse = ":%d\r\n"
var bulkStringResponse = "$%d\r\n%s\r\n"

var errLineTooShort = errors.New("line too short")
