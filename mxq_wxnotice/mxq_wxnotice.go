package mxq_wxnotice

import (
	"github.com/anndad/mywork/gftool"
	"github.com/anndad/mywork/mylog"

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

func SendWXNotice(msgGroup, msg string, At []string, fceSend ...bool) {
	taskid := mylog.GetGuid()
	checkLastSendTime := true
	if len(fceSend) > 0 {
		checkLastSendTime = !fceSend[0]
	}
	data := gjson.New(nil)
	data.Set("checkLastSendTime", checkLastSendTime)
	data.Set("msg.msgtype", "text")
	data.Set("msg.text.content", msg)
	if len(At) > 0 {
		data.Set("msg.text.mentioned_mobile_list", At)
	}
	SendToWeChat(taskid, msgGroup, data)
}

func SendToWeChat(taskid, fgroup string, msg *gjson.Json) bool {
	mylog.Info("[" + taskid + "]微信消息: " + msg.MustToJsonString())
	if config.GetInt("wechat.actived", 0) != 1 {
		mylog.Info("[" + taskid + "]未发送到微信: wechat.actived=0")
		return false
	}
	hook := getHook()
	if hook == "" {
		mylog.Info("[" + taskid + "]未发送到微信: 配置文件中没有找到节点[wechat.webhooks]")
		return false
	}

	if msg.GetBool("checkLastSendTime", true) {
		if fgroup == "" {
			fgroup = "default"
		}
		SendInterval := config.GetInt64("wechat.sendinterval_sec")
		if SendInterval > 0 {
			LastSentStr := config.GetString("wechat." + fgroup + ".lastsent")
			if LastSentStr != "" {
				t := gftool.SecondsBetweenNow(LastSentStr)
				if t < SendInterval {
					mylog.Info("[" + taskid + "]未发送到微信: 距离上次发送间隔太短(" + mylog.String(t) + "秒)")
					return false
				}
			}
		}
	}

	//todo: 区分类型,检查消息长度
	// if len(msg) >= 2000 {
	// 	mylog.Info("[" + taskid + "]文本被截断: " + msg)
	// 	data.Set("text.content", msg[:2000])
	// } else {
	// 	data.Set("text.content", msg)
	// }

	mylog.Info("[" + taskid + "]提交地址: " + hook)
	r, err := g.Client().Post(hook, msg.GetJson("msg").MustToJsonString())
	if err != nil {
		mylog.Error("["+taskid+"]提交出错: ", err)
		return false
	} else {
		defer r.Close()
		if fgroup != "" {
			config.Set("wechat."+fgroup+".lastsent", gftool.Now2Str())
		}
		html := r.ReadAllString()
		mylog.Info("[" + taskid + "]提交结果: " + html)
		return true
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
