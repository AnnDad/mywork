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

var (
	LogRequest bool = false
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
func rsp(fstate bool, fmsg string, fdata interface{}) TRSP {
	return TRSP{State: fstate, Msg: fmsg, Data: fdata}
}

func defaultResponse(req *ghttp.Request, fdata interface{}, abort ...bool) {
	if LogRequest {
		reqID := req.GetHeader("requestID")
		if reqID == "" {
			reqID = mylog.GetGuid()
		}
		mylog.Info("[" + reqID + "]url: " + req.Request.RequestURI)
		mylog.Info("[" + reqID + "]body: " + req.GetBodyString())
		mylog.Info("[" + reqID + "]response: " + mylog.String(fdata))
	}
	if myfunc.BooleanDef(true, abort...) {
		req.Response.WriteJsonExit(fdata)
	} else {
		req.Response.WriteJson(fdata)
	}
}

func ReturnData(req *ghttp.Request, fstate bool, fmsg string, fdata interface{}, abort ...bool) {
	defaultResponse(req, rsp(fstate, fmsg, fdata), abort...)
}

func ReturnError(req *ghttp.Request, fmsg string, abort ...bool) {
	defaultResponse(req, rsp(false, fmsg, nil), abort...)
}

func ReturnOK(req *ghttp.Request, fmsg string, abort ...bool) {
	defaultResponse(req, rsp(true, fmsg, nil), abort...)
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
