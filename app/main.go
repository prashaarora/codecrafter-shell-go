package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/handlers"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, handlers.MsgErrorReading, err)
			os.Exit(1)
		}
		input := parser.ParseCommand(command)
		if len(input) == 0 {
			continue
		}
		cmd := input[0]

		switch cmd {
		case "exit":
			handlers.HandleExit(input)
		case "echo":
			handlers.HandleEcho(input)
		case "type":
			handlers.HandleType(input)
		case "pwd":
			handlers.HandlePwd(input)
		case "cd":
			handlers.HandleCd(input)
		default:
			handlers.HandleExternal(input)

		}
	}
}
