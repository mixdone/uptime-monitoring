package logger

import (
	"fmt"

	"github.com/mixdone/uptime-monitoring/internal/config"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(args ...any)
	Infof(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Debug(args ...any)
	Debugf(format string, args ...any)
	WithField(key string, value any) Logger
	WithFields(fields map[string]any) Logger
	WithError(err error) Logger
}

type logrusLogger struct {
	entry *logrus.Entry
}

func NewLoggerAdapter(base *logrus.Logger) Logger {
	return &logrusLogger{
		entry: logrus.NewEntry(base),
	}
}

func InitStructuredLogger(cfg *config.Config) (Logger, error) {
	base, err := newLogrusBase(cfg)
	if err != nil {
		return nil, err
	}
	return NewLoggerAdapter(base), nil
}

func newLogrusBase(cfg *config.Config) (*logrus.Logger, error) {
	log := logrus.New()

	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", cfg.Log.Level, err)
	}
	log.SetLevel(level)

	if cfg.Log.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return log, nil
}

func (l *logrusLogger) Info(args ...any)                  { l.entry.Info(args...) }
func (l *logrusLogger) Infof(format string, args ...any)  { l.entry.Infof(format, args...) }
func (l *logrusLogger) Error(args ...any)                 { l.entry.Error(args...) }
func (l *logrusLogger) Errorf(format string, args ...any) { l.entry.Errorf(format, args...) }
func (l *logrusLogger) Warn(args ...any)                  { l.entry.Warn(args...) }
func (l *logrusLogger) Warnf(format string, args ...any)  { l.entry.Warnf(format, args...) }
func (l *logrusLogger) Debug(args ...any)                 { l.entry.Debug(args...) }
func (l *logrusLogger) Debugf(format string, args ...any) { l.entry.Debugf(format, args...) }

func (l *logrusLogger) WithField(key string, value any) Logger {
	return &logrusLogger{
		entry: l.entry.WithField(key, value),
	}
}

func (l *logrusLogger) WithFields(fields map[string]any) Logger {
	logrusFields := logrus.Fields{}
	for k, v := range fields {
		logrusFields[k] = v
	}
	return &logrusLogger{
		entry: l.entry.WithFields(logrusFields),
	}
}

func (l *logrusLogger) WithError(err error) Logger {
	return &logrusLogger{
		entry: l.entry.WithError(err),
	}
}
