package elog

import (
	"context"
	"log/slog"
)

type Event struct {
	msg   string
	slog  *slog.Logger
	level slog.Level
}

// Log the event, should be `defer`ed.
func (e *Event) Log() {
	e.slog.Log(context.Background(), e.level, e.msg)
}

// Add a attribute.
func (e *Event) With(args ...any) {
	*e.slog = *e.slog.With(args...)
}

// Add a `error` field and change the log level to "error".
func (e *Event) WithError(err error) {
	*e.slog = *e.slog.With("error", err)
	e.level = slog.LevelError
}

// Add s `warn` field and change the log level to "warn" (if an error was not already logged).
func (e *Event) WithWarn(err error) {
	*e.slog = *e.slog.With("warn", err)

	// Only set the level to warn
	if e.level < slog.LevelWarn {
		e.level = slog.LevelWarn
	}
}

// Set the log level manually.
func (e *Event) SetLevel(l slog.Level) {
	e.level = l
}

type Logger struct {
	slog *slog.Logger
}

func (l *Logger) NewEvent(msg string, args ...any) *Event {
	return &Event{
		// Copy the logger.
		slog: l.slog.With(args...),

		// Set the message.
		msg: msg,

		// Set the log level to info.
		level: slog.LevelInfo,
	}
}

func New(sl *slog.Logger) *Logger {
	return &Logger{
		slog: sl,
	}
}
