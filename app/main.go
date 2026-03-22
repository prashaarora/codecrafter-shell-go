package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/handlers"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
)

const (
	shellPrompt = "$ "
)

type Completer struct {
	lastPartial string
	tabPressCount int
}

func (c *Completer) Do(line []rune, pos int) ([][]rune, int) {
	lineStr := string(line[:pos])
	words := strings.Fields(lineStr)
	if len(words) == 0 {
		return nil, 0
	}
	partial := words[len(words)-1]
	completion := handlers.HandleAllCompletions(partial)
	switch len(completion) {
	case 0:
		fmt.Print("\a")
		c.lastPartial = ""
		c.tabPressCount = 0
		return nil, 0
	case 1:
		suffix := completion[0][len(partial):]
		c.lastPartial = ""
		c.tabPressCount = 0
		return [][]rune{[]rune(suffix + " ")}, len(partial)
	default:
		if c.lastPartial != partial || c.tabPressCount == 0 {
			fmt.Print("\a")
			c.lastPartial = partial
			c.tabPressCount = 1
			return nil, 0
		}
		if c.tabPressCount == 1 && c.lastPartial == partial {
			sort.Strings(completion)
			fmt.Printf("\n%s\n", strings.Join(completion, "  "))
			fmt.Printf("%s%s", shellPrompt, partial)
			c.tabPressCount = 0
			c.lastPartial = ""
			return nil, 0
		}
	}
	return nil, 0
}

func main() {
	completer := &Completer{}
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       shellPrompt,
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

		if h, ok := handlers.Registry[cmd]; ok {
			h.Execute(input)
		} else {
			handlers.HandleExternal(input)
		}
	}
}
