package mywebapp

import (
	//"anndad/mylog"
	//"log"

	"net/url"

	"github.com/anndad/mywork/myjson"

	"github.com/anndad/mywork/myfunc"
	"github.com/gogf/gf/net/ghttp"
)

type TRSP struct {
	State bool        `json:"state"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

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

func ReturnData(req *ghttp.Request, fstate bool, fmsg string, fdata interface{}, abort ...bool) {
	if myfunc.BooleanDef(true, abort...) {
		req.Response.WriteJsonExit(rsp(fstate, fmsg, fdata))
	} else {
		req.Response.WriteJson(rsp(fstate, fmsg, fdata))
	}
}

func ReturnError(req *ghttp.Request, fmsg string, abort ...bool) {
	if myfunc.BooleanDef(true, abort...) {
		req.Response.WriteJsonExit(rsp(false, fmsg, nil))
	} else {
		req.Response.WriteJson(rsp(false, fmsg, nil))
	}
}

func ReturnOK(req *ghttp.Request, fmsg string, abort ...bool) {
	if myfunc.BooleanDef(true, abort...) {
		req.Response.WriteJsonExit(rsp(true, fmsg, nil))
	} else {
		req.Response.WriteJson(rsp(true, fmsg, nil))
	}
}

func ReturnOK_Data(req *ghttp.Request, fdata interface{}, abort ...bool) {
	if myfunc.BooleanDef(true, abort...) {
		req.Response.WriteJsonExit(rsp(true, "ok", fdata))
	} else {
		req.Response.WriteJson(rsp(true, "ok", fdata))
	}
}

func ReturnCustomData(req *ghttp.Request, data interface{}, abort ...bool) {
	if myfunc.BooleanDef(true, abort...) {
		req.Response.WriteJsonExit(data)
	} else {
		req.Response.WriteJson(data)
	}
}

func CORSDefault(req *ghttp.Request) {
	req.Response.CORSDefault()
	req.Middleware.Next()
}
