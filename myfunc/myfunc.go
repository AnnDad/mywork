package myfunc

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/shopspring/decimal"

	//"sort"
	"io/ioutil"

	"math"
	"os/exec"
	"path"

	guid "github.com/typa01/go-utils"
)

type ArrayStr []string

type TDoWithFile func(file_path string)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type TMap map[string]interface{}

func Empty_Map() TMap {
	return TMap{}
}

func NewError(err string) error {
	return errors.New(err)
}

func WrapError(msg string, err error) error {
	return NewError(msg + err.Error())
}

func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func Abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

func GetGuid() string {
	return guid.GUID()
}

func CalcPages(count, size int) int {
	page := Div(count, size)
	if Mod(count, size) > 0 {
		page = page + 1
	}
	return page
}

func BooleanDef(def bool, arr ...bool) bool {
	result := def
	if len(arr) > 0 {
		result = arr[0]
	}
	return result
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Bool2Str(value bool) string {
	if value {
		return "true"
	} else {
		return "false"
	}
}

func QuoteS(str string) string {
	return "'" + str + "'"
}
func GetTickCount() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetNow_Unix() int64 {
	return time.Now().Unix()
}

func Mod(int1, int2 int) int {
	return (int1 % int2)
}
func Div(int1, int2 int) int {
	return (int1 / int2)
}

func Str2MD5(source string) string {
	d := []byte(source)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

func FormatNumber(length, n int) string {
	return fmt.Sprintf("%0"+String(length)+"d", n)
}

func Round2(num float64) float64 {
	var num2 float64
	decimalValue := decimal.NewFromFloat(0)
	if num < 0 {
		decimalValue = decimal.NewFromFloat(num - 0.005)

	} else {
		decimalValue = decimal.NewFromFloat(num + 0.005)
	}
	//乘100
	decimalValue = decimalValue.Mul(decimal.NewFromInt(100))
	res, _ := decimalValue.Float64()
	num3 := math.Trunc(res)
	decimalValue2 := decimal.NewFromFloat(num3)
	//除100
	decimalValue2 = decimalValue2.Div(decimal.NewFromInt(100))
	num2, _ = decimalValue2.Float64()
	return num2
}

func Split(source, sp_char string) ArrayStr {
	if strings.Trim(source, " ") == "" {
		return Empty_ArrayStr()
	}
	return strings.Split(source, sp_char)
}

func SetProxy_All(sock5_proxy string, debug ...bool) string {
	if len(debug) > 0 {
		if debug[0] {
			fmt.Println("SetProxy: ", sock5_proxy)
		}
	}
	err := os.Setenv("ALL_PROXY", sock5_proxy)
	if err != nil {
		return "err_ALL_PROXY: " + err.Error()
	}
	err = os.Setenv("HTTP_PROXY", sock5_proxy)
	if err != nil {
		return "err_HTTP_PROXY: " + err.Error()
	}
	err = os.Setenv("HTTPS_PROXY", sock5_proxy)
	if err != nil {
		return "err_HTTPS_PROXY: " + err.Error()
	}
	err = os.Setenv("NO_PROXY", "127.0.0.1;localhost;")
	return "ok"
}

func SetProxy_No(no_proxy string, debug ...bool) string {
	//fmt.Println(os.Getenv("NO_PROXY"))
	os.Setenv("NO_PROXY", no_proxy)
	return "ok"
}

func AppParams(index int) string {

	if index >= 0 && index < AppParamsCount() {
		return os.Args[index]
	} else {
		return ""
	}
}

func HttpProtocolSwap(url string) string {
	result := LowerCase(url)
	if LeftStr(result, 5) == "http:" || LeftStr(result, 6) == "https:" {
		if LeftStr(result, 5) == "http:" {
			return "https:" + SubStrRune(url, 5, 0)
		} else {
			return "http:" + SubStrRune(url, 6, 0)
		}

	} else {
		return "https://" + url
	}
}

func HttpProtocolExists(url string) bool {
	_url := LowerCase(url)
	return LeftStr(_url, 5) == "http:" || LeftStr(_url, 6) == "https:"
}

func AppParamsALL() []string {
	return os.Args
}

func AppParamsCount() int {
	return len(AppParamsALL())
}

func AppParamsExists() bool {
	return AppParamsCount() >= 2
}

func FileSize(path string) int64 {
	var result int64 = 0

	fi, err := os.Stat(path)
	if err == nil {
		result = fi.Size()

	}
	return result
}

func ThisName() string {
	return filepath.Base(ThisPath())
}

func ThisPath() string {
	selfPath, _ := exec.LookPath(os.Args[0])
	if selfPath != "" {
		selfPath, _ = filepath.Abs(selfPath)
	}
	if selfPath == "" {
		selfPath, _ = filepath.Abs(os.Args[0])
	}
	return selfPath
}

func IsWindows() bool {
	sysType := runtime.GOOS
	return sysType == "windows"
}

func ConvertRecoverError(RecoverErr interface{}) (err error) {
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

func ProtectDo(do func(), onErr func(err error)) {
	defer func() {
		if r := recover(); r != nil {
			err := ConvertRecoverError(r)
			if onErr != nil {
				onErr(err)
			} else {
				Println("Err_ProtectRun: ", err.Error())
			}
		}
	}()
	do()
}

func ProtectDoWithParams(do func(params ...interface{}), onErr func(err error), params ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			err := ConvertRecoverError(r)
			if onErr != nil {
				onErr(err)
			} else {
				Println("Err_ProtectRun: ", err.Error())
			}
		}
	}()
	do(params...)
}

func PathAutoSys(path string) string {
	return filepath.FromSlash(path)
}

func Replace(source, oldstr, newstr string) string {
	return strings.Replace(source, oldstr, newstr, -1)
}

func EndWithSlash(source string) string {
	result := strings.TrimRight(source, "/")
	return result + "/"

}

//UNICODE相关
type UnicodeStr []rune

func Str2Unicode(source string) UnicodeStr {
	return []rune(source)
}

func (this UnicodeStr) Char(index int) string {
	return string(this[index])
}

func (this UnicodeStr) Length() int {
	return len(this)
}

func (this UnicodeStr) Substr(tag1, tag2 string) UnicodeStr {

	return Empty_UnicodeStr()
}

func Empty_UnicodeStr() UnicodeStr {
	return make(UnicodeStr, 0)
}

func (this UnicodeStr) Copy(start, length int) string {
	if length == 0 {
		return ""
	}
	len_str := len(this)

	if start < 0 {
		start = len_str + start
	}
	if start > len_str {
		start = len_str
	}
	end := start + length
	if end > len_str {
		end = len_str
	}
	if length < 0 {
		end = len_str + length
	}
	if start > end {
		start, end = end, start
	}
	return string(this[start:end])
}

func (this *ArrayStr) Length() int {
	return len(*this)
}

func (this *ArrayStr) Clear() {
	*this = Empty_ArrayStr()
}

func (this *ArrayStr) Remove(index int) bool {
	result := false
	n := this.Length() - 1
	if index <= n {
		*this = append((*this)[:index], (*this)[index+1:]...)
		result = true
	}
	return result
}

func (this *ArrayStr) Exists(str string) bool {
	return StrInArray(str, *this)
}

func (this *ArrayStr) CountRange(range1, range2 string) int {
	result := 0
	n := this.Length()
	for i := 0; i < n; i++ {
		if (*this)[i] >= range1 && (*this)[i] <= range2 {
			result = result + 1
		}
	}
	return result
}

func (this *ArrayStr) RemoveLess(than string) int {
	result := *this
	*this = Empty_ArrayStr()
	n := result.Length()
	deleted := 0
	for i := 0; i < n; i++ {
		if result[i] >= than {
			*this = append(*this, result[i])
		} else {
			deleted = deleted + 1
		}
	}
	return deleted
}

func (this *ArrayStr) Join(char string) string {
	result := ""
	n := this.Length()
	for i := 0; i < n; i++ {
		result = AddStrWithDelimiter(result, (*this)[i], char)
	}
	return result
}

func (this *ArrayStr) String() string {
	result := ""
	n := this.Length()
	for i := 0; i < n; i++ {
		if result != "" {
			result = result + ","
		}
		result = result + "\"" + (*this)[i] + "\""
	}
	result = "[" + result + "]"
	return result
}

func (this *ArrayStr) Add(str string, CheckExists ...bool) bool {
	check := true
	if len(CheckExists) > 0 {
		check = CheckExists[0]
	}
	canAdd := true
	if check {
		canAdd = !this.Exists(str)
	}
	if canAdd {
		*this = append(*this, str)
		return true
	}
	return false
}

func Empty_ArrayStr() ArrayStr {
	return make(ArrayStr, 0)
}

func Empty_Bytes() []byte {
	return make([]byte, 0)
}

func StrInArray(str string, str_arr []string) bool {
	n := len(str_arr)
	for i := 0; i < n; i++ {
		if str_arr[i] == str {
			return true
		}
	}
	return false
}

func StrInStrs(str string, strs ...string) bool {
	n := len(strs)
	for i := 0; i < n; i++ {
		if strs[i] == str {
			return true
		}
	}
	return false
}

func Contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

func Join(strs []string, spliter string) string {
	result := ""
	n := len(strs)
	for i := 0; i < n; i++ {
		result = AddStrWithDelimiter(result, strs[i], spliter)
	}
	return result
}

func AddLeftStrIfNotExists(source, str string) string {
	if Contains(source, str) {
		return source
	} else {
		return str + source
	}
}
func AddRightStrIfNotExists(source, str string) string {
	if Contains(source, str) {
		return source
	} else {
		return source + str
	}
}
func ClearDir(APath string) {
	dir, _ := ioutil.ReadDir(APath)
	for _, d := range dir {
		item := path.Join([]string{APath, d.Name()}...)
		os.RemoveAll(item)
	}
}

func RemoveFile(filePath string) {
	if ExistPath(filePath) {
		os.Remove(filePath)
	}
}

func CreatePath(filePath string) error {
	if !ExistPath(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func Str2Reader(str string) *strings.Reader {
	return strings.NewReader(str)
}

func CopyFile(src, dst string, BUFFERSIZE int64, override bool) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file.", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	if !override {
		_, err = os.Stat(dst)
		if err == nil {
			return fmt.Errorf("File %s already exists.", dst)
		}
	}
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	if err != nil {
		panic(err)
	}
	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

func ArrToStr(spliter string, arr ...interface{}) string {
	result := ""
	for _, v := range arr {
		if result != "" {
			result = result + spliter
		}
		result = result + String(v)
	}
	return result
}

func CMD(cmd string, params ...string) error {
	command := exec.Command(cmd, params...)
	return command.Start()
}

func CMD2(cmd string, params ...string) error {
	command := exec.Command(cmd, params...)
	command.Stdout = os.Stdout
	return command.Run()

}

func OpenInExplorer(path string) {
	CMD("cmd", "/c", "explorer", "/select,"+path+"")
}
func OpenWithDefault(path string) {
	CMD("cmd", "/c", "explorer", "/open,"+path+"")
}
func ShellCmd(cmd string) (string, error) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	command := exec.Command("/bin/bash", "-c", cmd)
	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	command.Stdout = &out
	//Run执行c包含的命令，并阻塞直到完成。这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := command.Run()
	return out.String(), err
}

func OpenURL(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/C", "start", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {

	}
}

func AppPath(path ...string) string {
	dir := getCurrentAbPathByExecutable()
	if strings.Contains(dir, getTmpDir()) {
		dir = getCurrentAbPathByCaller()
	}
	if len(path) > 0 {
		if path[0] != "" {
			dir = dir + "/" + path[0]
		}
	}
	dir = PathAutoSys(dir + "/")
	if !ExistPath(dir) {
		CreatePath(dir)
	}
	return dir
}

func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		panic("getCurrentAbPathByExecutable: " + err.Error())
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

// 获取系统临时目录，兼容go run
func getTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

func Milliseconds(milliseconds int) time.Duration {
	return time.Duration(milliseconds) * time.Millisecond
}
func Seconds(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}

func Minutes(minutes int) time.Duration {
	return time.Duration(minutes) * time.Minute
}
func Hours(hours int) time.Duration {
	return time.Duration(hours) * time.Hour
}
func Days(days int) time.Duration {
	return time.Duration(days) * time.Hour * 24
}

func ExistPath(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func SaveFile(isnew bool, path, text string, writeRN ...bool) (err error) {
	newfile := isnew
	if !newfile {
		newfile = !ExistPath(path)
	}

	var f *os.File

	if newfile {
		f, err = os.Create(path)
	} else {
		f, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	}
	defer f.Close()
	if isnew && text == "" {
		//如果是新建并且文本为空, 那什么都不写入
	} else {
		_writeRN := true
		if len(writeRN) > 0 {
			_writeRN = writeRN[0]
		}
		if _writeRN {
			_, err = f.WriteString(text + "\r\n")
		} else {
			_, err = f.WriteString(text)
		}
	}

	return err
}

func ReadFileAsBytes(path string) ([]byte, error) {
	if ExistPath(path) {
		return ioutil.ReadFile(path)
	} else {
		return nil, NewError("Err_Not_Exists_Path: " + path)
	}
}

func ReadFileAsString(path string) (string, error) {
	bytes, err := ReadFileAsBytes(path)
	if err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}

func SaveFileBytes(isnew bool, path string, bytes []byte) (err error) {
	newfile := isnew
	if !newfile {
		newfile = !ExistPath(path)
	}

	var f *os.File
	defer f.Close()

	if newfile {
		f, err = os.Create(path)
	} else {
		f, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	}
	_, err = f.Write(bytes)
	return err
}

func GetLocalIP() ArrayStr {
	var result ArrayStr

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		result = append(result, "")
		return result
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if !result.Exists(ipnet.IP.String()) {
					result = append(result, ipnet.IP.String())
				}
			}
		}
	}
	return result
}

func Bool2StrEx(exp bool, value_true, value_false string) string {
	if exp {
		return value_true
	} else {
		return value_false
	}
}

func CaseWhenInt(exp bool, value_true, value_false int) int {
	if exp {
		return value_true
	} else {
		return value_false
	}
}
func CaseWhenStr(exp bool, value_true, value_false string) string {
	if exp {
		return value_true
	} else {
		return value_false
	}
}

func Switch(exp bool, value_true, value_false string) string {
	if exp {
		return value_true
	} else {
		return value_false
	}
}

func Random(min, max int) int {

	return rand.Intn(max+1-min) + min
}

func SleepSeconds(n int) {
	time.Sleep(Seconds(n))
}
func SleepMinutes(n int) {
	time.Sleep(Minutes(n))
}
func SleepMilliSeconds(n int) {
	time.Sleep(Milliseconds(n))
}

func SleepSecondsRandom(r1, r2 int, tip ...string) {
	n := Random(r1, r2)
	for i := 0; i < n; i++ {
		if len(tip) > 0 {
			if tip[0] != "" {
				Println(tip[0], n-i)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func Println(a ...interface{}) {
	var b []interface{}
	b = append(b, a...)
	fmt.Println(b...)
}

func PrintlnTip(tip string, a ...interface{}) {
	var b []interface{}
	b = append(b, tip)
	b = append(b, a...)
	fmt.Println(b...)
}

func KeyInMap(obj map[string]interface{}, key string) bool {
	_, ok := obj[key]
	return ok
}

func MapValue(obj map[string]interface{}, field string) string {
	v, ok := obj[field]
	if ok {
		return String(v)
	} else {
		return ""
	}
}

type apiString interface {
	String() string
}

// apiError is used for type assert api for Error().
type apiError interface {
	Error() string
}

func String(i interface{}) string {
	if i == nil {
		return ""
	}
	switch value := i.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	case time.Time:
		if value.IsZero() {
			return ""
		}
		return value.String()
	case *time.Time:
		if value == nil {
			return ""
		}
		return value.String()
	/*
		case gtime.Time:
			if value.IsZero() {
				return ""
			}
			return value.String()
		case *gtime.Time:
			if value == nil {
				return ""
			}
			return value.String()
	*/
	default:
		// Empty checks.
		if value == nil {
			return ""
		}
		if f, ok := value.(apiString); ok {
			// If the variable implements the String() interface,
			// then use that interface to perform the conversion
			return f.String()
		}
		if f, ok := value.(apiError); ok {
			// If the variable implements the Error() interface,
			// then use that interface to perform the conversion
			return f.Error()
		}
		// Reflect checks.
		var (
			rv   = reflect.ValueOf(value)
			kind = rv.Kind()
		)
		switch kind {
		case reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return ""
			}
		case reflect.String:
			return rv.String()
		}
		if kind == reflect.Ptr {
			return String(rv.Elem().Interface())
		}
		// Finally we use json.Marshal to convert.
		if jsonContent, err := json.Marshal(value); err != nil {
			return fmt.Sprint(value)
		} else {
			return string(jsonContent)
		}
	}
}

const (
	// NotFoundIndex is the position index for string not found in searching functions.
	NotFoundIndex = -1
)

func LowerCase(source string) string {
	return strings.ToLower(source)
}

func UpperCase(source string) string {
	return strings.ToUpper(source)
}

func IsMonth(month string) bool {
	fmonth := LowerCase(month)
	if StrInStrs(fmonth, "january", "february", "march", "april", "may", "june",
		"july", "august", "september", "october", "november", "december") {
		return true
	}
	return false
}

func MonthNameToNum(month string) int {
	result := 0
	switch LowerCase(month) {
	case "january":
		result = 1
	case "february":
		result = 2
	case "march":
		result = 3
	case "april":
		result = 4
	case "may":
		result = 5
	case "june":
		result = 6
	case "july":
		result = 7
	case "august":
		result = 8
	case "september":
		result = 9
	case "october":
		result = 10
	case "november":
		result = 11
	case "december":
		result = 12
	default:
		panic("Err_MonthNameToNum: " + month)
	}
	return result
}

func LeftStr(str string, length int) string {
	return SubStrRune(str, 0, length)
}

func RightStr(str string, length int) string {
	return SubStrRune(str, -length, length)
}

func ExtractRoot(ADomain string) string {
	result := ADomain
	result = strings.Replace(result, "http://", "", 1)
	result = strings.Replace(result, "https://", "", 1)
	result = strings.Replace(result, "www.", "", 1)
	if Contains(result, "/") {
		result = USubstrByTag(result, "/")
	}
	return result
}

func MergeUrl(url, url2 string) string {
	if RightStr(url, 1) == "/" {
		url = strings.TrimRight(url, "/")
	}
	if LeftStr(url2, 1) == "/" {
		url2 = strings.TrimLeft(url2, "/")
	}
	return url + "/" + url2
}

func ExtractCurrentPath(url string) string {
	//Println("ExtractCurrentPath_url: ", url)
	_countMax := 0
	if HttpProtocolExists(url) {
		_countMax = 2
	}
	_count := strings.Count(url, "/")
	urls := Split(url, "/")
	if _count > _countMax {
		urls.Remove(urls.Length() - 1)
	}
	//Println("ExtractCurrentPath_join: ", urls.Join("/"))
	return urls.Join("/") + "/"
}

func ExtractFilePath(source string) string {
	return PathAutoSys(filepath.Dir(source) + `\`)
}
func ExtractFileName(source string) string {
	return filepath.Base(source)
}
func ExtractFileExt(source string) string {
	ext := filepath.Ext(source)
	if ext != "" {
		ext = strings.ToLower(ext)
	}
	return ext
}

func SubStrRune(str string, start int, length int) (substr string) {
	// Converting to []rune to support unicode.
	var (
		runes       = []rune(str)
		runesLength = len(runes)
	)
	var end int
	// Simple border checks.

	if start < 0 {
		start = runesLength + start
	}

	if length == 0 {
		end = runesLength
	} else {
		if length < 0 {
			end = runesLength + length
		} else {
			end = start + length
		}
	}

	//最终范围检查
	if start < 0 {
		start = 0
	}
	if start > runesLength {
		start = runesLength
	}
	if end > runesLength {
		end = runesLength
	}
	if end < start {
		end = start
	}
	//fmt.Println("len: ",runesLength," ,start: ",start, " end: ",end)
	return string(runes[start:end])
}

func USubstr(Source, tag1, tag2 string, addtag ...bool) string {
	result := ""
	inctag := false
	if len(addtag) > 0 {
		inctag = addtag[0]
	}
	if n1 := UPos(Source, tag1, 0); n1 >= 0 {
		if !inctag {
			n1 = n1 + ULen(tag1)
		}
		tmp := SubStrRune(Source, n1, 0)
		if n2 := UPos(tmp, tag2, 0); n2 >= 0 {
			if inctag {
				n2 = n2 + ULen(tag2)
			}
			result = SubStrRune(tmp, 0, n2)
		}
	}
	return result
}

func USubstrByTag(Source, tag string, addtag ...bool) string {
	result := ""
	inctag := false
	if len(addtag) > 0 {
		inctag = addtag[0]
	}
	n1 := UPos(Source, tag, 0)

	if n1 >= 0 {
		if inctag {
			n1 = n1 + ULen(tag)
		}
		result = SubStrRune(Source, 0, n1)
	}

	return result
}

func Pos(source, substr string, startOffset ...int) int {
	length := len(source)
	offset := 0
	if len(startOffset) > 0 {
		offset = startOffset[0]
	}
	if length == 0 || offset > length || -offset > length {
		return -1
	}
	if offset < 0 {
		offset += length
	}
	str := source[offset:]

	pos := strings.Index(str, substr)
	if pos == NotFoundIndex {
		return NotFoundIndex
	}
	return pos + offset
}

func UPos(source, substr string, startOffset ...int) int {
	pos := Pos(source, substr, startOffset...)

	if pos < 3 {
		return pos
	}
	return len([]rune(source[:pos]))
}

func ULen(str string) int {
	return utf8.RuneCountInString(str)
}

//filepath.WalkFunc: func(path string, info os.FileInfo, err error) error {}
func EnumDir(path string, do TDoWithFile) error {
	rd, err := ioutil.ReadDir(path)
	for _, fi := range rd {
		if fi.IsDir() {
			EnumDir(path+fi.Name()+string(os.PathSeparator), do)
		} else {
			do(path + fi.Name())
		}
	}
	return err
}

// UnsafeStrToBytes converts string to []byte without memory copy.
// Note that, if you completely sure you will never use `s` variable in the feature,
// you can use this unsafe function to implement type conversion in high performance.
func UnsafeStrToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// UnsafeBytesToStr converts []byte to string without memory copy.
// Note that, if you completely sure you will never use `b` variable in the feature,
// you can use this unsafe function to implement type conversion in high performance.
func UnsafeBytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func GetDirAllEntryPaths(dirname string, incDir bool) ([]string, error) {
	// Remove the trailing path separator if dirname has.
	dirname = strings.TrimSuffix(dirname, string(os.PathSeparator))

	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(infos))
	// Include current dir.
	if incDir {
		paths = append(paths, dirname)
	}

	for _, info := range infos {
		path := dirname + string(os.PathSeparator) + info.Name()
		if info.IsDir() {
			tmp, err := GetDirAllEntryPaths(path, incDir)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			continue
		}
		paths = append(paths, path)
	}
	return paths, nil
}

// GetDirAllEntryPathsFollowSymlink gets all the file or dir paths in the specified directory recursively.
func GetDirAllEntryPathsFollowSymlink(dirname string, incDir bool) ([]string, error) {
	// Remove the trailing path separator if dirname has.
	dirname = strings.TrimSuffix(dirname, string(os.PathSeparator))

	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(infos))
	// Include current dir.
	if incDir {
		paths = append(paths, dirname)
	}

	for _, info := range infos {
		path := dirname + string(os.PathSeparator) + info.Name()
		realInfo, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if realInfo.IsDir() {
			tmp, err := GetDirAllEntryPathsFollowSymlink(path, incDir)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			continue
		}
		paths = append(paths, path)
	}
	return paths, nil
}

//在字符串尾部添加新的字符串及间隔符
func AddStrWithDelimiter(items, item, delimiter string) string {
	if items == "" {
		return item
	} else {
		return items + delimiter + item
	}
}

func TypeOf(v interface{}) string {
	return fmt.Sprintf("{%T}:", v)
}

func DumpArry(v []interface{}) {

	for key, value := range v {
		fmt.Printf("key: " + String(key) + " value: ")
		fmt.Println(value)
	}
}

func TrimStr(source string, char ...string) string {
	chr := " "
	if len(char) > 0 {
		chr = char[0]
	}
	return strings.TrimRight(strings.TrimLeft(source, chr), chr)
}

func DumpArryStr(v []string) {
	for key, value := range v {
		fmt.Printf("key: " + String(key) + " value: ")
		fmt.Println(value)
	}
}
