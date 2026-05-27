package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"trbb/internal/config"
)

type Logger struct {
	z *zap.SugaredLogger
}

func New(cfg config.LogConfig) *Logger {
	if err := os.MkdirAll(cfg.Dir, 0755); err != nil {
		panic(fmt.Sprintf("cannot create log dir: %v", err))
	}

	level := zapcore.InfoLevel
	if cfg.Level == "debug" {
		level = zapcore.DebugLevel
	}

	// JSON encoder config
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "time"
	encCfg.MessageKey = "msg"
	encCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	enc := zapcore.NewJSONEncoder(encCfg)

	// Daily rotating file writer
	fileWriter := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Dir, fmt.Sprintf("trbb-%s.log", time.Now().Format("2006-01-02"))),
		MaxSize:    100, // MB
		MaxBackups: 30,
		MaxAge:     30, // days
		Compress:   true,
	}

	// Multi-writer: file + stdout
	multiWriter := io.MultiWriter(os.Stdout, fileWriter)
	writeSyncer := zapcore.AddSync(multiWriter)

	core := zapcore.NewCore(enc, writeSyncer, level)
	z := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &Logger{z: z.Sugar()}
}

func (l *Logger) Info(msg string, args ...any)  { l.z.Infow(msg, args...) }
func (l *Logger) Debug(msg string, args ...any) { l.z.Debugw(msg, args...) }
func (l *Logger) Warn(msg string, args ...any)  { l.z.Warnw(msg, args...) }
func (l *Logger) Error(msg string, args ...any) { l.z.Errorw(msg, args...) }
func (l *Logger) Fatal(msg string, args ...any) { l.z.Fatalw(msg, args...) }
