package builtins

// Builtin describes a built-in URScript function.
type Builtin struct {
	Name      string
	Signature string
	Doc       string
	Module    string
}

// Registry holds the built-in function definitions.
type Registry struct {
	funcs map[string]Builtin
}

// NewRegistry creates a registry pre-populated with standard UR built-ins.
func NewRegistry() *Registry {
	r := &Registry{funcs: make(map[string]Builtin)}
	groups := [][]Builtin{
		motionBuiltins(),
		internalsBuiltins(),
		urmathBuiltins(),
		interfacesBuiltins(),
		ioconfigurationBuiltins(),
		processpathBuiltins(),
	}
	for _, g := range groups {
		for _, b := range g {
			r.funcs[b.Name] = b
		}
	}
	return r
}

// Lookup returns the built-in with the given name, if any.
func (r *Registry) Lookup(name string) (Builtin, bool) {
	b, ok := r.funcs[name]
	return b, ok
}

// All returns all built-in functions.
func (r *Registry) All() []Builtin {
	out := make([]Builtin, 0, len(r.funcs))
	for _, b := range r.funcs {
		out = append(out, b)
	}
	return out
}
