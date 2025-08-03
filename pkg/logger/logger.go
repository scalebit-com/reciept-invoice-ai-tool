package logger

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/scalebit-com/reciept-invoice-ai-tool/pkg/interfaces"
)

// ColorLogger implements the Logger interface with colored output
type ColorLogger struct {
	infoColor  *color.Color
	errorColor *color.Color
	warnColor  *color.Color
	debugColor *color.Color
	timeColor  *color.Color
}

// NewColorLogger creates a new ColorLogger instance
func NewColorLogger() interfaces.Logger {
	return &ColorLogger{
		infoColor:  color.New(color.FgGreen),
		errorColor: color.New(color.FgRed),
		warnColor:  color.New(color.FgYellow),
		debugColor: color.New(color.FgCyan),
		timeColor:  color.New(color.FgWhite),
	}
}

// formatMessage formats a log message with timestamp and level
func (l *ColorLogger) formatMessage(level string, levelColor *color.Color, msg string, args ...interface{}) {
	timestamp := l.timeColor.Sprintf("[%s]", time.Now().Format("2006-01-02 15:04:05"))
	levelStr := levelColor.Sprintf("[%s]", level)
	message := fmt.Sprintf(msg, args...)
	fmt.Printf("%s %s %s\n", timestamp, levelStr, message)
}

// Info logs an info message
func (l *ColorLogger) Info(msg string, args ...interface{}) {
	l.formatMessage("INFO", l.infoColor, msg, args...)
}

// Error logs an error message
func (l *ColorLogger) Error(msg string, args ...interface{}) {
	l.formatMessage("ERROR", l.errorColor, msg, args...)
}

// Warn logs a warning message
func (l *ColorLogger) Warn(msg string, args ...interface{}) {
	l.formatMessage("WARN", l.warnColor, msg, args...)
}

// Debug logs a debug message
func (l *ColorLogger) Debug(msg string, args ...interface{}) {
	l.formatMessage("DEBUG", l.debugColor, msg, args...)
}