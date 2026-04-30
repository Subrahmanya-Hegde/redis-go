package command

import (
	"strings"

	"github.com/subrahmanyahegde/redis-go/app/resp"
	"github.com/subrahmanyahegde/redis-go/app/storage"
)

type Handler func(context *Context) error

const (
	GET   = "GET"
	SET   = "SET"
	ECHO  = "ECHO"
	PING  = "PING"
	RPUSH = "RPUSH"
)

var registry = map[string]Handler{
	GET:   handleGet,
	SET:   handleSet,
	ECHO:  handleEcho,
	PING:  handlePing,
	RPUSH: handleRPush,
}

type Context struct {
	Command string
	Args    []resp.Value
	Writer  *resp.Writer
	Store   *storage.Store
}

func Execute(context *Context) error {
	handler, ok := registry[strings.ToUpper(context.Command)]
	if !ok {
		return context.Writer.WriteError("unknown command: " + context.Command)
	}
	return handler(context)
}
