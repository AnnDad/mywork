package mytask

import (
	"github.com/anndad/mywork/myfunc"

	"github.com/anndad/mywork/myjson"
)

//线程执行函数
//输入: taskData 任务参数
//输出: string 执行结果, error 执行过程中是否有错误
type TOnExecute func(taskData string) (string, error)

//某个任务执行完成后
//输入: data 执行结果, err 执行过程中是否有错误
//输出: bool 是否继续等待完成, true->继续等待,false->不等待

type TOnDone func(allDone *myjson.TJson, data string, err error) bool

func Execute(Func_Exec TOnExecute, taskData ...string) *myjson.TJson {
	TaskCount := len(taskData)
	chDone := make(chan string, TaskCount)
	allDone := myjson.New()
	for i, item := range taskData {
		go execInThread("thread"+myfunc.String(i), item, Func_Exec, chDone)
	}
	//wait
	for i := 0; i < TaskCount; i++ {
		result := <-chDone
		jsResult := myjson.New(result)
		err_msg := jsResult.GetString("task.msg")
		if err_msg == "ok" {
			allDone.Append("done_ok", jsResult.Get("task"))
		} else {
			allDone.Append("done_err", jsResult.Get("task"))
		}
		allDone.Append("tasks", jsResult.Get("task"))
	}
	return allDone
}

func ExecuteWithDone(Func_Exec TOnExecute, Func_Done TOnDone, taskData ...string) *myjson.TJson {
	TaskCount := len(taskData)
	chDone := make(chan string, TaskCount)
	allDone := myjson.New()
	for i, item := range taskData {
		go execInThread("thread"+myfunc.String(i), item, Func_Exec, chDone)
	}
	//wait
	for i := 0; i < TaskCount; i++ {
		result := <-chDone
		jsResult := myjson.New(result)
		var err error
		err_msg := jsResult.GetString("task.msg")
		data := jsResult.GetString("task.data")
		if err_msg == "ok" {
			allDone.Append("done_ok", jsResult.Get("task"))
		} else {
			err = myfunc.NewError(err_msg)
			allDone.Append("done_err", jsResult.Get("task"))
		}
		allDone.Append("tasks", jsResult.Get("task"))
		is_Continue := Func_Done(allDone, data, err)
		if !is_Continue {
			break
		}

	}
	return allDone
}

func execInThread(threadid, taskData string, Func_Exec TOnExecute, chDone chan string) {
	result := myjson.New()
	result.Set("task.source", taskData)
	defer func() {
		chDone <- result.MustToJsonString()
	}()
	data, err := Func_Exec(taskData)
	result.Set("task.data", data)
	msg := "ok"
	if err != nil {
		msg = "err: " + err.Error()
	}
	result.Set("task.msg", msg)
}
