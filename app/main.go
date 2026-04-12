package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/subrahmanyahegde/redis-go/app/resp"
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

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error while accepting the connection", err.Error())
			os.Exit(1)
		}
		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
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
			writer.WriteError(err.Error())
			return
		}
		commands := value.Array
		switch commands[0].String {
		case "PING":
			writer.WriteSimpleString("PONG")
		case "ECHO":
			writer.WriteBulkString(commands[1].String)
		default:
			writer.WriteError("unrecognized command: " + commands[0].String)
		}
	}
}
