package testcase

import (
	"net/http"
	"testing"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugarLogger *zap.SugaredLogger

func getLogWriter() zapcore.WriteSyncer {
	// zap logger 加入 Lumberjack 纸质
	lumberJackLogger := &lumberjack.Logger{
		// 日志文件的位置
		Filename: "test.log",
		// 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxSize: 10,
		// 保留旧文件的最大个数
		MaxBackups: 5,
		// 保留旧文件的最大天数
		MaxAge:   30,
		Compress: false,
	}

	return zapcore.AddSync(lumberJackLogger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core)
	sugarLogger = logger.Sugar()
}

func simpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}

func TestZapLogger(t *testing.T) {
	InitLogger()
	defer sugarLogger.Sync()

	simpleHttpGet("www.sogo.com")
	simpleHttpGet("http://www.sogo.com")
}
