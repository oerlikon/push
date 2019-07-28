package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	logDir    = "./logs/"
	logPrefix = "push.log."
)

var logger struct {
	mutex     sync.Mutex
	writer    io.Writer
	startTime time.Time
	echoMutex sync.Mutex
}

func init() {
	logger.startTime = time.Now()
}

func initlog() {
	if logger.startTime.IsZero() {
		logger.startTime = time.Now()
	}
	if err := os.MkdirAll(logDir, 0755); err != nil {
		os.RemoveAll(logDir)
		logger.writer = ioutil.Discard
		return
	}
	n := filepath.Join(logDir, logPrefix+logger.startTime.Format("20060102.150405.00000"))
	f, err := os.Create(n)
	if err != nil {
		logger.writer = ioutil.Discard
		return
	}
	logger.writer = f
}

func Echo(format string, a ...interface{}) {
	logger.echoMutex.Lock()
	fmt.Fprintf(os.Stderr, format, a...)
	logger.echoMutex.Unlock()
}

func Echoln(format string, a ...interface{}) {
	logger.echoMutex.Lock()
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	logger.echoMutex.Unlock()
}

func Logf(format string, a ...interface{}) {
	s := fmt.Sprintf(strings.TrimSpace(format), a...)
	if len(s) == 0 {
		return
	}
	logger.mutex.Lock()
	if logger.writer == nil {
		initlog()
	}
	fmt.Fprintln(logger.writer, time.Time{}.Add(time.Since(logger.startTime)).Format("15:04:05.000"), s)
	logger.mutex.Unlock()
}

func LogEcho(format string, a ...interface{}) {
	Logf(format, a...)
	Echo(format, a...)
}

func LogEcholn(format string, a ...interface{}) {
	Logf(format, a...)
	Echoln(format, a...)
}
