package util

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var (
	cmdFunc = map[string]func(){
		"linux": func() {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
		},
		"darwin": func() {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
		},
		"windows": func() {
			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
		},
	}
)

func clearScreen() {
	if f, ok := cmdFunc[runtime.GOOS]; !ok {
		panic("Your platform is unsupported!")
	} else {
		f()
	}
}

func ClearPrint(param string) {
	clearScreen()
	fmt.Println(param)
}

func ClearPrintWithGap(param string) {
	clearScreen()
	fmt.Printf("\n%v\n", param)
}
