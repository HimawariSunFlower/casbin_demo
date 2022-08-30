package util

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"time"
)

var Logger *zap.SugaredLogger

func InitLogger() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	Logger = zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(config),
		zapcore.NewMultiWriteSyncer(rotateLog(), zapcore.AddSync(os.Stdout)),
		zap.DebugLevel), zap.AddCaller()).Sugar()
}

func rotateLog() zapcore.WriteSyncer {
	fileWriter, _ := rotatelogs.New(
		path.Join("./log/", "%Y-%m-%d.log"),
		rotatelogs.WithClock(rotatelogs.Local),
		rotatelogs.WithMaxAge(time.Duration(24)*24*time.Hour), // 日志留存时间
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	return zapcore.AddSync(fileWriter)
}
