package handlers

import (
	"fmt"
	"os"
	"os/exec"
)

func HandleType(input []string) {
	if len(input) < 2 {
		fmt.Fprintln(os.Stderr, "type"+MsgMissingArg)
		return
	}
	cmdName := input[1]

	if Builtins[cmdName] {
		fmt.Println(cmdName + MsgIsBuiltin)
	} else {
		path, err := exec.LookPath(cmdName)
		if err == nil {
			fmt.Println(cmdName + MsgIs + path)
		} else {
			fmt.Fprintln(os.Stderr, cmdName+MsgNotFound)
		}
	}
}
