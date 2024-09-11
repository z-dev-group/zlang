package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var cmd *exec.Cmd

var sourceFileMap = make(map[string]time.Time)

func RunDev() {
	if len(os.Args) < 3 {
		return
	}
	ch := make(chan int)
	cmd = exec.Command(os.Args[0], "run", os.Args[2])
	go runCommand(cmd)
	go watchDir(cmd)

	fmt.Println("exit", <-ch)
}

func watchDir(cmd *exec.Cmd) {
	if len(os.Args) < 4 {
		return
	}
	currentDir, _ := os.Getwd()
	watchDir := currentDir + "/" + os.Args[3]
	fmt.Println("watch dir is: " + watchDir)
	for {
		time.Sleep(1 * time.Second)
		isFileChanged := checkFileIsChange(watchDir)
		if isFileChanged {
			fmt.Println("file changed restart program")
			err := cmd.Process.Kill()
			if err != nil {
				fmt.Println("exit error:", err.Error())
			} else {
				fmt.Println("stop process successed")
				fmt.Println("start new process")
				cmd = exec.Command(os.Args[0], "run", os.Args[2])
				go runCommand(cmd)
			}
		}
	}
}

func checkFileIsChange(dir string) bool {
	isChanged := false
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	for _, file := range files {
		fileName := filepath.Join(dir, file.Name())
		if file.Type().IsRegular() {
			fileInfo, _ := file.Info()
			oldTime, ok := sourceFileMap[fileName]
			if !ok {
				sourceFileMap[fileName] = fileInfo.ModTime()
			} else {
				if oldTime.Before(fileInfo.ModTime()) {
					fmt.Println(fileName + " ---> change.....")
					isChanged = true
					sourceFileMap[fileName] = fileInfo.ModTime()
				}
			}
		} else {
			checkDirFileIsChanged := checkFileIsChange(fileName)
			if !isChanged {
				isChanged = checkDirFileIsChanged
			}
		}

	}
	return isChanged
}
func runCommand(cmd *exec.Cmd) {
	stdout, _ := cmd.StdoutPipe()
	go func() {
		// 读取标准输出
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Println(err)
				}
				break
			}
			fmt.Print(string(buf[:n]))
		}
	}()
	err := cmd.Start()
	if err != nil {
		fmt.Println(err.Error())
	}
}
