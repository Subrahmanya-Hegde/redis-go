package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

const (
	STRING      = '+'
	ERROR       = '-'
	INTEGER     = ':'
	BULK_STRING = '$'
	ARRAY       = '*'
)

var _ = net.Listen
var _ = os.Exit
var ErrLineTooShort = errors.New("line too short")
var ErrorResponse = "-ERR %s\n"
var SuccessResponse = "$%d\r\n%s\r\n"

func main() {
	fmt.Println("Local: http://0.0.0.0:6379")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error while accepting the connection.", err.Error())
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
	reader := bufio.NewReader(connection)

	for {
		commands, err := getCommandsArray(reader)
		if err != nil {
			if err == io.EOF {
				fmt.Print("client disconnected.")
				return
			}
			writeToConnection(connection, fmt.Sprintf(ErrorResponse, err.Error()))
			return
		}
		commandsLength := len(commands)
		if commandsLength == 0 {
			fmt.Println("No commands received, closing connection.")
			return
		}
		if commandsLength > 2 {
			writeToConnection(connection, fmt.Sprintf(ErrorResponse, "more than 2 commands received."))
			return
		}

		if commandsLength == 1 {
			if commands[0] != "PING" {
				writeToConnection(connection, fmt.Sprintf(ErrorResponse, "unrecognized command"))
				return
			}
			writeToConnection(connection, fmt.Sprintf("+PONG\r\n"))
			return
		}
		if commands[0] != "ECHO" {
			writeToConnection(connection, fmt.Sprintf(ErrorResponse, "unrecognized command."))
			return
		}
		writeToConnection(connection, fmt.Sprintf(SuccessResponse, len(commands[1]), commands[1]))
	}
}

func writeToConnection(connection net.Conn, line string) {
	_, _ = connection.Write([]byte(line))
}

// Sample RESP (REdis Serialization Protocol): *2\r\n$3\r\nGET\r\n$3\r\nkey\r\n
// * for array, $ for bulk string, + for simple string, : for integer, - for error
func getCommandsArray(reader *bufio.Reader) ([]string, error) {
	fmt.Println("Reading commands from connection...")

	firstByte, err := reader.ReadByte()
	var commands []string
	if err != nil {
		if err == io.EOF {
			return commands, err
		}
		fmt.Printf("Error while reading first byte: %s\n", err.Error())
		return commands, err
	}
	if firstByte != ARRAY {
		return commands, fmt.Errorf("expected array, got %c", firstByte)
	}
	line, err := readLine(reader)
	if err != nil {
		return commands, err
	}
	commandsLength, err := strconv.Atoi(line)
	if err != nil {
		return commands, err
	}
	commands = make([]string, commandsLength)

	for i := range commandsLength {
		command, err := readBulkString(reader)
		if err != nil {
			return commands, err
		}
		commands[i] = command
	}
	return commands, nil
}

func readBulkString(reader *bufio.Reader) (string, error) {
	bulkStringByte, err := reader.ReadByte()
	if err != nil {
		return "", err
	}
	if bulkStringByte != BULK_STRING {
		return "", fmt.Errorf("expected bulk string, got %c", bulkStringByte)
	}
	line, err := readLine(reader)
	if err != nil {
		return "", err
	}
	bulkStringLength, err := strconv.Atoi(line)
	if err != nil {
		return "", err
	}
	bulkStringBytes := make([]byte, bulkStringLength)
	_, err = io.ReadFull(reader, bulkStringBytes)
	if err != nil {
		return "", err
	}
	_, err = readLine(reader)
	if err != nil {
		if errors.Is(err, ErrLineTooShort) {
			return "", fmt.Errorf(fmt.Sprintf("expected at least %d bytes.", bulkStringLength))
		}
		return "", err
	}
	return string(bulkStringBytes), nil
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading line: %s", err.Error())
	}
	lineLength := len(line)
	if lineLength < 2 {
		return "-ERR", ErrLineTooShort
	}
	return line[:lineLength-2], nil //Line without \r\n
}
