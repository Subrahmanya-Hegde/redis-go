package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/subrahmanyahegde/redis-go/app/command"
	"github.com/subrahmanyahegde/redis-go/app/resp"
	"github.com/subrahmanyahegde/redis-go/app/storage"
)

var _ = net.Listen
var _ = os.Exit

func main() {
	fmt.Println("Local: http://0.0.0.0:6379")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("Failed to close listener")
		}
	}(listener)
	store := storage.NewStore()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error while accepting the connection", err.Error())
			os.Exit(1)
		}
		go handleConnection(connection, store)
	}
}

func handleConnection(connection net.Conn, store *storage.Store) {
	defer func(connection net.Conn) {
		err := connection.Close()
		if err != nil {
			fmt.Printf("Error while closing connection: %s\n", err.Error())
		}
	}(connection)

	reader := resp.NewReader(connection)
	writer := resp.NewWriter(connection)

	for {
		value, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by remote host")
				return
			}
			_ = writer.WriteError(err.Error())
			return
		}

		commandAndArgs := value.Array
		context := &command.Context{
			Args:    commandAndArgs[1:],
			Command: commandAndArgs[0].String,
			Writer:  writer,
			Store:   store,
		}
		err = command.Execute(context)
		if err != nil {
			_ = writer.WriteError(err.Error())
		}
	}
}
