package app

import (
	"context"
	"fmt"
	"github.com/natefinch/lumberjack"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type SugaredLogger struct {
	*zap.SugaredLogger
}

var Log SugaredLogger

func (s *SugaredLogger) WithContext(ctx *context.Context) *zap.SugaredLogger {
	defer func() {
		*ctx = context.WithValue(*ctx, "time", time.Now())
	}()
	s.initContext(ctx)
	fields := make([]interface{}, 0)
	fields = append(fields, zap.Int64("timeline", time.Now().UnixNano()))
	fields = append(fields, zap.Duration("duration", time.Since(((*ctx).Value("time")).(time.Time))))
	for _, s2 := range []string{"request", "response", "context", "category", "ip", "type", "sub_type", "trace_id"} {
		if (*ctx).Value(s2) != nil {
			fields = append(fields, zap.Any(s2, (*ctx).Value(s2)))
		}
	}
	return s.With(fields...)
}

func (s SugaredLogger) initContext(ctx *context.Context) {
	if (*ctx).Value("time") == nil {
		*ctx = context.WithValue(*ctx, "time", time.Now())
	}
	if (*ctx).Value("trace_id") == nil {
		*ctx = context.WithValue(*ctx, "trace_id", uuid.NewV4())
	}
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		Log.Info(fmt.Sprintf("Could not find local user folder. Error: %v\n", err))
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	//encoder := zapcore.NewJSONEncoder(encoderConfig) // 使用JSON的编码器
	encoder := zapcore.NewConsoleEncoder(encoderConfig) // 使用制表符分隔的编码器
	lumberJackLogger := &lumberjack.Logger{
		Filename:   userHomeDir + string(os.PathSeparator) + "mypass.log",
		MaxSize:    500,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
	}
	writerSyncer := zapcore.AddSync(lumberJackLogger)
	//consoleSyncer := zapcore.AddSync(os.Stdout)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel),
		//zapcore.NewCore(encoder, consoleSyncer, zapcore.DebugLevel),
	)
	log := zap.New(core, zap.AddCaller())
	Log.SugaredLogger = log.Sugar()
}
