package handlers

import (
	"os"
	"strings"
)

func HandleCompletions(partialString string) string {
	seen := make(map[string]bool)
	var completions []string

	for cmd := range Builtins {
		if strings.HasPrefix(cmd, partialString) {
			seen[cmd] = true
			completions = append(completions, cmd)
		}
	}

	for _, name := range handlePathCompletions(partialString) {
		if !seen[name] {
			seen[name] = true
			completions = append(completions, name)
		}
	}

	if len(completions) == 1 {
		return completions[0]
	}
	return ""
}

func getPathDirs() []string {
	pathEnv := os.Getenv("PATH")
	return strings.Split(pathEnv, string(os.PathListSeparator))
}

func handlePathCompletions(partialString string) []string {
	completions := []string{}
	seen := make(map[string]bool)
	pathDirs := getPathDirs()
	for _, dir := range pathDirs{
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			name := entry.Name()
			if !strings.HasPrefix(name, partialString) {
				continue
			}
			if entry.IsDir() {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				continue
			}
			if info.Mode().Perm()&0111 != 0 {
				seen[name] = true
				completions = append(completions, name)	
			}
		}
	}
	return completions
}
