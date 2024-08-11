package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"z/ast"
	"z/build"
	"z/lexer"
	"z/object"
	"z/parser"
)

func BuildSourceCode(sourceCode string, sourceFile string) {
	var duration time.Duration
	start := time.Now()
	fmt.Println("begin to build code")
	compiledCode := parseAstToC(sourceCode)
	fileName := filepath.Base(sourceFile)
	outFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	compileC(compiledCode, outFileName)
	duration = time.Since(start)
	fmt.Printf("build execute time is :%s\n", duration)
}

func parseAstToC(sourceCode string) string {
	l := lexer.New(sourceCode)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	compiledCode := generateProgram(program, env)
	return compiledCode
}

func ConvertZToC(sourceCode string, isWrapper bool) string {
	l := lexer.New(sourceCode)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	compiledCCode := ""
	if isWrapper {
		compiledCCode = generateProgram(program, env)
	} else {
		compiledCCode = generateCompiledCode(program, env)
	}
	return compiledCCode
}

func generateProgram(program *ast.Program, env *object.Environment) string {
	generatedCompiledCode := ""
	compiledCode := generateCompiledCode(program, env)
	generatedCompiledCode += "#include <stdio.h>\n"
	generatedCompiledCode += "#include <stdlib.h>\n"
	generatedCompiledCode += "int main() {\n"
	generatedCompiledCode += compiledCode + "\n"
	generatedCompiledCode += "}\n"
	return generatedCompiledCode
}

func generateCompiledCode(program *ast.Program, env *object.Environment) string {
	var compiledCode string
	for _, statement := range program.Statements {
		fmt.Println(statement)
		_, code := build.Eval(statement, env)
		compiledCode = compiledCode + code
	}
	return compiledCode
}

func compileC(code string, outFile string) {
	tempFileName := "./temp.c"
	file, err := os.Create(tempFileName)
	if err != nil {
		panic(err)
	}
	_, _ = file.WriteString(code)
	file.Close()
	cmd := exec.Command("gcc", tempFileName, "-o", outFile)
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("build success!")
		os.Remove(tempFileName)
	}
}
