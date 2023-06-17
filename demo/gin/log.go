package gin

import (
	"bookstore/demo/gin/config"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger *zap.Logger

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLoggerWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}

	return zapcore.AddSync(lumberJackLogger)
}

// InitLogger 初始化 Logger
func InitLogger(cfg *config.LogConfig) error {
	writerSyncer := getLoggerWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	encoder := getEncoder()

	l := new(zapcore.Level)
	err := l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return err
	}

	core := zapcore.NewCore(encoder, writerSyncer, l)
	logger = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可

	return nil
}

// GinLogger 基于 zap 的中间件
// 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(context *gin.Context) {
		start := time.Now()
		path := context.Request.URL.Path
		query := context.Request.URL.RawQuery
		context.Next()

		cost := time.Since(start)
		logger.Info(path,
			zap.Int("status", context.Writer.Status()),
			zap.String("method", context.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", context.ClientIP()),
			zap.String("user-agent", context.Request.UserAgent()),
			zap.String("errors", context.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("duration", cost),
		)
	}
}
