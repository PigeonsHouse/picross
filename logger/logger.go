package logger

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type LogLebel int

const (
	Debug LogLebel = iota
	Info
	Warn
	Error
	Progress
)

type Logger struct {
	level LogLebel
}

var logger Logger

func Init(level LogLebel) {
	logger = Logger{
		level: level,
	}
}

func CleanTerminal() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls") // Windows
	} else {
		cmd = exec.Command("clear") // Linux, macOS
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func AppLogLevel() LogLebel {
	return logger.level
}

func DebugLog(message string) {
	if logger.level > Debug {
		return
	}
	message = strings.ReplaceAll(message, "\n", "\n        ")
	fmt.Println("[DEBUG]", message)
}

func InfoLog(message string) {
	if logger.level > Info {
		return
	}
	message = strings.ReplaceAll(message, "\n", "\n       ")
	fmt.Println("[INFO]", message)
}

func WarnLog(message string) {
	if logger.level > Warn {
		return
	}
	message = strings.ReplaceAll(message, "\n", "\n       ")
	fmt.Println("[WARN]", message)
}

func ErrorLog(message string) {
	if logger.level > Error {
		return
	}
	message = strings.ReplaceAll(message, "\n", "\n        ")
	fmt.Println("[ERROR]", message)
}
