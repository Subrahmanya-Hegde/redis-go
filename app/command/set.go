package command

func handleSet(context *Context) error {
	if len(context.Args) != 2 {
		return context.Writer.WriteError("SET command expects 2 arguments")
	}
	context.Store.Set(context.Args[0].String, context.Args[1].String)
	return context.Writer.WriteSimpleString("OK")
}
