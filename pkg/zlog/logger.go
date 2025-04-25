package zlog

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go_chat/internal/config"
	"os"
	"path"
	"runtime"
)

var logger *zap.Logger
var logPath string

func init() {
	//new一个
	encoderConfig := zap.NewProductionEncoderConfig()
	//设置日志时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//格式化为json
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	//加载toml中的log级别信息
	conf := config.GetConfig()
	logPath = conf.LogPath
	//指定日志文件以及写入操作
	file, _ := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 644)
	fileWriteSyncer := zapcore.AddSync(file)
	//日志会同时输出到控制台（os.Stdout）和文件  两个 NewCore 合并为一个 Tee
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel),
	)
	logger = zap.New(core)
}

// 日志分割
func getFileLogWriter() (writerSyncer zapcore.WriteSyncer) {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100,
		MaxBackups: 60,
		MaxAge:     1,
		Compress:   false,
	}
	return zapcore.AddSync(lumberjackLogger)
}

// 获取函数名、文件名、行号记录在日志中
func getCallerInfoLog() (callerFields []zap.Field) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return
	}
	//获取函数名
	funcName := path.Base(runtime.FuncForPC(pc).Name())

	callerFields = append(callerFields, zap.String("func", funcName), zap.String("file", file), zap.Int("line", line))
	return
}

func Info(message string, fields ...zap.Field) {
	callerInfoLog := getCallerInfoLog()
	fields = append(fields, callerInfoLog...)
	logger.Info(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	callerInfoLog := getCallerInfoLog()
	fields = append(fields, callerInfoLog...)
	logger.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	callerInfoLog := getCallerInfoLog()
	fields = append(fields, callerInfoLog...)
	logger.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	callerInfoLog := getCallerInfoLog()
	fields = append(fields, callerInfoLog...)
	logger.Fatal(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	callerInfoLog := getCallerInfoLog()
	fields = append(fields, callerInfoLog...)
	logger.Debug(message, fields...)
}
