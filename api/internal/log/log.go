package log

import (
	"os"
	"strings"
	"sync"

	"github.com/Kudryavkaz/sztuea-api/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	once   sync.Once
	logger *zap.Logger

	zapLevelMapping = map[string]zapcore.Level{
		"DEBUG":  zapcore.DebugLevel,
		"INFO":   zapcore.InfoLevel,
		"WARN":   zapcore.WarnLevel,
		"ERROR":  zapcore.ErrorLevel,
		"DPANIC": zapcore.DPanicLevel,
		"PANIC":  zapcore.PanicLevel,
		"FATAL":  zapcore.FatalLevel,
	}
	DebugLevel  = zap.DebugLevel
	InfoLevel   = zap.InfoLevel
	WarnLevel   = zap.WarnLevel
	ErrorLevel  = zap.ErrorLevel
	DPanicLevel = zap.DPanicLevel
	PanicLevel  = zap.PanicLevel
	FatalLevel  = zap.FatalLevel
)

func Logger() *zap.Logger {
	once.Do(func() {
		logLevel, ok := zapLevelMapping[strings.ToUpper(config.Config.GetString("log.level"))]
		if !ok {
			logLevel = zapcore.InfoLevel
		}

		var cfg zap.Config
		if logLevel == zap.DebugLevel {
			cfg = zap.NewDevelopmentConfig()
		} else {
			cfg = zap.NewProductionConfig()
		}

		logBasePath := "/data/sztuea-api/logs/sztuea"
		// General log directory Rotation
		logFileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logBasePath + ".log",
			MaxSize:    20, // Log volume / MB
			MaxBackups: 3,
			MaxAge:     7, // Log retention days
			Compress:   false,
		})
		// Error log directory Rotation
		errLogWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logBasePath + "_error.log",
			MaxSize:    20, // Log volume / MB
			MaxBackups: 3,
			MaxAge:     7, // Log retention days
			Compress:   false,
		})
		panicLogWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logBasePath + "_panic.log",
			MaxSize:    20, // Log volume / MB
			MaxBackups: 3,
			MaxAge:     7, // Log retention days
			Compress:   false,
		})

		// Set log level
		cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		cfg.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
		cfg.EncoderConfig.ConsoleSeparator = "  "
		core := zapcore.NewTee(
			zapcore.NewCore(zapcore.NewConsoleEncoder(cfg.EncoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), logFileWriter), logLevel),
			// Error log / Panic log output to err.log
			zapcore.NewCore(zapcore.NewConsoleEncoder(cfg.EncoderConfig), zapcore.NewMultiWriteSyncer(errLogWriter), ErrorLevel),
			// PanicLog output to panic.log
			zapcore.NewCore(zapcore.NewConsoleEncoder(cfg.EncoderConfig), zapcore.NewMultiWriteSyncer(panicLogWriter), DPanicLevel),
		)

		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.PanicLevel))
	})
	return logger
}
