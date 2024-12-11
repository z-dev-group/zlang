package object

import (
	"os"
)

func init() {
	Builtins = append(Builtins, filePutContent())
	Builtins = append(Builtins, fileGetContent())
}

func filePutContent() BuiltinFn {
	return BuiltinFn{
		"file_put_contents",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != STRING_OBJ {
				return newError("argument 1 must be String, got=%s", args[0].Type())
			}
			filePath := args[0].(*String).Value
			if args[1].Type() != STRING_OBJ {
				return newError("argument 2 must be string. got=%s", args[1].Type())
			}
			contentStr := args[1].(*String).Value
			os.WriteFile(filePath, []byte(contentStr), 0644)
			return &Boolean{Value: true}
		}},
	}
}

func fileGetContent() BuiltinFn {
	return BuiltinFn{
		"file_get_contents",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != STRING_OBJ {
				return newError("argument 1 must be String, got=%s", args[0].Type())
			}
			filePath := args[0].(*String).Value
			contents, err := os.ReadFile(filePath)
			contentStr := ""
			if err == nil {
				contentStr = string(contents)
			}
			return &String{Value: contentStr}
		}},
	}
}
