package handlers

var Builtins map[string]bool

func init() {
	Builtins = make(map[string]bool, len(Registry))
	for name := range Registry {
		Builtins[name] = true
	}
}
