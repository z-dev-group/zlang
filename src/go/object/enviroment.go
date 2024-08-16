package object

func NewEnclosedEnviroment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string, packageName string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && packageName != "" {
		varName := packageName + "." + name
		obj, ok = e.store[varName]
	}
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name, packageName)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object, packageName string) Object {
	if packageName != "" {
		name = packageName + "." + name
	}
	e.store[name] = val
	return val
}

func (e *Environment) GetAll() map[string]Object {
	return e.store
}
