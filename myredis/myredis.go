package myredis

import (
	"anndad/mywxnotice"

	"github.com/anndad/mywork/myfunc"
	"github.com/anndad/mywork/mylog"

	"github.com/gogf/gf/frame/g"
)

func RedisAccess(redis string, commandName string, args ...interface{}) (string, error) {
	tip := commandName
	if len(args) > 0 {
		tip = tip + " " + myfunc.String(args[0])
	}
	mylog.Info("RedisAccess[" + tip + "]...")
	connRedis := g.Redis(redis)
	result := ""
	for n := 1; n <= 50; n++ {
		data, err := connRedis.DoVarWithTimeout(myfunc.Seconds(30), commandName, args...)
		if err != nil {
			mylog.Info("RedisAccess[" + tip + "] Error(try" + myfunc.String(n) + "): " + err.Error())
			if n >= 50 {
				mywxnotice.Send("", "redis访问失败超过50次: "+err.Error(), true)
			}
			continue
		} else {
			mylog.Info("RedisAccess[" + tip + "] Successed(try" + myfunc.String(n) + ")")
			if data != nil {
				result = data.String()
			}
			return result, nil
		}
	}
	return "", myfunc.NewError("redis访问失败超过50次")
}
