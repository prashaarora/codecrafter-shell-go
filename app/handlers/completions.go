package handlers

import (
	"os"
	"strings"
)

func HandleCompletions(partialString string) string {
	var completions []string
	for cmd := range Builtins {
		if strings.HasPrefix(cmd, partialString) {
			completions = append(completions, cmd)
		}
	}
	pathMatches := handlePathCompletions(partialString)
	completions = append(completions, pathMatches...)
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
