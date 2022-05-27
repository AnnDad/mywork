package mywebapp

import (
	//"anndad/mylog"
	//"log"

	//"github.com/gogf/gf/encoding/gjson"
	"net/url"

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

func rsp(fstate bool, fmsg string, fdata interface{}) TRSP {
	return TRSP{State: fstate, Msg: fmsg, Data: fdata}
}

func ReturnData(req *ghttp.Request, fstate bool, fmsg string, fdata interface{}) {
	req.Response.WriteJsonExit(rsp(fstate, fmsg, fdata))
}

func ReturnError(req *ghttp.Request, fmsg string) {
	req.Response.WriteJsonExit(rsp(false, fmsg, nil))
}

func ReturnOK(req *ghttp.Request, fmsg string) {
	req.Response.WriteJsonExit(rsp(true, fmsg, nil))
}

func ReturnOK_Data(req *ghttp.Request, fdata interface{}) {
	req.Response.WriteJsonExit(rsp(true, "ok", fdata))
}

func ReturnCustomData(req *ghttp.Request, data interface{}) {
	req.Response.WriteJsonExit(data)
}

func ResultData(req *ghttp.Request, fstate bool, fmsg string, fdata interface{}) {
	req.Response.WriteJson(rsp(fstate, fmsg, fdata))
}

func ResultError(req *ghttp.Request, fmsg string) {
	req.Response.WriteJson(rsp(false, fmsg, nil))
}

func ResultOK(req *ghttp.Request, fmsg string) {
	req.Response.WriteJson(rsp(true, fmsg, nil))
}

func ResultOK_Data(req *ghttp.Request, fdata interface{}) {
	req.Response.WriteJson(rsp(true, "ok", fdata))
}

func CORSDefault(req *ghttp.Request) {
	req.Response.CORSDefault()
	req.Middleware.Next()
}
