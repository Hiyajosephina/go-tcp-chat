// Package api here provides methods that both client and server will use
package api

import (
	"fmt"
	"time"
)

const (
	Red     = 31
	Green   = 32
	Yellow  = 33
	Magenta = 35
	Cyan    = 36
)

// Err prints the string in red and proceeded by the current timestamp
func Err(msg string) {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", Red, time.Now().Format(time.Stamp)+": "+msg)
	fmt.Print(colored)
}

// Stat prints the string in green and proceeded by the current timestamp
func Stat(msg string) {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", Green, time.Now().Format(time.Stamp)+": "+msg)
	fmt.Print(colored)
}

// Broadcast prints the string in cyan
func Broadcast(msg string) {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", Cyan, msg)
	fmt.Print(colored)
}

// DirectMessage prints the string in magenta
func DirectMessage(msg string) {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", Magenta, msg)
	fmt.Print(colored)
}

// Log prints the string in yellow and proceeded by the current timestamp
func Log(msg string) {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", Yellow, time.Now().Format(time.Stamp)+": "+msg)
	fmt.Print(colored)
}

// Println prints the string
func Println(msg string) {
	fmt.Println(msg)
}

// Println prints the string
func Print(msg string) {
	fmt.Print(msg)
}
