package evaluator

import (
	"z/object"
)

var Builtins = map[string]*object.Builtin{
	"len":               object.GetBuiltinByName("len"),
	"push":              object.GetBuiltinByName("push"),
	"puts":              object.GetBuiltinByName("puts"),
	"execute":           object.GetBuiltinByName("execute"),
	"mysql_query":       object.GetBuiltinByName("mysql_query"),
	"mysql_init":        object.GetBuiltinByName("mysql_init"),
	"typeof":            object.GetBuiltinByName("typeof"),
	"fetch":             object.GetBuiltinByName("fetch"),
	"json_encode":       object.GetBuiltinByName("json_encode"),
	"with_error":        object.GetBuiltinByName("with_error"),
	"is_with_error":     object.GetBuiltinByName("is_with_error"),
	"get_error_message": object.GetBuiltinByName("get_error_message"),
	"syscall":           object.GetBuiltinByName("syscall"),
}
