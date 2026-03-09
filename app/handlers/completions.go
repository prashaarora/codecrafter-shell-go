package handlers

import (
	"strings"
)

func HandleCompletions(partialString string) string {
	var completions []string
	for cmd := range Builtins {
		if strings.HasPrefix(cmd, partialString) {
			completions = append(completions, cmd)
		}
	}

	if len(completions) == 1 {
		return completions[0]
	}
	return ""
}
