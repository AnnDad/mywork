package myredis

import (
	"github.com/anndad/mywork/myfunc"
	"github.com/anndad/mywork/mylog"
	"github.com/gogf/gf/frame/g"
)

var (
	MaxRetry int = 3
)

func Access(taskid, redis string, commandName string, args ...interface{}) (string, error) {
	var err error
	tip := commandName

	if len(args) > 0 {
		tip = tip + " " + myfunc.String(args[0])
	}

	if len(args) > 1 {
		value := myfunc.String(args[1])
		tip = tip + " " + value
	}
	tip = "[" + taskid + "][" + redis + "]RedisAccess[" + tip + "], "
	connRedis := g.Redis(redis)
	result := ""
	for n := 1; n <= MaxRetry; n++ {
		data, err := connRedis.DoVarWithTimeout(myfunc.Seconds(30), commandName, args...)
		if err != nil {
			mylog.Error(tip+"Error(try"+myfunc.String(n)+"): ", err)
			continue
		} else {
			mylog.Info(tip + "Successed(try" + myfunc.String(n) + ")")
			if data != nil {
				result = data.String()
			}
			return result, nil
		}
	}
	return "", err
}
