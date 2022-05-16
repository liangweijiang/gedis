package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Config struct {
	Path       string `yaml:"path"`
	Name       string `yaml:"name"`
	Ext        string `yaml:"ext"`
	TimeFormat string `yaml:"time-format"`
}

var (
	logger             *log.Logger
	logFile            *os.File
	defaultPrefix      = ""
	defaultCallerDepth = 2
	lock               sync.Mutex
	logPrefix          = ""
	levelFlags         = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

type logLevel int

const (
	DEBUG logLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

const flags = log.LstdFlags

func init() {
	logger = log.New(os.Stdout, defaultPrefix, flags)
}

func InitLog(conf *Config) {
	var err error
	dir := conf.Path
	fileName := fmt.Sprintf("%s-%s.%s", conf.Name, time.Now().Format(conf.TimeFormat), conf.Ext)
	logFile, err = mustOpen(fileName, dir)
	if err != nil {
		log.Fatalf("init log err:%+v", err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	logger = log.New(mw, defaultPrefix, flags)
}

func setPrefix(level logLevel) {
	// 函数的返回值为调用栈标识符、带路径的完整文件名、该调用在文件中的行号。如果无法获得信息，ok会被设为false。
	_, file, line, ok := runtime.Caller(defaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}
	logger.SetPrefix(logPrefix)
}

// Debug prints debug log
func Debug(v ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	setPrefix(DEBUG)
	logger.Println(v...)
}

// Info prints normal log
func Info(v ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	setPrefix(INFO)
	logger.Println(v...)
}

// Warn prints warning log
func Warn(v ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	setPrefix(WARNING)
	logger.Println(v...)
}

// Error prints error log
func Error(v ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	setPrefix(ERROR)
	logger.Println(v...)
}

// Fatal prints error log then stop the program
func Fatal(v ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	setPrefix(FATAL)
	logger.Fatalln(v...)
}

// Debugf prints debug log
func Debugf(format string, v ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	setPrefix(DEBUG)
	logger.Printf(format, v...)
}

// Infof prints normal log
func Infof(format string, v ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	setPrefix(INFO)
	logger.Printf(format, v...)
}

// Warnf prints warning log
func Warnf(format string, v ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	setPrefix(WARNING)
	logger.Printf(format, v...)
}

// Errorf prints error log
func Errorf(format string, v ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	setPrefix(ERROR)
	logger.Printf(format, v...)
}

// Fatalf prints error log then stop the program
func Fatalf(format string, v ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	setPrefix(FATAL)
	logger.Fatalf(format, v...)
}
