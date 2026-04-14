package command

func handleEcho(context *Context) error {
	if len(context.Args) != 1 {
		return context.Writer.WriteError("arguments mismatch. Usage ECHO <message>")
	}
	return context.Writer.WriteBulkString(context.Args[0].String)
}
