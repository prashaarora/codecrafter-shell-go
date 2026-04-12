package parser

import (
	"strings"
)

// CommandParser incrementally parses a shell command string into arguments.
type CommandParser struct {
	result              []string
	currentArg          strings.Builder
	activeQuote         rune
	inEscape            bool
	inDoubleQuoteEscape bool
}

// NewCommandParser creates a parser instance with empty state.
func NewCommandParser() *CommandParser {
	return &CommandParser{
		result:              make([]string, 0),
		activeQuote:         0,
		inEscape:            false,
		inDoubleQuoteEscape: false,
	}
}

// ParseCommand parses a command string into shell-style arguments.
func ParseCommand(cmd string) []string {
	p := NewCommandParser()
	return p.Parse(cmd)
}

func (p *CommandParser) reset() {
	p.result = p.result[:0]
	p.currentArg.Reset()
	p.activeQuote = 0
	p.inEscape = false
	p.inDoubleQuoteEscape = false
}

// Parse tokenizes cmd and returns the parsed arguments.
func (p *CommandParser) Parse(cmd string) []string {
	p.reset()
	for _, c := range cmd {
		p.processCharacter(c)
	}

	if p.currentArg.Len() > 0 {
		p.result = append(p.result, p.currentArg.String())
	}

	return p.result
}

func (p *CommandParser) processCharacter(c rune) {
	if p.handleEscapedChar(c) {
		return
	}

	if p.handleBackslash(c) {
		return
	}

	if p.handleQuote(c) {
		return
	}

	if p.handleSpecialChars(c) {
		return
	}

	p.currentArg.WriteRune(c)
}

func (p *CommandParser) handleEscapedChar(c rune) bool {
	if p.inEscape {
		p.currentArg.WriteRune(c)
		p.inEscape = false
		return true
	}
	if p.inDoubleQuoteEscape {
		if c == '"' || c == '\\' {
			p.currentArg.WriteRune(c)
		} else {
			p.currentArg.WriteRune('\\')
			p.currentArg.WriteRune(c)
		}
		p.inDoubleQuoteEscape = false
		return true
	}

	return false
}

func (p *CommandParser) handleBackslash(c rune) bool {
	if c != '\\' {
		return false
	}

	if p.activeQuote == 0 {
		p.inEscape = true
		return true
	}
	if p.activeQuote == '"' {
		p.inDoubleQuoteEscape = true
		return true
	}
	return false
}

func (p *CommandParser) handleQuote(c rune) bool {
	if c == '\'' || c == '"' {
		switch p.activeQuote {
		case 0:
			p.activeQuote = c
		case c:
			p.activeQuote = 0
		default:
			p.currentArg.WriteRune(c)
		}
		return true
	}
	return false
}

func (p *CommandParser) handleSpecialChars(c rune) bool {
	if c == '\n' {
		return true
	}
	if c == ' ' {
		if p.activeQuote != 0 {
			p.currentArg.WriteRune(c)
		} else {
			if p.currentArg.Len() > 0 {
				p.result = append(p.result, p.currentArg.String())
				p.currentArg.Reset()
			}
		}
		return true
	}
	return false
}
