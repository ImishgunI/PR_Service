package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger interface {
	LogLeveled
	LogFormatted
	LogWithFields
}

type LogLeveled interface {
	Debug(args any)
	Info(args any)
	Warn(args any)
	Error(args any)
	Fatal(args any)
}

type LogFormatted interface {
	Errorf(format string, args any)
	Infof(format string, args any)
	Fatalf(format string, args any)
}

type LogWithFields interface {
	WithField(key string, value any) Logger
}

type LogrusLogger struct {
	entry *logrus.Entry
}

func New() Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetLevel(logrus.DebugLevel)
	return &LogrusLogger{
		entry: logrus.NewEntry(log),
	}
}

func (l *LogrusLogger) Debug(args any) {
	l.entry.Debug(args)
}

func (l *LogrusLogger) Info(args any) {
	l.entry.Info(args)
}

func (l *LogrusLogger) Warn(args any) {
	l.entry.Warn(args)
}

func (l *LogrusLogger) Error(args any) {
	l.entry.Error(args)
}

func (l *LogrusLogger) Fatal(args any) {
	l.entry.Fatal(args)
}

func (l *LogrusLogger) Errorf(format string, args any) {
	l.entry.Errorf(format, args)
}

func (l *LogrusLogger) Infof(format string, args any) {
	l.entry.Infof(format, args)
}

func (l *LogrusLogger) Fatalf(format string, args any) {
	l.entry.Fatalf(format, args)
}

func (l *LogrusLogger) WithField(key string, value any) Logger {
	return &LogrusLogger{
		entry: l.entry.WithField(key, value),
	}
}
