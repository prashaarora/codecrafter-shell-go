package handlers

// Handler defines the contract for shell command handlers.
type Handler interface {
	Execute(input []string)
}

// HandlerFunc adapts a function to implement the Handler interface.
type HandlerFunc func(input []string)

// Execute runs the adapted function.
func (f HandlerFunc) Execute(input []string) { f(input) }

// Registry maps command names to their Handler.
var Registry = map[string]Handler{
	"exit": HandlerFunc(HandleExit),
	"echo": HandlerFunc(HandleEcho),
	"type": HandlerFunc(HandleType),
	"pwd":  HandlerFunc(HandlePwd),
	"cd":   HandlerFunc(HandleCd),
}
