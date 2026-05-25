package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

const (
	MaxLogSize = 10 * 1024 * 1024 // 10MB
)

var logFile *os.File
var currentLogPath string

func InitLogger(logPath string) {
	var err error
	currentLogPath = logPath
	// TODO: an absolute monstrosity fix this later
	path := filepath.Join(append([]string{"/"}, strings.Split(logPath, "/")...)...)
	if filepath.Dir(path) != "." {
		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			fmt.Printf("Critical: Could not create directory: %v\n", err)
			os.Exit(1)
		}
	}
	logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Critical: Could not open log file: %v\n", err)
		os.Exit(1)
	}
}

func rotateLog() {
	if logFile != nil {
		logFile.Close()
	}
	timestamp := time.Now().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.%s", currentLogPath, timestamp)
	os.Rename(currentLogPath, backupPath)
	InitLogger(currentLogPath)
}

func Log(level, message string) {
	if logFile == nil {
		InitLogger("app.log")
	}

	stat, err := logFile.Stat()
	if err == nil && stat.Size() > MaxLogSize {
		rotateLog()
	}

	_, file, line, ok := runtime.Caller(1)
	location := "unknown:0"
	if ok {
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				file = file[i+1:]
				break
			}
		}
		location = fmt.Sprintf("%s:%d", file, line)
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(logFile, "[%s] [%s] [%s] %s\n", level, location, timestamp, message)
}
