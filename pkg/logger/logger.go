package logger

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

type Logger struct {
	log *logrus.Logger
}

func NewLogger(logLevel string) *Logger {
	log := logrus.New()
	customFormatter := new(logrus.JSONFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"

	log.SetFormatter(customFormatter)
	log.SetReportCaller(true)

	setLogLevel(logLevel, log)

	return &Logger{log: log}
}

func setLogLevel(logLevel string, log *logrus.Logger) {
	switch strings.ToLower(logLevel) {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.DebugLevel)
	}
}

func (l *Logger) LogInfo(c *gin.Context, statusCode int, message string) {
	logrus.WithFields(logrus.Fields{
		"code":         statusCode,
		"X-Process-ID": c.GetHeader("X-Process-ID"),
		"path":         c.Request.RequestURI,
		"message":      message,
	}).Info(message)
}

func (l *Logger) LogError(c *gin.Context, statusCode int, message string, err error) {
	logrus.WithFields(logrus.Fields{
		"code":         statusCode,
		"X-Process-ID": c.GetHeader("X-Process-ID"),
		"path":         c.Request.RequestURI,
		"message":      err.Error(),
	}).Error(message)
}

func (l *Logger) LogFatal(c *gin.Context, statusCode int, message string, err error) {
	logrus.WithFields(logrus.Fields{
		"code":    statusCode,
		"path":    c.Request.RequestURI,
		"message": err.Error(),
	}).Fatal(message)
}

func (l *Logger) LogDebug(c *gin.Context, message string) {
	logrus.WithFields(logrus.Fields{
		"X-Process-ID": c.GetHeader("X-Process-ID"),
		"path":         c.Request.RequestURI,
	}).Debug(message)
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.log.Debugf(msg, args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.log.Panicln(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log.Fatalln(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.log.Infoln(args...)
}

func (l *Logger) Errorf(error string, args ...interface{}) {
	l.log.Errorf(error, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.log.Warnln(args...)
}

func (l *Logger) Println(args ...interface{}) {
	l.log.Println(args...)
}
