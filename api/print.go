package api

import "fmt"

const (
	RED    = 31
	GREEN  = 32
	YELLOW = 33
	CYAN   = 36
)

func Err(msg string) {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", RED, msg)
	fmt.Print(colored)
}

func Stat(msg string) {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", GREEN, msg)
	fmt.Print(colored)
}

func Broadcast(msg string) {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", CYAN, msg)
	fmt.Print(colored)
}

func Log(msg string) {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", YELLOW, msg)
	fmt.Print(colored)
}

func Println(msg string) {
	fmt.Println(msg)
}

func Print(msg string) {
	fmt.Print(msg)
}
