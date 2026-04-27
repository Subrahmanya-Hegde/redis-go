package command

import "time"

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
	return context.Writer.WriteBulkString(data.Data)
}
