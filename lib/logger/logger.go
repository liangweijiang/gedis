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

// logLevel defines the log levels
type logLevel int

// Log level constants
const (
	DEBUG logLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// Color constants for terminal output
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Blue   = "\033[34m"
	Yellow = "\033[33m"
)

// Default values for logger configuration
const (
	defaultCallerDepth = 2
	maxBufferSize      = 1e5
	defaultDateFormat  = "2006-01-02"
)

// Level flags for log messages
var levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

// Settings holds the configuration for the logger
type Settings struct {
	Path      string // Path to the log files
	Name      string // Name of the log file
	WithColor bool   // Whether to use colored output
	WithJson  bool   // Whether to use JSON format for logs
}

// loggerEntry represents a single log entry
type loggerEntry struct {
	level   logLevel // Log level of the entry
	content string   // Content of the log message
}

// jsonContent represents the structure of a JSON log entry
type jsonContent struct {
	Level   string `json:"level"`   // Log level
	Caller  string `json:"caller"`  // Caller information
	Content string `json:"content"` // Log message content
}

// Logger is the main logger struct
type Logger struct {
	ctx       context.Context   // Context for the logger
	withJson  bool              // Whether to use JSON format
	logFile   *os.File          // File handle for the log file
	logger    *log.Logger       // Standard library logger
	entryChan chan *loggerEntry // Channel for log entries
	entryPool *sync.Pool        // Pool for log entries
}

// WriteContent writes a log message with the given level and content
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

// formatJsonLog formats the log message as JSON
func (l *Logger) formatJsonLog(level logLevel, callerInfo, content string) string {
	msg := jsonContent{
		Level:   levelFlags[level],
		Caller:  callerInfo,
		Content: content,
	}
	msgB, _ := json.Marshal(msg)
	return string(msgB)
}

// defaultLogger is the default logger instance
var defaultLogger = newStdLogger()

// Setup initializes the logger with the given settings
func Setup(settings Settings) {
	logger, err := newFileLogger(settings)
	if err != nil {
		panic(err)
	}
	defaultLogger = logger
}

// newStdLogger creates a new standard logger that outputs to stdout
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

// newFileLogger creates a new logger that outputs to a file
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

// Debug logs a debug message
func Debug(v ...interface{}) {
	content := fmt.Sprintln(v)
	defaultLogger.WriteContent(DEBUG, defaultCallerDepth, content)
}

// Debugf logs a formatted debug message
func Debugf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v)
	defaultLogger.WriteContent(DEBUG, defaultCallerDepth, content)
}

// Info logs an info message
func Info(v ...interface{}) {
	content := fmt.Sprintln(v)
	defaultLogger.WriteContent(INFO, defaultCallerDepth, content)
}

// Infof logs a formatted info message
func Infof(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v)
	defaultLogger.WriteContent(INFO, defaultCallerDepth, content)
}

// Warn logs a warning message
func Warn(v ...interface{}) {
	content := fmt.Sprintln(v)
	defaultLogger.WriteContent(WARN, defaultCallerDepth, content)
}

// Warnf logs a formatted warning message
func Warnf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v)
	defaultLogger.WriteContent(WARN, defaultCallerDepth, content)
}

// Error logs an error message
func Error(v ...interface{}) {
	content := fmt.Sprintln(v)
	defaultLogger.WriteContent(ERROR, defaultCallerDepth, content)
}

// Errorf logs a formatted error message
func Errorf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v)
	defaultLogger.WriteContent(ERROR, defaultCallerDepth, content)
}

// Fatal logs a fatal message and exits the program
func Fatal(v ...interface{}) {
	content := fmt.Sprintln(v)
	defaultLogger.WriteContent(FATAL, defaultCallerDepth, content)
	os.Exit(1)
}

// Fatalf logs a formatted fatal message and exits the program
func Fatalf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v)
	defaultLogger.WriteContent(FATAL, defaultCallerDepth, content)
	os.Exit(1)
}
