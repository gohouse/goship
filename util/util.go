package util

import (
	"flag"
	"os"
	"os/exec"
	"strings"
)

func NewFileWithMod(fileName string, mod int) *os.File {
	f, err := os.OpenFile(fileName, mod, 0766)
	if err != nil {
		panic(err.Error())
	}
	return f
}

func CamelCase(str string) string {
	var text string
	//for _, p := range strings.Split(name, "_") {
	for _, p := range strings.Split(str, "_") {
		// 字段首字母大写的同时, 是否要把其他字母转换为小写
		switch len(p) {
		case 0:
		case 1:
			text += strings.ToUpper(p[0:1])
		default:
			text += strings.ToUpper(p[0:1]) + p[1:]
		}
	}
	return text
}

func GofmtW(file string)  {
	cmd := exec.Command("sh","-c", "gofmt -w "+file)
	cmd.Run()
}

func GoModTidy()  {
	cmd := exec.Command("sh","-c", "go mod tidy")
	cmd.Run()
}

func SwagInit()  {
	cmd := exec.Command("sh","-c", "swag")
	cmd.Run()
}

func RunCmd(cmdstr string) error {
	cmd := exec.Command("sh","-c", cmdstr)
	return cmd.Run()
}

func RunCmds(cmdstr []string) error {
	cmd := exec.Command("sh","-c", strings.Join(cmdstr, " && "))
	return cmd.Run()
}
// CheckCommandExists 检查命令是否存在
func CheckCommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func IsFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}