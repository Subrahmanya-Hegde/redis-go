package command

import (
	"github.com/subrahmanyahegde/redis-go/app/resp"
	"github.com/subrahmanyahegde/redis-go/app/storage"
)

func handleRPush(context *Context) error {
	if len(context.Args) < 2 {
		return context.Writer.WriteError("RPUSH expects at least 2 arguments")
	}

	key := context.Args[0].String
	existingKey, ok := context.Store.Get(key)
	if ok && existingKey.Type != storage.TypeList {
		return context.Writer.WriteError("Wrong type operation against a key holding the wrong type")
	}

	values := context.Args[1:]
	var list []string
	list = append(list, getListElements(values)...)

	context.Store.Set(key, storage.Data{
		Type: storage.TypeList,
		List: list,
	})
	return context.Writer.WriteInteger(len(list))
}

func getListElements(values []resp.Value) []string {
	newItems := make([]string, len(values))
	for i, value := range values {
		newItems[i] = value.String
	}
	return newItems
}
