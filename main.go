package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Formatter struct {
	*sync.WaitGroup
}

var d string

// 格式化某个文件夹下所有Go文件的代码
func main() {

	flag.StringVar(&d, "d", "", "directory you need to format")
	flag.Parse()
	if d == "" {
		fmt.Println("param d can't be null")
		return
	}
	formatter := &Formatter{WaitGroup: &sync.WaitGroup{}}
	formatter.formatDir(d)
	formatter.Wait()
}

func (f *Formatter) getFullPath(path, fileName string) string {
	if strings.HasSuffix(path, "/") {
		return path + fileName
	} else {
		return path + "/" + fileName
	}
}

func (f *Formatter) formatDir(path string) {
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("open directory error,", err.Error())
		return
	}
	for _, fi := range fileInfo {
		f.Add(1)
		go func(file os.FileInfo) {
			defer f.Done()
			fileName := f.getFullPath(path, file.Name())
			if file.IsDir() {
				if file.Name() != "vendor" {
					f.formatDir(fileName)
				}
			} else {
				if strings.HasSuffix(file.Name(), "go") {
					fmt.Println("format go file :", fileName)
					f.goFormatFile(fileName)
				}
			}
		}(fi)
	}
}

func (f *Formatter) goFormatFile(fileName string) {
	cmd := exec.Command("go", "fmt", fileName)
	cmd.Stdout = os.Stdout
	cmd.Run()
}
