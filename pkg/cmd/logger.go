package cmd

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

const (
	logSuffix  = ".log"
	timeFormat = "2006-01-02 15:04:05.000"
)

// 配置日志行为
func initLogger(path string, level int, color, caller bool) {
	log.SetLevel(log.AllLevels[level%len(log.AllLevels)])
	log.SetReportCaller(caller)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:   color,
		TimestampFormat: timeFormat,
	})
	log.AddHook(fileLoggerHook(path))
}

// 文件日志
func fileLoggerHook(logPath string) log.Hook {
	infoWriter, _ := rotatelogs.New(logPath+"%Y%m%d.info"+logSuffix,
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour))
	warnWriter, _ := rotatelogs.New(logPath+"%Y%m%d.warn"+logSuffix,
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour))
	errWriter, _ := rotatelogs.New(logPath+"%Y%m%d.err"+logSuffix,
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour))

	fileFormatter := &log.TextFormatter{
		DisableColors:   true,
		TimestampFormat: timeFormat,
	}

	return lfshook.NewHook(lfshook.WriterMap{
		log.InfoLevel:  infoWriter,
		log.WarnLevel:  warnWriter,
		log.ErrorLevel: errWriter,
	}, fileFormatter)
}
