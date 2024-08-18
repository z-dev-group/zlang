package object

func NewEnclosedEnviroment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	Context := make(map[string]string)
	return &Environment{store: s, outer: nil, Context: Context}
}

type Environment struct {
	store   map[string]Object
	Context map[string]string
	outer   *Environment
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

func (e *Environment) IsFormOuter(name string, packageName string) bool {
	isFromOuter := false
	queryName := name
	if packageName != "" {
		queryName = packageName + "." + name
	}
	_, ok := e.store[queryName]
	if !ok && e.outer != nil {
		_, ok = e.outer.Get(name, packageName)
		if ok {
			isFromOuter = true
		} else {
			return e.outer.IsFormOuter(name, packageName)
		}
	}
	return isFromOuter
}

func (e *Environment) Set(name string, val Object, packageName string) Object {
	if packageName != "" {
		name = packageName + "." + name
	}
	e.store[name] = val
	return val
}

func (e *Environment) OuterSet(name string, val Object, packageName string) Object {
	if e.outer == nil {
		return newError("outer is not exists")
	}
	queryName := name
	if packageName != "" {
		queryName = packageName + "." + name
	}
	_, ok := e.outer.store[queryName]
	if ok {
		e.outer.store[queryName] = val
	} else {
		return e.outer.OuterSet(name, val, packageName)
	}
	return val
}

func (e *Environment) GetAll() map[string]Object {
	return e.store
}

func (e *Environment) Outer() *Environment {
	if e.outer != nil {
		return e.outer
	}
	return e
}
