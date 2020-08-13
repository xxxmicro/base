package log

import(
	"os/exec"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(env string) {
	if env == "dev" {
		logConfig := zap.NewDevelopmentConfig()
		logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		log, _ = logConfig.Build()
		log = log.WithOptions(
			zap.Hooks(func(entry zapcore.Entry) error {
				cmd := exec.Command("/usr/bin/say", entry.Message)
				// 执行命令，返回命令是否执行成功
				_ = cmd.Run()
				return nil
			}),
		)
	} else if env != "production" {
		logConfig := zap.NewDevelopmentConfig()
		logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		log, _ = logConfig.Build()
		log = log.WithOptions(
			zap.Hooks(func(entry zapcore.Entry) error {
				return nil
			}),
		)
	} else {
		log, _ = zap.NewProduction()
	}
}

func Debug(f string, args ...interface{}) {
	log.Debug(f)
}

func Error(e error) {
	
}

func Info(f string) {
	log.Info(f)
}

func Fatal(f string, args ...interface{}) {
	log.Panic(f)
}

func Panic(f string) {
	log.Panic(f)
}
