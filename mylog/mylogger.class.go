package mylog

import (
	"anndad/gftool"
	"anndad/myfunc"
	"fmt"

	"errors"
	"strings"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gregex"
)

type TMylogger struct {
	isScreenOut bool
	RotateSize  int64
	conf        glog.Config
	logger      *glog.Logger
}

func Newlogger(name string) *TMylogger {
	fname := name
	if fname == "" {
		fname = myfunc.ThisName()
	}
	result := new(TMylogger)
	result.isScreenOut = true
	result.logger = g.Log(fname)
	Path := myfunc.AppPath("logs")
	Println("logPath:", Path)
	result.conf = glog.DefaultConfig()
	result.conf.Path = Path
	//取消默认的时间
	result.conf.Flags = 0
	result.conf.Prefix = ""
	result.conf.File = fname + ".{Ymd}.{H}.log"
	result.conf.StdoutPrint = false

	//切分为30M
	result.conf.RotateSize = 0
	//result.conf.RotateBackupLimit = 9
	result.conf.RotateBackupExpire = time.Hour * 72
	result.logger.SetConfig(result.conf)
	result.RotateSize = 1024 * 10
	return result
}

func (this *TMylogger) logFilePath(now *gtime.Time) string {
	// Content containing "{}" in the file name is formatted using gtime.
	file, _ := gregex.ReplaceStringFunc(`{.+?}`, this.conf.File, func(s string) string {
		return now.Format(strings.Trim(s, "{}"))
	})
	file = gfile.Join(this.conf.Path, file)
	return file
}

func (this *TMylogger) Options_SaveLog() bool {
	return g.Config("").GetInt("log_options.save_log", 1) == 1
}
func (this *TMylogger) Options_SaveDetail() bool {
	if this.Options_SaveLog() {
		return g.Config("").GetInt("log_options.save_doc", 0) == 1
	} else {
		return false
	}

}

func (this *TMylogger) doRotate() {
	return
	if this.RotateSize > 0 {
		path := this.logFilePath(gftool.Now())
		bytes := myfunc.FileSize(path)
		if bytes >= this.RotateSize {
			r := gfile.Rename(path, this.conf.Path+gftool.Now2StrWithMS()+".log")
			fmt.Println("doRotate", r)
			myfunc.SleepSeconds(1)
		} else {
			fmt.Println("less size: ", bytes)
		}
	}
}

func (this *TMylogger) SetBackupLimit(n int) {
	this.conf.RotateBackupLimit = n
	this.logger.SetConfig(this.conf)
}

func (this *TMylogger) SetBackupExpire(n time.Duration) {
	this.conf.RotateBackupExpire = n
	this.logger.SetConfig(this.conf)
}

func (this *TMylogger) SetScreenOut(value bool) {
	this.isScreenOut = value
}

func (this *TMylogger) SetPrefix(prefix string) {
	this.logger.SetPrefix("(" + prefix + ")>")
}

func (this *TMylogger) GetGuid() string {
	return myfunc.GetGuid()
}

func (this *TMylogger) String(i interface{}) string {
	return myfunc.String(i)
}

func (this *TMylogger) Info(info string, docs ...string) {
	msg := gftool.Now2StrWithMS() + "> " + info
	if len(docs) > 0 {
		fpath := "doc_" + this.GetGuid() + ".txt"
		msg = msg + "[" + fpath + "]"
		fpath = myfunc.AppPath("logs/docs") + fpath
		if this.Options_SaveDetail() {
			myfunc.SaveFile(true, fpath, docs[0])
		} else {
			myfunc.Println(docs)
		}
	}
	this.StdoutPrint(msg)
	this.doRotate()
	if this.Options_SaveLog() {
		this.logger.Info(msg)
	}

}

func (this *TMylogger) Error(tip string, err error) {
	this.Info(tip + err.Error())
	this.logger.Error(err)

}

func (this *TMylogger) Fatal(tip string, err error) {
	this.Info(tip + err.Error())
	this.logger.Fatal(err)
}

func (this *TMylogger) Debug(tip string) {
	this.Info(tip)
	this.logger.Debug(tip)
}

func (this *TMylogger) Clear() {
	myfunc.ClearDir(this.logger.GetPath())
}

func (this *TMylogger) ScreenPrint(tip string) {
	fmt.Println(tip)
}

func (this *TMylogger) StdoutPrint(tip string) {
	if this.isScreenOut {
		this.ScreenPrint(tip)
	}
}

func (this *TMylogger) Export(i ...interface{}) string {
	tip := g.Export(i...)
	return tip
}

func (this *TMylogger) ProtectDoWithLog(do func()) {
	defer func() {
		if r := recover(); r != nil {
			err := this.ConvertRecoverError(r)
			this.Error("Err_ProtectRun: ", err)
		}
	}()
	do()
}

func (this *TMylogger) ConvertRecoverError(RecoverErr interface{}) (err error) {
	switch x := RecoverErr.(type) {
	case string:
		err = errors.New(x)
	case error:
		err = x
	default:
		err = errors.New("unknow error")
	}
	return err
}

func (this *TMylogger) CatchError(ATip string, OnErr func(ainfo string, err error)) {

	if r := recover(); r != nil {
		err := this.ConvertRecoverError(r)
		OnErr(ATip, err)
	}

}

func (this *TMylogger) Println(a ...interface{}) {
	var b []interface{}
	b = append(b, gftool.Now2StrWithMS())
	b = append(b, a...)
	fmt.Println(b...)
}

func (this *TMylogger) DumpArryStr(v []string) {
	for key, value := range v {
		msg := "key: " + String(key) + " value: " + value
		this.Info(msg)
	}
}

func (this *TMylogger) Dump(tip string, v ...interface{}) {
	n := len(v)
	fmt.Println("=====开始====="+tip+"======长度[", n, "]========")
	for i := 0; i < n; i++ {
		vtype := fmt.Sprintf("%T", v[i])
		fmt.Println("[", i, "]类型(", vtype, ")=>", v[i], vtype == "[]string")
		fmt.Println("-----------------------")
	}
	fmt.Println("=====结束=====" + tip + "==============")
}
