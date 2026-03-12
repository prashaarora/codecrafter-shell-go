package handlers

type Handler interface {
	Execute(input []string)
}

type HandlerFunc func(input []string)

func (f HandlerFunc) Execute(input []string) { f(input) }

// Registry maps command names to their Handler.
var Registry = map[string]Handler{
	"exit": HandlerFunc(HandleExit),
	"echo": HandlerFunc(HandleEcho),
	"type": HandlerFunc(HandleType),
	"pwd":  HandlerFunc(HandlePwd),
	"cd":   HandlerFunc(HandleCd),
}
