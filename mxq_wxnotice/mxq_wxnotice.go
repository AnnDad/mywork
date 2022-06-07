package mxq_wxnotice

import (
	"github.com/anndad/mywork/gftool"
	"github.com/anndad/mywork/mylog"
	"github.com/anndad/mywork/mytool"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcfg"
)

var (
	config *gcfg.Config
)

func init() {
	config = gcfg.New("mxq_wxnotice.cfg.toml")
}

func SetConfigPath(path string) {

	config.SetPath(path)
}

func SendWXNotice(msgGroup, msg string, fceSend ...bool) {
	taskid := mytool.GetGuid()
	mylog.Info("[" + taskid + "]微信消息: " + msg)
	if config.GetInt("wechat.actived", 0) != 1 {
		mylog.Info("wechat.actived=0,未发送到微信")
		return
	}
	fGroup := msgGroup
	if fGroup == "" {
		fGroup = "default"
	}
	checkLastSendTime := true
	if len(fceSend) > 0 {
		checkLastSendTime = !fceSend[0]
	}
	if checkLastSendTime {
		SendInterval := config.GetInt64("wechat.sendinterval_sec")
		if SendInterval != 0 {
			LastSentStr := config.GetString("wechat." + fGroup + ".lastsent")
			if LastSentStr != "" {
				t := gftool.SecondsBetweenNow(LastSentStr)
				if t < SendInterval {
					mylog.Info("距离上次发送间隔" + gftool.Int64ToStr(t) + "秒，跳过...")
					return
				}
			}
		}
	}

	data := gjson.New(nil)
	data.Set("msgtype", "text")
	if len(msg) >= 2000 {
		mylog.Info("[" + taskid + "]文本被截断: " + msg)
		data.Set("text.content", msg[:2000])
	} else {
		data.Set("text.content", msg)
	}
	hook := getHook()
	if hook == "" {
		mylog.Info("配置文件中没有找到节点[wechat.webhooks]")
		return
	}
	mylog.Info("[" + taskid + "]提交地址: " + hook)
	//mylog.Info("[" + taskid + "]提交数据: " + data.MustToJsonString())
	r, err := g.Client().Post(hook, data.MustToJsonString())
	if err != nil {
		mylog.Info("[" + taskid + "]提交出错: " + err.Error())
		return
	} else {
		config.Set("wechat."+fGroup+".lastsent", gftool.Now2Str())
		defer r.Close()
		html := r.ReadAllString()
		mylog.Info("[" + taskid + "]提交结果: " + html)
	}
}

func getHook() string {
	hooks := config.GetArray("wechat.webhooks")
	n := len(hooks)
	if n == 0 {
		return ""
	}
	index := config.GetInt("wechat.lasthook", -1) + 1
	if index >= n {
		index = 0
	}
	config.Set("wechat.lasthook", index)
	return hooks[index].(string)
}
