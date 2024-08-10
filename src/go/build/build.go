package build

import (
	"fmt"
	"z/ast"
	"z/evaluator"
	"z/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) (object.Object, string) {
	switch node := node.(type) {
	case *ast.LetStatement:
		object, val := Eval(node.Value, env)
		env.Set(node.Name.Value, object)
		objectType := object.Type()
		switch objectType {
		case "STRING":
			return nil, "char *" + node.Name.Value + "=" + val + ";\n"
		case "INTEGER":
			return nil, "int " + node.Name.Value + "=" + val + ";\n"
		case "FLOAT":
			return nil, "double " + node.Name.Value + "=" + val + ";\n"
		}
	case *ast.Identifier:
		object := evalIdentifier(node, env)
		return object, node.Value
	case *ast.FunctionLiteral:
		return nil, node.Name
	case *ast.CallExpression:
		_, function := Eval(node.Function, env)
		callString := ""
		switch function {
		case "puts":
			for _, argument := range node.Arguments {
				object, argumentRes := Eval(argument, env)
				if object.Type() == "STRING" {
					callString += "printf(\"%s\"," + argumentRes + ");\n"
				}
				if object.Type() == "INTEGER" {
					callString += "printf(\"%d\"," + argumentRes + ");\n"
				}
				if object.Type() == "FLOAT" {
					callString += "printf(\"%f\"," + argumentRes + ");\n"
				}
			}
		}
		return nil, callString
	case *ast.StringLiteral:
		return &object.String{}, string("\"" + node.Value + "\"")
	case *ast.IntegerLiteral:
		return &object.Integer{}, node.String()
	case *ast.FloatLiteral:
		return &object.Float{}, node.String()
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	}
	return nil, ""
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if builtin, ok := evaluator.Builtins[node.Value]; ok {
		return builtin
	}
	if !ok {
		fmt.Println("identifier not found:", node.Value)
	}
	return val
}
