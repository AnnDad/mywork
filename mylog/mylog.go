package mylog

import (
	"time"
)

var defLogger *TMylogger

func init() {
	defLogger = Newlogger("")
}

func SetBackupLimit(n int) {
	defLogger.SetBackupLimit(n)
}
func SetBackupExpire(n time.Duration) {
	defLogger.SetBackupExpire(n)
}

func SetScreenOut(Enable bool) {
	defLogger.SetScreenOut(Enable)
}

func SetPrefix(prefix string) {
	defLogger.SetPrefix(prefix)
}

func GetGuid() string {
	return defLogger.GetGuid()
}

func String(i interface{}) string {
	return defLogger.String(i)
}

func Info(info string, docs ...string) {
	defLogger.Info(info, docs...)
}

func Error(tip string, err error) {
	defLogger.Error(tip, err)

}

func Fatal(tip string, err error) {
	defLogger.Fatal(tip, err)
}

func Debug(tip string) {
	defLogger.Debug(tip)
}

func Clear() {
	defLogger.Clear()
}

func ScreenPrint(tip string) {
	defLogger.ScreenPrint(tip)
}

func Export(i ...interface{}) string {
	return defLogger.Export(i...)
}

func ProtectDoWithLog(do func()) {
	defLogger.ProtectDoWithLog(do)
}

func ConvertRecoverError(RecoverErr interface{}) (err error) {
	return defLogger.ConvertRecoverError(RecoverErr)
}

func CatchError(tip string, OnErr func(ainfo string, err error)) {

	defLogger.CatchError(tip, OnErr)

}

func Println(v ...interface{}) {
	defLogger.Println(v...)
}

func Dump(tip string, v ...interface{}) {
	defLogger.Dump(tip, v...)
}

func DumpArryStr(v []string) {
	defLogger.DumpArryStr(v)
}
