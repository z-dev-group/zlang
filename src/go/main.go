package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"z/cli"
	"z/repl"
)

func main() {
	var operation string
	var fileName string
	mode := "vm"
	if len(os.Args) > 1 {
		operation = os.Args[1]
		if operation == "run" || operation == "build" {
			if len(os.Args) == 2 {
				fmt.Println("please input source file")
				return
			}
			fileName = os.Args[2]
			if len(os.Args) == 4 {
				mode = os.Args[3]
			}
		} else {
			operation = "run"
			fileName = os.Args[1]
			if len(os.Args) == 3 {
				mode = os.Args[2]
			}
		}
		fileContent, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
		sourceCode := string(fileContent)
		sourceCodeLines := strings.Split(sourceCode, "\n")
		if len(sourceCodeLines) > 0 {
			dir, _ := exec.LookPath(os.Args[0])
			buildSourceDir := dir
			envZRoot, ok := os.LookupEnv("Z_ROOT")
			if ok {
				buildSourceDir = envZRoot
			}
			builtinLine := `import "` + buildSourceDir + `/standard/builtin.z";`
			if strings.Contains(sourceCodeLines[0], "package") {
				runSourceCodeLines := []string{}
				runSourceCodeLines = append(runSourceCodeLines, sourceCodeLines[0])
				runSourceCodeLines = append(runSourceCodeLines, builtinLine)
				runSourceCodeLines = append(runSourceCodeLines, sourceCodeLines[1:]...)
				sourceCode = strings.Join(runSourceCodeLines, "\n")
			} else {
				sourceCode = builtinLine + sourceCode
			}
		}
		switch operation {
		case "run":
			cli.RunSourceCode(sourceCode, mode, fileName)
			return
		case "build":
			cli.BuildSourceCode(sourceCode, fileName)
			return
		}
	}
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Println("Welcome to z language , user is: ", user.Username, " type code to execute")
	repl.CStart(os.Stdin, os.Stdout)
}
