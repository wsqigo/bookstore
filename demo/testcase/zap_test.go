package testcase

import (
	"bookstore/demo/gin/config"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/natefinch/lumberjack"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

// https://www.cnblogs.com/jiujuan/p/17304844.html

var url = "https://www.cnblogs.com/jiujuan/p/17304844.html"

func TestSugaredLogger(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // zap 底层有缓冲。在任何情况下执行 defer logger.Sync() 是一个很好的习惯

	sugar := logger.Sugar()
	sugar.Info("failed to fetch URL",
		// 字段是松散类型，不是强类型
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
}

// 当性能和类型安全很重要时，请使用 Logger。它比 SugaredLogger 更快，分配的资源更少，但它只支持结构化日志和强类型字段。
func TestLogger(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("failed to fetch URL",
		// 字段是强类型，不是松散类型
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}

func TestNewExample(t *testing.T) {
	logger := zap.NewExample()
	logger.Debug("this is debug message")
	logger.Info("this is info message")
	logger.Info("this is info message with fields",
		zap.Int("age", 37),
		zap.String("agent", "man"),
	)
	logger.Warn("this is warn message")
	logger.Error("this is error message")

	// Output:
	// {"level":"debug","msg":"this is debug message"}
	// {"level":"info","msg":"this is info message"}
	// {"level":"info","msg":"this is info message with fields","age":37,"agender":"man"}
	// {"level":"warn","msg":"this is warn message"}
	// {"level":"error","msg":"this is error message"}
}

func TestNewDevelopment(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	logger.Info("failed to fetch URL",
		// 强类型字段
		zap.String("url", "http://example.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	logger.With( // 强类型字段
		zap.String("url", "http://development.com"),
		zap.Int("attempt", 4),
		zap.Duration("duration", time.Second*5),
	).Info("[With] failed to fetch url")
	// Output:
	// 2023-07-05T01:01:01.261+0800	INFO	testcase/zap_test.go:64	failed to fetch URL	{"url": "http://example.com", "attempt": 3, "backoff": "1s"}
	// 2023-07-05T01:01:01.277+0800	INFO	testcase/zap_test.go:75	[With] failed to fetch url	{"url": "http://development.com", "attempt": 4, "duration": "5s"}
}

func TestNewProduction(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	url := "http://zap.uber.io"
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		"url", url,
		"attempt", 3,
		"time", time.Second,
	)

	sugar.Infof("Failed to fetch URL: %s", url)

	// Output:
	// {"level":"info","ts":1679472893.2944522,"caller":"zapdemos/newproduction1.go:16","msg":"failed to fetch URL","url":"http://zap.uber.io","attempt":3,"time":1}
	// {"level":"info","ts":1679472893.294975,"caller":"zapdemos/newproduction1.go:22","msg":"Failed to fetch URL: http://zap.uber.io"}
}

// 使用配置
func TestConfig(t *testing.T) {
	logger, _ := zap.NewProduction(zap.Fields(
		zap.String("log_name", "testlog"),
		zap.String("log_author", "prometheus"),
	))
	defer logger.Sync()

	logger.Info("test fields output")

	logger.Warn("warn info")

	// Output
	// {"level":"info","ts":1688574134.260133,"caller":"testcase/zap_test.go:108","msg":"test fields output","log_name":"testlog","log_author":"prometheus"}
	// {"level":"warn","ts":1688574134.260133,"caller":"testcase/zap_test.go:110","msg":"warn info","log_name":"testlog","log_author":"prometheus"}
}

func TestHook(t *testing.T) {
	logger := zap.NewExample(zap.Hooks(func(entry zapcore.Entry) error {
		fmt.Println("[zap.Hooks]test Hooks")
		return nil
	}))
	defer logger.Sync()

	logger.Info("test output")

	logger.Warn("warn info")
}

func TestWriteFile(t *testing.T) {
	// 设置一些配置参数
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(cfg)

	defaultLogLevel := zapcore.DebugLevel // 设置 loglevel

	logFile, _ := os.OpenFile("./log-test-zap.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	// or os.Create()
	writer := zapcore.AddSync(logFile)

	logger := zap.New(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	defer logger.Sync()

	url := "http://www.test.com"
	logger.Info("write log to file",
		zap.String("url", url),
		zap.Int("attemp", 3),
	)
}

func TestBasicConfig(t *testing.T) {
	// 表示 zap.Config 的 json 原始编码
	// outputPath: 设置日志输出路径，日志内容输出到标准输出和文件 logs.log
	// errorOutputPaths：设置错误日志输出路径
	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "./logs.log"],
	  "errorOutputPaths": ["stderr"],
	  "initialFields": {"foo": "bar"},
	  "encoderConfig": {
		"messageKey": "message-customer",
		"levelKey": "level",
		"levelEncoder": "lowercase"
	  }
	}`)

	// 把 json 格式数据解析到 zap.Config struct
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	// cfg.Build() 为配置对象创建一个 Logger
	// zap.Must() 封装了 Logger, Must() 函数如果返回值不是 nil，就会报 panic
	// 也就是检查 Build 是否错误
	logger := zap.Must(cfg.Build())
	defer logger.Sync()

	logger.Info("logger construction succeeded")
}

func TestLevelFile(t *testing.T) {
	// 设置配置
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(cfg)

	logFile, _ := os.OpenFile("./log-debug-zap.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666) // 日志记录debug信息
	errFile, _ := os.OpenFile("./log-error-zap.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666) // 日志记录error信息

	teecore := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zap.DebugLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(errFile), zap.ErrorLevel),
	)

	logger := zap.New(teecore, zap.AddCaller())
	defer logger.Sync()

	url := "http://www.diff-log-level.com"
	logger.Info("write log to file",
		zap.String("url", url),
		zap.Int("time", 3),
	)

	logger.With(
		zap.String("url", url),
		zap.String("name", "jimmmyr"),
	).Error("test error ")
}

func TestLumberJackLogger(t *testing.T) {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./log-rotate-test.json",
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // disabled by default
	}
	defer lumberJackLogger.Close()

	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder // 设置时间格式
	fileEncoder := zapcore.NewJSONEncoder(cfg)

	// 编码设置 输出到文件 日志等级
	core := zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberJackLogger), zap.InfoLevel)

	logger := zap.New(core)
	defer logger.Sync()

	// 测试分割日志
	for i := 0; i < 8000; i++ {
		logger.With(
			zap.String("url", fmt.Sprintf("www.test%d.com", i)),
			zap.String("name", "jimmmyr"),
			zap.Int("age", 23),
			zap.String("agradege", "no111-000222"),
		).Info("test info")
	}

}

var lg *zap.Logger

// InitLogger 初始化 Logger
func InitLogger(cfg *config.LogConfig) error {
	writeSyncer := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	encoder := getEncoder()

	var l = new(zapcore.Level)
	err := l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return err
	}

	core := zapcore.NewCore(encoder, writeSyncer, l)

	lg = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(lg) // 替换 zap 包中全局的 logger 实例，后续在其他包中只需使用 zap.L() 调用即可
	return nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time" // 记录日志时间的键名，默认为 level
	// 日志编码的一些设置项
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}

	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger 接收 gin 框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		ctx.Next()

		cost := time.Since(start)
		zap.L().Info(path,
			zap.Int("status", ctx.Writer.Status()),
			zap.String("method", ctx.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", ctx.ClientIP()),
			zap.String("user-agent", ctx.Request.UserAgent()),
			zap.String("errors", ctx.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recovery
