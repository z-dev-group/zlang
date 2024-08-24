package evaluator

import (
	"encoding/json"
	"fmt"
	"z/ast"
	"z/object"
)

var (
	NULL         = &object.Null{}
	TRUE         = &object.Boolean{Value: true}
	FALSE        = &object.Boolean{Value: false}
	initedEnv    object.Environment
	withBreakKey = "is_with_break"
	isWithBreak  = "Y"
	notWithBreak = "N"
	breadKeyWord = "break"
)

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		initedEnv = *env
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		env.Set(node.Name.Value, val, node.PackageName)
		return val
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) && node.Operator != "->" {
			return right
		}
		if node.Operator == "->" { // object get not need eval
			right = &object.String{Value: node.Right.String()}
		}
		infixValue := evalInfixExpression(node.Operator, left, right)
		resetOperators := []string{
			"+=",
			"-=",
			"*=",
			"/=",
			"++",
			"--",
			"=",
		}
		if isInStringArray(resetOperators, node.Operator) { // need reset env data
			leftIdentifier, ok := node.Left.(*ast.Identifier)
			if ok {
				isFromOuter := env.IsFormOuter(leftIdentifier.Value, leftIdentifier.PackageName)
				if isFromOuter {
					env.OuterSet(leftIdentifier.Value, infixValue, leftIdentifier.PackageName)
				} else {
					env.Set(leftIdentifier.Value, infixValue, leftIdentifier.PackageName)
				}
			}
		}
		return infixValue
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.WhileExpression:
		return evalWhileExpression(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		function := &object.Function{Parameters: params, Env: env, Body: body}
		if node.Name != "" {
			env.Set(node.Name, function, node.PackageName)
		}
		return function
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)

		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.HashAssignExpress:
		hashObject, ok := env.Get(node.Hash.Value, node.Hash.PackageName)
		if !ok {
			return newError("hash variable " + node.Hash.Value + " not found")
		}
		index := Eval(node.Index, env)
		val := Eval(node.Value, env)
		hash, ok := hashObject.(*object.Hash)
		if !ok {
			return newError("object is not hash")
		}
		hashKey, ok := index.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s,", index.Type())
		}
		hashed := hashKey.HashKey()
		hash.Pairs[hashed] = object.HashPair{Key: index, Value: val}

		return val
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.ForExpression:
		return evalForExpression(node, env)
	case *ast.ClassExpress:
		return evalClassExpression(node, env)
	case *ast.ObjectExpress:
		return evalObjectExpression(node, env)
	}
	return nil
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s,", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpress(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s ", left.Type())
	}
}

func evalStringIndexExpress(str, index object.Object) object.Object {
	stringObject, _ := str.(*object.String)
	key, _ := index.(*object.Integer)
	idx := key.Value
	max := len(stringObject.Value)
	if idx < 0 || idx > int64(max) {
		return NULL
	}
	returnStringObject := object.String{}
	singleString := string(stringObject.Value[idx])
	returnStringObject.Value = singleString
	return &returnStringObject
}
func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)

	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]

	if !ok {
		return NULL
	}
	return pair.Value
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObj.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObj.Elements[idx]
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		if result := fn.Fn(args...); result != nil {
			return result
		}
		return NULL
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnviroment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx], "")
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluted := Eval(e, env)
		if isError(evaluted) {
			return []object.Object{evaluted}
		}
		result = append(result, evaluted)
	}
	return result
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	env = object.NewEnclosedEnviroment(env)
	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		consequenceVal := Eval(ie.Consequence, env)
		if env.Context[withBreakKey] == isWithBreak {
			env.Outer().Context[withBreakKey] = isWithBreak
		}
		return consequenceVal
	} else if ie.Alternative != nil {
		alternativeVal := Eval(ie.Alternative, env)
		if env.Context[withBreakKey] == isWithBreak {
			env.Outer().Context[withBreakKey] = isWithBreak
		}
		return alternativeVal
	} else {
		return NULL
	}
}

func evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
	env = object.NewEnclosedEnviroment(env)
	condition := Eval(we.Condition, env)
	if isError(condition) {
		return condition
	}
	env.Context[withBreakKey] = notWithBreak

	for isTruthy(condition) {
		bodyResult := Eval(we.Body, env)
		if env.Context[withBreakKey] == isWithBreak {
			break
		}
		condition := Eval(we.Condition, env)
		if !isTruthy(condition) {
			return bodyResult
		}
	}
	return NULL
}

func isTruthy(obj object.Object) bool {
	boolean, ok := obj.(*object.Boolean) // question ?
	if ok {
		return boolean.Value
	}
	switch obj {
	case NULL:
		return false
	default:
		return true
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case operator == "=":
		return right
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case operator == "->":
		return evalObjectGetInfixExpress(left, right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		fallthrough
	case "++":
		fallthrough
	case "+=":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		fallthrough
	case "--":
		fallthrough
	case "-=":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		fallthrough
	case "*=":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		fallthrough
	case "/=":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return NULL
	}
}

func evalFloatInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return NULL
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperationExpression(right)
	default:
		return newError("unkown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperationExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		if statement.TokenLiteral() == breadKeyWord {
			env.Context[withBreakKey] = isWithBreak
			return result
		}
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}

		}
	}
	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value, node.PackageName)

	if builtin, ok := Builtins[node.Value]; ok {
		return builtin
	}
	if node.Value == "http_server" {
		return init_builtin_http_server()
	}
	if !ok {
		return newError("identifier not found:%s", node.Value)
	}
	return val
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}
	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isInStringArray(array []string, findStr string) bool {
	isFind := false
	for _, str := range array {
		if str == findStr {
			isFind = true
			break
		}
	}
	return isFind
}

func evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
	env = object.NewEnclosedEnviroment(env)
	Eval(fe.Initor, env)
	condition := Eval(fe.Condition, env)
	if isError(condition) {
		return condition
	}
	env.Context[withBreakKey] = notWithBreak

	for isTruthy(condition) {
		bodyResult := Eval(fe.Body, env)
		Eval(fe.After, env)
		if env.Context[withBreakKey] == isWithBreak {
			break
		}
		condition := Eval(fe.Condition, env)
		if !isTruthy(condition) {
			return bodyResult
		}
	}
	return NULL
}

func evalClassExpression(ce *ast.ClassExpress, env *object.Environment) object.Object {
	classObject := &object.Class{
		Name:    ce.Name.Value,
		Parents: []*object.Class{},
	}
	classEnv := object.NewEnclosedEnviroment(env)
	for _, letStatement := range ce.LetStatements {
		value := Eval(letStatement, classEnv)
		classEnv.Set(letStatement.Name.Value, value, "")
	}
	for _, parent := range ce.Parents {
		parentObj, ok := env.Get(parent.Value, "")
		if !ok {
			return newError("parent class not exists")
		}
		parentClassObject, ok := parentObj.(*object.Class)
		if !ok {
			return newError("parent is not a class")
		}
		classObject.Parents = append(classObject.Parents, parentClassObject)
	}
	classObject.Environment = classEnv
	env.Set(ce.Name.Value, classObject, "")
	return classObject
}

func evalObjectExpression(oe *ast.ObjectExpress, env *object.Environment) object.Object {
	objectInstance := &object.ObjectInstance{}
	class, ok := env.Get(oe.Class.Value, "")
	objectEnv := object.NewEnclosedEnviroment(env)
	if ok {
		instanceClass, ok := class.(*object.Class)
		if ok {
			objectInstance.InstanceClass = instanceClass
			for _, parent := range instanceClass.Parents {
				parentClassObject, ok := env.Get(parent.Name, "")
				if ok {
					parentClass, ok := parentClassObject.(*object.Class)
					if ok {
						copyClassProperties(parentClass, objectEnv)
					}
				}
			}
			copyClassProperties(instanceClass, objectEnv)
			objectInstance.Environment = objectEnv
		}
	}
	return objectInstance
}

func copyClassProperties(class *object.Class, env *object.Environment) bool {
	clasEvnProperties := class.Environment.GetAll()
	for name, classProperty := range clasEvnProperties {
		newClassProperty := classProperty
		buffer, _ := json.Marshal(&classProperty)
		json.Unmarshal([]byte(buffer), newClassProperty)
		env.Set(name, newClassProperty, "")
	}
	return true
}

func evalObjectGetInfixExpress(left object.Object, right object.Object) object.Object {
	objectInstance, ok := left.(*object.ObjectInstance)
	if !ok {
		return newError("object is not object")
	}
	isStatic := false
	if isStatic {
		return nil
	}
	return getObjectInstanceValue(objectInstance, right)
}

func getObjectInstanceValue(objectInstance *object.ObjectInstance, right object.Object) object.Object {
	rightIdentifier, ok := right.(*object.String)
	if !ok {
		return newError("right is not string")
	}
	value, ok := objectInstance.Environment.Get(rightIdentifier.Value, "")
	if ok {
		functionValue, ok := value.(*object.Function)
		if ok {
			functionValue.Env = objectInstance.Environment
			return functionValue
		}
		return value
	}
	return nil
}
