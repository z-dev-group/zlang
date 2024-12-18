package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
	"z/ast"
	"z/code"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
	Json() string
}

const (
	INTEGER_OBJ              = "INTEGER"
	BOOLEAN_OBJ              = "BOOLEAN"
	NULL_OBJ                 = "NULL"
	ERROR_OBJ                = "ERROR"
	RETURN_VALUE_OBJ         = "RETURN_VALUE"
	FUNCTION_OBJ             = "FUNCTION"
	STRING_OBJ               = "STRING"
	BUILTIN_OBJ              = "BUILTIN"
	ARRAY_OBJ                = "ARRAY"
	HASH_OBJ                 = "OBJECT" // hash_obj as object
	COMPILED_FUNCTION_OBJECT = "COMPILED_FUNCTION_OBJECT"
	CLOSURE_OBJ              = "CLOSURE"
	FLOAT_OBJ                = "FLOAT"
	INTERFACE_OBJ            = "INTERFACE"
	CLASS_OBJ                = "CLASS"
	OBJECT_INSTANCE          = "OBJECT"
)

type Integer struct {
	Value int64
	Error *Error
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Json() string     { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Boolean struct {
	Value bool
	Error *Error
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Json() string     { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

type Null struct {
}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Json() string     { return "\"null\"" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) Json() string     { return "\"error:" + e.Message + "\"" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Json() string     { return rv.Value.Json() }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Name       string
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(("} {\n"))
	out.WriteString(f.Body.String())
	out.WriteString("\n")

	return out.String()
}
func (rv *Function) Json() string { return "\"function\"" }

type String struct {
	Value string
	Error *Error
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) Json() string     { return "\"" + s.Value + "\"" }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type BuiltinFunction = func(args ...Object) Object
type Builtin struct {
	Fn       BuiltinFunction
	FilePath string
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Json() string     { return "builtin function" }

type Array struct {
	Elements []Object
	Error    *Error
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}

	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
func (ao *Array) Json() string {
	var out bytes.Buffer
	elements := []string{}

	for _, e := range ao.Elements {
		elements = append(elements, e.Json())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
	Index int8
}

type Hash struct {
	Pairs    map[HashKey]HashPair
	Error    *Error
	MaxIndex int8
}

func (h *Hash) Type() ObjectType {
	return HASH_OBJ
}
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

func (h *Hash) Json() string {
	var out bytes.Buffer
	pairs := []string{}
	maxIndex := h.MaxIndex
	for i := 1; i <= int(maxIndex); i++ {
		for _, pair := range h.Pairs {
			if pair.Index == int8(i) {
				pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Json(), pair.Value.Json()))
			}
		}
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type Hashable interface {
	HashKey() HashKey
}

type CompiledFunction struct {
	Instructions  code.Instructions
	NumLocals     int
	NumParameters int
}

func (cf *CompiledFunction) Type() ObjectType { return COMPILED_FUNCTION_OBJECT }
func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
}
func (cf *CompiledFunction) Json() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
}

type Closure struct {
	Fn   *CompiledFunction
	Free []Object
}

func (c *Closure) Type() ObjectType { return CLOSURE_OBJ }
func (c *Closure) Inspect() string {
	return fmt.Sprintf("Closure[%p]", c)
}
func (c *Closure) Json() string {
	return fmt.Sprintf("Closure[%p]", c)
}

type Float struct {
	Value float64
	Error *Error
}

func (f *Float) Inspect() string  { return fmt.Sprintf("%v", f.Value) }
func (f *Float) Json() string     { return fmt.Sprintf("%v", f.Value) }
func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) HashKey() HashKey {
	return HashKey{Type: f.Type(), Value: uint64(f.Value)}
}

type Interface struct {
	Name      string
	Functions []*Function
}

func (i *Interface) Inspect() string  { return fmt.Sprintf("interface %s", i.Name) }
func (i *Interface) Json() string     { return fmt.Sprintf("interface %s", i.Name) }
func (i *Interface) Type() ObjectType { return CLASS_OBJ }

type Class struct {
	Name        string
	Parents     []*Class
	Interface   *Interface
	Environment *Environment
}

func (c *Class) Inspect() string  { return fmt.Sprintf("interface %s", c.Name) }
func (c *Class) Json() string     { return fmt.Sprintf("interface %s", c.Name) }
func (c *Class) Type() ObjectType { return CLASS_OBJ }

type ObjectInstance struct {
	InstanceClass *Class
	Environment   *Environment
}

func (oi *ObjectInstance) Inspect() string {
	return fmt.Sprintf("object %s", oi.InstanceClass.Name)
}
func (oi *ObjectInstance) Json() string     { return fmt.Sprintf("object %s", oi.InstanceClass.Name) }
func (oi *ObjectInstance) Type() ObjectType { return OBJECT_INSTANCE }
