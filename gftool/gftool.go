package gftool

import (
	"github.com/anndad/mywork/myfunc"
	"time"

	// "crypto/md5"
	// "encoding/base64"
	// "fmt"
	// "io"
	// "math/rand"
	// "os"
	// syspath "path"
	// "path/filepath"
	// "runtime"
	"strconv"
	// "strings"
	// "time"
	// "unsafe"
	// "github.com/axgle/mahonia"
	"github.com/gogf/gf/frame/g"
	// "github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"

	//"github.com/gogf/gf/database/gredis"
	// "github.com/typa01/go-utils"
	// "golang.org/x/net/html/charset"
	//"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/gtime"
	//"github.com/gogf/gf/encoding/gjson"
)

type TMap = g.Map

func GetGuid() string {
	return myfunc.GetGuid()
}
func GetTickount() int64 {
	return gtime.Millisecond()
}

func GetTimestamp_second() int64 {
	return gtime.Second()
}

func SecondsBetween(datetime1, datetime2 string) int64 {
	d1 := gtime.New(datetime1).Unix()
	d2 := gtime.New(datetime2).Unix()
	result := myfunc.Abs(d2 - d1)
	return result
}

func SecondsBetweenNow(datetime string) int64 {
	return SecondsBetween(datetime, Now2Str())
}

func MinutesBetween(datetime1, datetime2 string) int64 {
	return SecondsBetween(datetime1, datetime2) / 60
}

func MinutesBetweenNow(datetime string) int64 {
	return MinutesBetween(datetime, Now2Str())
}

func HoursBetween(datetime1, datetime2 string) int64 {
	return SecondsBetween(datetime1, datetime2) / 3600
}

func HoursBetweenNow(datetime string) int64 {
	return HoursBetween(datetime, Now2Str())
}
func DaysBetween(str_day1, str_day2 string) int64 {
	return SecondsBetween(str_day1, str_day2) / 86400
}

func DaysBetweenToday(str_day string) int64 {
	return DaysBetween(Today2Str(), str_day)
}

func IntToStr(v int) string {
	return gconv.String(v)
}

func Int64ToStr(v int64) string {
	return gconv.String(v)
}

func Interface2Str(v interface{}) string {
	return gconv.String(v)
}

func Any2Str(v interface{}) string {
	return gconv.String(v)
}

func Str2Int(v string) int {
	return gconv.Int(v)
}

func Str2Int64(v string) int64 {
	return gconv.Int64(v)
}

func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func Str2IntDef(v string, def int) int {
	if IsNum(v) {
		return Str2Int(v)
	} else {
		return def
	}
}
func Str2Int64Def(v string, def int64) int64 {
	if IsNum(v) {
		return Str2Int64(v)
	} else {
		return def
	}
}

func Str2Float32(v string) float32 {
	return gconv.Float32(v)
}

func StrDateTime2Unix(str_datetime string) int64 {
	return gtime.New(str_datetime).Unix()
}

func StrDateTime2UnixNano(str_datetime string) int64 {
	return gtime.New(str_datetime).UnixNano()
}

func StrDateTime2GTime(str_datetime string) *gtime.Time {
	return gtime.New(str_datetime)
}

func AddDate(source *gtime.Time, year, month, day int) *gtime.Time {
	return source.AddDate(year, month, day)
}

func AddTime(source *gtime.Time, hour, minute, second int) *gtime.Time {
	fix := time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute + time.Duration(second)*time.Second
	return source.Add(fix)
}

func EncodeDateS(year, month, day int) string {
	fyear := IntToStr(year)
	if len(fyear) == 2 {
		fyear = "20" + fyear
	}
	fmonth := IntToStr(month)
	if len(fmonth) == 1 {
		fmonth = "0" + fmonth
	}
	fday := IntToStr(day)
	if len(fday) == 1 {
		fday = "0" + fday
	}
	return fyear + "-" + fmonth + "-" + fday
}
func EncodeTimeS(hour, minute, second int) string {
	fhour := IntToStr(hour)
	if len(fhour) == 1 {
		fhour = "0" + fhour
	}
	fminute := IntToStr(minute)
	if len(fminute) == 1 {
		fminute = "0" + fminute
	}
	fsecond := IntToStr(second)
	if len(fsecond) == 1 {
		fsecond = "0" + fsecond
	}
	return fhour + ":" + fminute + ":" + fsecond
}
func EncodeTime(hour, minute, second int) *gtime.Time {
	return StrDateTime2GTime(EncodeTimeS(hour, minute, second))
}

func EncodeDate(year, month, day int) *gtime.Time {
	return StrDateTime2GTime(EncodeDateS(year, month, day))
}

func EncodeDateTime(year, month, day, hour, minute, second int) *gtime.Time {
	return StrDateTime2GTime(EncodeDateS(year, month, day) + " " + EncodeTimeS(hour, minute, second))
}

func Date2Str(source *gtime.Time) string {
	return source.Format("Y-m-d")
}

func YearOfDate(source *gtime.Time) int {
	return Str2Int(source.Format("Y"))
}

func Time2Str(source *gtime.Time) string {
	return source.Format("H:i:s")
}

func Minute2Str(source *gtime.Time) string {
	return source.Format("H:i:00")
}

func DateTime2Str(source *gtime.Time) string {
	return source.Format("Y-m-d H:i:s")
}

func DateTime2Week(source *gtime.Time) int {
	return Str2Int(source.Format("w"))
}

func StrDateTime2Week(source string) int {
	return Str2Int(StrDateTime2GTime(source).Format("w"))
}

func Now() *gtime.Time {
	return gtime.Now()
}

func Today() *gtime.Time {
	return gtime.New(gtime.Date())
}

func Now2Str() string {
	return DateTime2Str(Now())
}

func Now2StrWithMS() string {
	return Now().Format("Y-m-d H:i:s.u")
}

func Today2Str() string {
	return gtime.Date()
}

func MB_Pos(str, substr string, offset int) int {
	foffset := offset
	if offset > 0 {
		item := gstr.SubStrRune(str, 0, offset)
		foffset = len(item)
	}
	return gstr.PosRune(str, substr, foffset)
}

func MB_Substr(str, tag1, tag2 string, offset int, addtag bool) (string, int) {

	n1 := MB_Pos(str, tag1, offset)
	if n1 < 0 {
		return "", -1
	}
	n1 = n1 + gstr.LenRune(tag1)
	n2 := MB_Pos(str, tag2, n1)
	if n2 < 0 {
		return "", -1
	}
	result := gstr.SubStrRune(str, n1, n2-n1)
	if addtag {
		result = tag1 + result + tag2
	}
	return result, n2 + gstr.LenRune(tag2)
}

func MB_DelSubstr(str, tag1, tag2 string, offset int) (string, int) {
	n1 := MB_Pos(str, tag1, offset)
	if n1 < 0 {
		return "", -1
	}
	n1 = n1 + gstr.LenRune(tag1)
	n2 := MB_Pos(str, tag2, n1)
	if n2 < 0 {
		return "", -1
	}
	return gstr.SubStrRune(str, n1, n2-n1), n2 + gstr.LenRune(tag2)
}
