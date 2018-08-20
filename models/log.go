package models

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	ErrorLevel int = iota
	WarnLevel
	InfoLevel
	DebugLevel
)

type LogFile struct {
	level    int
	logTime  int64
	fileName string
	fileFd   io.Writer
}

var logFile LogFile

func Config(logFolder string, level int) {
	fileFd, err := os.Create(logFolder)
	if err != nil {
		log.Fatalln("open log file error")
	}

	writers := []io.Writer{
		fileFd,
		os.Stdout,
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	logFile.fileName = logFolder
	logFile.level = level
	logFile.fileFd = fileAndStdoutWriter
	log.SetOutput(logFile.fileFd)
	log.SetFlags(log.Llongfile | log.LstdFlags)
}
func SetLevel(level int) {
	logFile.level = level
}

func Debugf(format string, args ...interface{}) {
	if logFile.level >= DebugLevel {
		log.SetPrefix("[Debug]")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Infof(format string, args ...interface{}) {
	if logFile.level >= InfoLevel {
		log.SetPrefix("[Info]")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Warnf(format string, args ...interface{}) {
	if logFile.level >= WarnLevel {
		log.SetPrefix("[Warning]")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Errorf(format string, args ...interface{}) {
	if logFile.level >= ErrorLevel {
		log.SetPrefix("[Error]")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Fatalf(format string, args ...interface{}) {
	log.SetPrefix("[Error]")
	log.Output(2, fmt.Sprintf(format, args...))
	os.Exit(1)
}
