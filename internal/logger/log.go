package logger

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
)

//LogConfig is a configuration for logger
type LogConfig struct {
	Type      string `yaml:"type" json:"type" toml:"type"`
	Severity  string `yaml:"severity" json:"severity" toml:"severity"`
	DebugMode bool   `yaml:"debug mode" json:"debug_mode" toml:"debug_mode"`
}

type ctxlog struct{}

//WithLogger put logger to context
func WithLogger(ctx context.Context, logger *log.Logger) context.Context {
	return context.WithValue(ctx, ctxlog{}, logger)
}

//WithEntry put logger.Entry to context
func WithEntry(ctx context.Context, logger *log.Entry) context.Context {
	return context.WithValue(ctx, ctxlog{}, logger)
}

//Logger get logger from context
func Logger(ctx context.Context) *log.Logger {
	//	l, _ := ctx.Value("logger").(*log.Logger)
	l, ok := ctx.Value(ctxlog{}).(*log.Logger)
	if !ok {
		//return DefaultLogger
		l = initLogger(LogConfig{Type: "stdout", Severity: "LOG_INFO"})
	}
	return l
}

//Entry get logger from context
func Entry(ctx context.Context) *log.Entry {
	//	l, _ := ctx.Value("logger").(*log.Entry)
	l, ok := ctx.Value(ctxlog{}).(*log.Entry)
	if !ok {
		//return DefaultLogger
		l = initLogger(LogConfig{Type: "stdout", Severity: "LOG_INFO"}).WithField("", "")
	}
	return l
}

//CreateLogger from config
func CreateLogger(config LogConfig) *log.Logger {
	logger := initLogger(config)
	return logger
}

func initLogger(config LogConfig) *log.Logger {
	logger := log.New()
	switch config.Type {
	case "stdout":
		logger.Out = os.Stdout
		return logger
	case "stderr":
		logger.Out = os.Stderr
		return logger
	default:
	}

	if config.DebugMode {
		logger.Out = os.Stdout
	}
	logger.Formatter = &log.TextFormatter{}
	logger.Level = logLevel[config.Severity]
	return logger
}

var logLevel = map[string]log.Level{
	"LOG_EMERG":   log.PanicLevel,
	"LOG_ALERT":   log.PanicLevel,
	"LOG_CRIT":    log.FatalLevel,
	"LOG_ERR":     log.ErrorLevel,
	"LOG_WARNING": log.WarnLevel,
	"LOG_NOTICE":  log.InfoLevel,
	"LOG_INFO":    log.InfoLevel,
	"LOG_DEBUG":   log.DebugLevel,
}
