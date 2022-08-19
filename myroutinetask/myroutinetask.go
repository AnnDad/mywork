package myroutinetask

import (
	"github.com/anndad/mywork/myfunc"

	"github.com/anndad/mywork/myjson"
)

type TRoutineTask struct {
	TaskCount  int
	chDone     chan string
	chTaskList chan string
	onExecute  TOnExecute
	results    *myjson.TJson
}

//线程执行函数
//输入: taskData 任务参数
//输出: string 执行结果, error 执行过程中是否有错误
type TOnExecute func(taskData string) (string, error)

//某个任务执行完成后
//输入: data 执行结果, err 执行过程中是否有错误
//输出: bool 是否继续等待完成, true->继续等待,false->不等待
type TOnDone func(data string, err error) bool

func NewRoutineTask() *TRoutineTask {
	this := new(TRoutineTask)
	this.chDone = make(chan string)
	this.results = myjson.New()
	return this
}

func (this *TRoutineTask) AddTask(taskData ...string) {
	this.TaskCount = len(taskData)
	this.chTaskList = make(chan string, this.TaskCount)
	for _, item := range taskData {
		this.chTaskList <- item
	}
}

func (this *TRoutineTask) Execute(onExecute TOnExecute) {
	this.onExecute = onExecute
	for i := 0; i < this.TaskCount; i++ {
		taskData := <-this.chTaskList
		go this.ExecuteInThread("thread"+myfunc.String(i), taskData)
	}
}

func (this *TRoutineTask) Wait(onDone ...TOnDone) {
	for i := 0; i < this.TaskCount; i++ {
		result := <-this.chDone
		json := myjson.New(result)
		if len(onDone) > 0 {
			var err error
			err_msg := json.GetString("task.err")
			data := json.GetString("task.data")
			if err_msg != "" {
				err = myfunc.NewError(err_msg)
			}
			is_abort := onDone[0](data, err)
			//myfunc.Println("isContinue: ", isContinue)
			if is_abort {
				//myfunc.Println("abort")
				break
			}
		}
		this.results.Append("results", json.Get("task"))
	}
}

func (this *TRoutineTask) ExecuteInThread(threadid, taskData string) {
	result := myjson.New()
	result.Set("task.source", taskData)
	defer func() {
		this.chDone <- result.MustToJsonString()
	}()
	data, err := this.onExecute(taskData)
	result.Set("task.data", data)
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	result.Set("task.err", msg)
}
