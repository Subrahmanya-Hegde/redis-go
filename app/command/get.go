package command

import (
	"time"

	"github.com/subrahmanyahegde/redis-go/app/storage"
)

func handleGet(context *Context) error {
	if len(context.Args) != 1 {
		return context.Writer.WriteError("Argument mismatch. usage: get <key>")
	}
	data, ok := context.Store.Get(context.Args[0].String)
	if !ok {
		return context.Writer.WriteNilString()
	}
	if !data.Expiry.IsZero() && time.Now().After(data.Expiry) {
		return context.Writer.WriteNilString()
	}
	return write(context, data)
}

func write(context *Context, data storage.Data) error {
	switch data.Type {
	case storage.TypeString:
		return context.Writer.WriteBulkString(data.String)
	case storage.TypeList:
		return context.Writer.WriteArray(data.List)
	}
	return context.Writer.WriteError("Unknown type: " + data.Type)
}
