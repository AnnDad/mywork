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
type TOnDone func(data string, err error) bool

func Execute(Func_Exec TOnExecute, Func_Done TOnDone, taskData ...string) (*myjson.TJson, error) {
	chDone := make(chan string)
	results := myjson.New()
	TaskCount := len(taskData)
	for i, item := range taskData {
		go execInThread("thread"+myfunc.String(i), item, Func_Exec, chDone)
	}
	//wait
	for i := 0; i < this.TaskCount; i++ {
		result := <-this.chDone
		json := myjson.New(result)
		if Func_Done != nil {
			var err error
			err_msg := json.GetString("task.err")
			data := json.GetString("task.data")
			if err_msg != "ok" {
				err = myfunc.NewError(err_msg)
			}
			is_abort := Func_Done(data, err)
			if !is_abort {
				break
			}
		}
		results.Append("results", json.Get("task"))
	}
	return results
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
	result.Set("task.err", msg)
}
