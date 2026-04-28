package command

import (
	"strconv"
	"strings"
	"time"

	"github.com/subrahmanyahegde/redis-go/app/storage"
)

func handleSet(context *Context) error {
	commandArguments := len(context.Args)
	if commandArguments < 2 {
		return context.Writer.WriteError("SET expects at least 2 arguments")
	}
	var expiry time.Time
	key := context.Args[0].String
	value := context.Args[1].String

	i := 2
	for i < commandArguments {
		arg := strings.ToUpper(context.Args[i].String)
		switch arg {
		case "EX":
			if (i + 1) >= commandArguments {
				return context.Writer.WriteError("EX expects at least 2 arguments")
			}
			seconds, err := strconv.Atoi(context.Args[i+1].String)
			if err != nil {
				return context.Writer.WriteError(err.Error())
			}
			expiry = time.Now().Add(time.Duration(seconds) * time.Second)
			i += 2
		case "PX":
			if (i + 1) >= commandArguments {
				return context.Writer.WriteError("PX expects at least 2 arguments")
			}
			ms, err := strconv.Atoi(context.Args[i+1].String)
			if err != nil {
				return context.Writer.WriteError(err.Error())
			}
			expiry = time.Now().Add(time.Duration(ms) * time.Millisecond)
			i += 2
		case "NX":
			i += 1
		case "XX":
			i += 1
		default:
			return context.Writer.WriteError("unknown argument " + arg)
		}
	}
	context.Store.Set(key, storage.Data{
		String: value,
		Type:   "string",
		Expiry: expiry,
	})
	return context.Writer.WriteSimpleString("OK")
}
