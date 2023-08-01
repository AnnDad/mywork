package mywebapp

import (
	//"anndad/mylog"
	//"log"

	"net/url"

	"github.com/anndad/mywork/myfunc"
	"github.com/anndad/mywork/myjson"
	"github.com/anndad/mywork/mylog"
	"github.com/gogf/gf/net/ghttp"
)

type TRSP struct {
	State bool        `json:"state"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}
type TLogFunc = func(req *ghttp.Request, fdata interface{})

var (
	LogRequest  bool = false
	LogFunction TLogFunc
)

func UrlEncode(s string) string {
	return url.QueryEscape(s)
}

func UrlDecode(s string) string {
	r, _ := url.QueryUnescape(s)
	return r
}

func RequestAsJson(req *ghttp.Request) *myjson.TJson {
	jsReq, err := req.GetJson()
	if err != nil {
		ReturnError(req, "err_request_not_valid_json")
		return nil
	}
	return jsReq
}

func RequestFieldCheck(req *ghttp.Request, js *myjson.TJson, field string) {
	if !js.Contains(field) {
		ReturnError(req, "err_miss_field: "+field)
		return
	}
}
func rsp(fstate bool, msg string, fdata interface{}) TRSP {
	return TRSP{State: fstate, Msg: msg, Data: fdata}
}

func defaultResponse(req *ghttp.Request, fdata interface{}, abort ...bool) {
	if LogFunction == nil {
		//默认日志
		reqID := req.GetHeader("Requestid")
		if reqID == "" {
			reqID = mylog.GetGuid()
		}
		mylog.Info("def[" + reqID + "]url: " + req.Request.RequestURI)
		mylog.Info("def[" + reqID + "]body: " + req.GetBodyString())
		mylog.Info("def[" + reqID + "]response: " + mylog.String(fdata))
	} else {
		LogFunction(req, fdata)
	}
	if myfunc.BooleanDef(true, abort...) {
		req.Response.WriteJsonExit(fdata)
	} else {
		req.Response.WriteJson(fdata)
	}
}

func ReturnData(req *ghttp.Request, fstate bool, msg string, fdata interface{}, abort ...bool) {
	defaultResponse(req, rsp(fstate, msg, fdata), abort...)
}

func ReturnError_Data(req *ghttp.Request, msg string, fdata interface{}, abort ...bool) {
	defaultResponse(req, rsp(false, msg, fdata), abort...)
}

func ReturnError(req *ghttp.Request, msg string, abort ...bool) {
	defaultResponse(req, rsp(false, msg, nil), abort...)
}

func ReturnErrorIf(req *ghttp.Request, tip string, err error, abort ...bool) {
	if err != nil {
		defaultResponse(req, rsp(false, tip+" "+err.Error(), nil), abort...)
	}
}

func ReturnOK(req *ghttp.Request, msg string, abort ...bool) {
	defaultResponse(req, rsp(true, msg, nil), abort...)
}

func ReturnOK_Data(req *ghttp.Request, fdata interface{}, abort ...bool) {
	defaultResponse(req, rsp(true, "ok", fdata), abort...)
}

func ReturnCustomData(req *ghttp.Request, data interface{}, abort ...bool) {
	defaultResponse(req, data, abort...)
}

func CORSDefault(req *ghttp.Request) {
	req.Response.CORSDefault()
	req.Middleware.Next()
}
