package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/config"
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	With(fields ...Field) Logger
	Sync() error
}

type Field = zap.Field

var (
	String  = zap.String
	Int     = zap.Int
	Int64   = zap.Int64
	Float64 = zap.Float64
	Bool    = zap.Bool
	Time    = zap.Time
	Duration = zap.Duration
	Any     = zap.Any
	Error   = zap.Error
)

type zapLogger struct {
	logger *zap.Logger
}

func NewLogger(cfg config.Config) (Logger, error) {
	loggingCfg := cfg.GetLoggingConfig()
	
	logLevel, err := parseLogLevel(loggingCfg.Level)
	if err != nil {
		return nil, err
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	var encoder zapcore.Encoder
	switch loggingCfg.Format {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	var writeSyncer zapcore.WriteSyncer
	switch loggingCfg.Output {
	case "stdout":
		writeSyncer = zapcore.AddSync(os.Stdout)
	case "stderr":
		writeSyncer = zapcore.AddSync(os.Stderr)
	default:
		file, err := os.OpenFile(loggingCfg.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			writeSyncer = zapcore.AddSync(os.Stdout)
		} else {
			writeSyncer = zapcore.AddSync(file)
		}
	}

	core := zapcore.NewCore(encoder, writeSyncer, logLevel)

	appCfg := cfg.GetAppConfig()
	var loggerOptions []zap.Option

	if appCfg.Environment == "development" {
		loggerOptions = append(loggerOptions, zap.AddCaller())
	}
	
	loggerOptions = append(loggerOptions, zap.AddStacktrace(zapcore.ErrorLevel))
	
	zapLoggerInstance := zap.New(core, loggerOptions...)

	return &zapLogger{logger: zapLoggerInstance}, nil
}

func parseLogLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, nil
	}
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{logger: l.logger.With(fields...)}
}

func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}

func WithRequestID(logger Logger, requestID string) Logger {
	return logger.With(String("request_id", requestID))
}

func WithOperation(logger Logger, operation string) Logger {
	return logger.With(String("operation", operation))
}

func WithDuration(logger Logger, duration float64) Logger {
	return logger.With(Float64("duration_ms", duration))
}

func WithError(logger Logger, err error) Logger {
	return logger.With(Error(err))
}

func WithComponent(logger Logger, component string) Logger {
	return logger.With(String("component", component))
}

func WithMethod(logger Logger, method string) Logger {
	return logger.With(String("method", method))
}

func WithPath(logger Logger, path string) Logger {
	return logger.With(String("path", path))
}

func WithStatusCode(logger Logger, statusCode int) Logger {
	return logger.With(Int("status_code", statusCode))
}

func WithIP(logger Logger, ip string) Logger {
	return logger.With(String("ip", ip))
}

func WithEnvironment(logger Logger, environment string) Logger {
	return logger.With(String("environment", environment))
}

func WithAppName(logger Logger, appName string) Logger {
	return logger.With(String("app_name", appName))
}

func WithVersion(logger Logger, version string) Logger {
	return logger.With(String("version", version))
} 