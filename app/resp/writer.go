package resp

import (
	"bufio"
	"fmt"
	"io"
)

type Writer struct {
	bw *bufio.Writer
}

func NewWriter(ioWriter io.Writer) *Writer {
	return &Writer{bw: bufio.NewWriter(ioWriter)}
}

func (w *Writer) Write(v Value) error {
	err := w.write(v)
	if err != nil {
		return err
	}
	return w.bw.Flush()
}

func (w *Writer) WriteSimpleString(v string) error {
	return w.Write(Value{Type: String, String: v})
}

func (w *Writer) WriteError(v string) error {
	return w.Write(Value{Type: Error, String: v})
}

func (w *Writer) WriteBulkString(v string) error {
	return w.Write(Value{Type: BulkString, String: v})
}

func (w *Writer) WriteNilString() error {
	_, err := fmt.Fprintf(w.bw, nilString)
	return err
}

func (w *Writer) write(v Value) error {
	switch v.Type {
	case String:
		_, err := fmt.Fprintf(w.bw, stringResponse, v.String)
		return err
	case Integer:
		_, err := fmt.Fprintf(w.bw, integerResponse, v.Number)
		return err
	case Error:
		_, err := fmt.Fprintf(w.bw, errorResponse, v.String)
		return err
	case BulkString:
		_, err := fmt.Fprintf(w.bw, bulkStringResponse, len(v.String), v.String)
		return err
	case Array:
		_, err := fmt.Fprintf(w.bw, "*%d\r\n", len(v.Array))
		if err != nil {
			return err
		}
		for _, val := range v.Array {
			err := w.write(val)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("unknown type %d", v.Type)
}
