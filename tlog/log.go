package tlog

import (
	"github.com/astaxie/beego/logs"
	"strings"
)

var Log *logs.BeeLogger

func init() {
	Log = logs.NewLogger(10000)
	Log.SetLogger("console", "")
	Log.EnableFuncCallDepth(true)
	Log.SetLogFuncCallDepth(3)
}

// SetLevel sets the global log level used by the simple logger.
func SetLogLevel(l int) {
	Log.SetLevel(l)
}

// SetLogFuncCall set the CallDepth, default is 3
func SetLogFuncCall(b bool) {
	Log.EnableFuncCallDepth(b)
	Log.SetLogFuncCallDepth(3)
}

// SetLogger sets a new logger.
func SetLogger(adaptername string, config string) error {
	err := Log.SetLogger(adaptername, config)
	if err != nil {
		return err
	}
	return nil
}

// Emergency logs a message at emergency level.
func Emergency(v ...interface{}) {
	Log.Emergency(generateFmtStr(len(v)), v...)
}

// Alert logs a message at alert level.
func Alert(v ...interface{}) {
	Log.Alert(generateFmtStr(len(v)), v...)
}

// Critical logs a message at critical level.
func Critical(v ...interface{}) {
	Log.Critical(generateFmtStr(len(v)), v...)
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	Log.Error(generateFmtStr(len(v)), v...)
}

func LOGE(v ...interface{}) {
	Log.Error(generateFmtStr(len(v)), v...)
}

// Warning logs a message at warning level.
func Warning(v ...interface{}) {
	Log.Warning(generateFmtStr(len(v)), v...)
}

// Warn compatibility alias for Warning()
func Warn(v ...interface{}) {
	Log.Warn(generateFmtStr(len(v)), v...)
}

func LOGW(v ...interface{}) {
	Log.Warn(generateFmtStr(len(v)), v...)
}

// Notice logs a message at notice level.
func Notice(v ...interface{}) {
	Log.Notice(generateFmtStr(len(v)), v...)
}

// Informational logs a message at info level.
func Informational(v ...interface{}) {
	Log.Informational(generateFmtStr(len(v)), v...)
}

// Info compatibility alias for Warning()
func Info(v ...interface{}) {
	Log.Info(generateFmtStr(len(v)), v...)
}

func LOGI(v ...interface{}) {
	Log.Info(generateFmtStr(len(v)), v...)
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	Log.Debug(generateFmtStr(len(v)), v...)
}

func LOGD(v ...interface{}) {
	Log.Debug(generateFmtStr(len(v)), v...)
}

// Trace logs a message at trace level.
// compatibility alias for Warning()
func Trace(v ...interface{}) {
	Log.Trace(generateFmtStr(len(v)), v...)
}

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}
