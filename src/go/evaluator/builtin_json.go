package evaluator

import (
	"z/lexer"
	"z/object"
	"z/parser"
)

func init_builtin_json_decode() *object.Builtin {
	jsonDecode := &object.Builtin{Fn: func(args ...object.Object) object.Object {
		if args[0].Type() != object.STRING_OBJ {
			return newError("argument 1 to `mysql_init` must be String, got=%s", args[0].Type())
		}
		jsonObj, _ := args[0].(*object.String)
		l := lexer.New(jsonObj.Value)
		p := parser.New(l)
		env := object.NewEnvironment()
		program := p.ParseProgram()
		result := Eval(program, env)
		return result
	},
	}
	return jsonDecode
}
