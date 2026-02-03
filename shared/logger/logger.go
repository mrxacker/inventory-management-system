package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	Panic(msg string, fields ...interface{})
	With(fields ...interface{}) Logger
	Sync() error
}

type zapLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

// NewLogger creates a new logger instance
func NewLogger() Logger {
	return NewLoggerWithConfig(LogConfig{
		Level:       "info",
		Environment: "development",
		OutputPaths: []string{"stdout"},
		ErrorPaths:  []string{"stderr"},
	})
}

// LogConfig holds logger configuration
type LogConfig struct {
	Level       string
	Environment string
	OutputPaths []string
	ErrorPaths  []string
}

// NewLoggerWithConfig creates a logger with custom configuration
func NewLoggerWithConfig(config LogConfig) Logger {
	var zapConfig zap.Config

	if config.Environment == "production" {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Set log level
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(config.Level)); err == nil {
		zapConfig.Level = zap.NewAtomicLevelAt(level)
	}

	// Set output paths
	if len(config.OutputPaths) > 0 {
		zapConfig.OutputPaths = config.OutputPaths
	}

	if len(config.ErrorPaths) > 0 {
		zapConfig.ErrorOutputPaths = config.ErrorPaths
	}

	// Customize encoder config
	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	zapConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	logger, err := zapConfig.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		panic(err)
	}

	return &zapLogger{
		logger: logger,
		sugar:  logger.Sugar(),
	}
}

// NewFileLogger creates a logger that writes to a file
func NewFileLogger(filename string, level string, environment string) Logger {
	return NewLoggerWithConfig(LogConfig{
		Level:       level,
		Environment: environment,
		OutputPaths: []string{filename, "stdout"},
		ErrorPaths:  []string{filename, "stderr"},
	})
}

func (l *zapLogger) Debug(msg string, fields ...interface{}) {
	l.sugar.Debugw(msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...interface{}) {
	l.sugar.Infow(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...interface{}) {
	l.sugar.Warnw(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...interface{}) {
	l.sugar.Errorw(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...interface{}) {
	l.sugar.Fatalw(msg, fields...)
}

func (l *zapLogger) Panic(msg string, fields ...interface{}) {
	l.sugar.Panicw(msg, fields...)
}

func (l *zapLogger) With(fields ...interface{}) Logger {
	return &zapLogger{
		logger: l.logger,
		sugar:  l.sugar.With(fields...),
	}
}

func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}

// GetZapLogger returns the underlying zap.Logger for advanced usage
func (l *zapLogger) GetZapLogger() *zap.Logger {
	return l.logger
}

// GetSugaredLogger returns the underlying zap.SugaredLogger
func (l *zapLogger) GetSugaredLogger() *zap.SugaredLogger {
	return l.sugar
}
