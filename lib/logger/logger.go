package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/liangweijiang/gedis/lib/utils"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type logLevel int

const (
	DEBUG logLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Blue   = "\033[34m"
	Yellow = "\033[33m"
)

const (
	defaultCallerDepth = 2
	maxBufferSize      = 1e5
	defaultDateFormat  = "2006-01-02"
)

var levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

type Settings struct {
	Path      string
	Name      string
	WithColor bool
	WithJson  bool
}

type loggerEntry struct {
	level   logLevel
	content string
}

type jsonContent struct {
	Level   string `json:"level"`
	Caller  string `json:"caller,omitempty"`
	Content string `json:"content"`
}

type Logger struct {
	ctx       context.Context
	withJson  bool
	logFile   *os.File
	logger    *log.Logger
	entryChan chan *loggerEntry
	entryPool *sync.Pool
}

func (l *Logger) WriteContent(level logLevel, callerDepth int, content string) {
	var formattedContent string
	_, file, line, ok := runtime.Caller(callerDepth)

	callerInfo := ""
	if ok {
		callerInfo = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	if l.withJson {
		formattedContent = l.formatJsonLog(level, callerInfo, content)
	} else {
		if ok {
			formattedContent = fmt.Sprintf("[%s][%s] %s", levelFlags[level], callerInfo, content)
		} else {
			formattedContent = fmt.Sprintf("[%s] %s", levelFlags[level], content)
		}
	}

	entry := l.entryPool.Get().(*loggerEntry)
	entry.level = level
	entry.content = formattedContent
	l.entryChan <- entry
}

func (l *Logger) formatJsonLog(level logLevel, callerInfo, content string) string {
	msg := jsonContent{
		Level:   levelFlags[level],
		Caller:  callerInfo,
		Content: content,
	}
	msgB, _ := json.Marshal(msg)
	return string(msgB)
}

var defaultLogger = newStdLogger()

func Setup(settings Settings) {
	logger, err := newFileLogger(settings)
	if err != nil {
		panic(err)
	}
	defaultLogger = logger
}

func newStdLogger() *Logger {
	stdLogger := &Logger{
		logFile:   nil,
		logger:    log.New(os.Stdout, "", log.LstdFlags),
		entryChan: make(chan *loggerEntry, maxBufferSize),
		entryPool: &sync.Pool{
			New: func() interface{} {
				return &loggerEntry{}
			},
		},
	}
	go func() {
		for entry := range stdLogger.entryChan {
			content := entry.content
			switch entry.level {
			case DEBUG:
				content = Blue + content + Reset
			case INFO:
				content = Green + content + Reset
			case WARN:
				content = Yellow + content + Reset
			case ERROR, FATAL:
				content = Red + content + Reset
			}
			_ = stdLogger.logger.Output(0, content)
			stdLogger.entryPool.Put(entry)
		}
	}()
	return stdLogger
}

func newFileLogger(settings Settings) (*Logger, error) {
	fileName := fmt.Sprintf("%s-%s.log", settings.Name, time.Now().Format(defaultDateFormat))
	logFile, err := utils.MustOpen(settings.Path, fileName)
	if err != nil {
		return nil, err
	}
	logger := &Logger{
		logFile:   logFile,
		logger:    log.New(io.MultiWriter(os.Stdout, logFile), "", log.LstdFlags),
		entryChan: make(chan *loggerEntry, maxBufferSize),
		entryPool: &sync.Pool{
			New: func() interface{} {
				return &loggerEntry{}
			},
		},
		withJson: settings.WithJson,
	}
	go func() {
		for entry := range logger.entryChan {
			logFilename := fmt.Sprintf("%s-%s.log", settings.Name, time.Now().Format(defaultDateFormat))
			if path.Join(settings.Path, logFilename) != logger.logFile.Name() {
				newFile, err := utils.MustOpen(settings.Path, logFilename)
				if err != nil {
					panic("open log " + logFilename + " failed: " + err.Error())
				}
				logger.logFile = newFile
				logger.logger = log.New(io.MultiWriter(os.Stdout, logFile), "", log.LstdFlags)
			}

			content := entry.content
			if settings.WithColor {
				switch entry.level {
				case DEBUG:
					content = Blue + content + Reset
				case INFO:
					content = Green + content + Reset
				case WARN:
					content = Yellow + content + Reset
				case ERROR, FATAL:
					content = Red + content + Reset
				}
			}

			_ = logger.logger.Output(0, content)
		}
	}()
	return logger, nil
}

func Debug(v ...interface{}) {
	content := fmt.Sprintln(v)
	defaultLogger.WriteContent(DEBUG, defaultCallerDepth, content)
}

func Debugf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v)
	defaultLogger.WriteContent(DEBUG, defaultCallerDepth, content)
}

func Info(v ...interface{}) {
	content := fmt.Sprintln(v)
	defaultLogger.WriteContent(INFO, defaultCallerDepth, content)
}

func Infof(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v)
	defaultLogger.WriteContent(INFO, defaultCallerDepth, content)
}

func Warn(v ...interface{}) {
	content := fmt.Sprintln(v)
	defaultLogger.WriteContent(WARN, defaultCallerDepth, content)
}

func Warnf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v)
	defaultLogger.WriteContent(WARN, defaultCallerDepth, content)
}

func Error(v ...interface{}) {
	content := fmt.Sprintln(v)
	defaultLogger.WriteContent(ERROR, defaultCallerDepth, content)
}

func Errorf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v)
	defaultLogger.WriteContent(ERROR, defaultCallerDepth, content)
}

func Fatal(v ...interface{}) {
	content := fmt.Sprintln(v)
	defaultLogger.WriteContent(FATAL, defaultCallerDepth, content)
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v)
	defaultLogger.WriteContent(FATAL, defaultCallerDepth, content)
	os.Exit(1)
}
