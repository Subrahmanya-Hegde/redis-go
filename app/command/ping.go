package command

func handlePing(context *Context) error {
	return context.Writer.WriteSimpleString("PONG")
}
