package main

import (
	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/handlers"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
	"os"
	"strings"
)

type Completer struct{}

func (c *Completer) Do(line []rune, pos int) ([][]rune, int) {
	lineStr := string(line[:pos])
	words := strings.Fields(lineStr)
	if len(words) == 0 {
		return nil, 0
	}

	partial := words[len(words)-1]
	completion := handlers.HandleCompletions(partial)

	if completion != "" {
		suffix := completion[len(partial):]
		return [][]rune{[]rune(suffix + " ")}, len(partial)
	}

	return nil, 0
}

func main() {
	completer := &Completer{}
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "$ ",
		AutoComplete: completer,
	})
	if err != nil {
		os.Exit(1)
	}
	defer rl.Close()

	for {
		command, err := rl.Readline()
		if err != nil {
			os.Exit(0)
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
