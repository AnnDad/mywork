package myes

import (
	"errors"
	"strings"
	"time"

	"github.com/anndad/mywork/gftool"
	"github.com/anndad/mywork/myfunc"
	"github.com/anndad/mywork/mylog"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	_ "github.com/gogf/gf/os/gfile"
)

const (
	nodeHost  = "esHost"
	nodeIndex = "esIndex"
	nodeURL   = "esURL"
	levChar   = "/"
)

var (
	es_proxy   string = ""
	es_timeout int    = 60
	es_retry          = 9
	es_debug   bool   = false
)

type TESQueryExp struct {
	qeType   int //0--普通，1--range
	qeName   string
	qeName2  string
	qeValue  string
	qeValue2 string
}

type TESDataExp struct {
	qeName string
	qeData interface{}
}

func SetRetry(retry int) {
	es_retry = retry
	mylog.Info("ES重试次数：" + mylog.String(es_retry))
}

func SetProxy(proxy string) {
	es_proxy = proxy
	mylog.Info("ES启用代理：" + es_proxy)
}

func SetDebug(debug bool) {
	es_debug = debug
}
func SetTimeout(second int) {
	es_timeout = second
}

func intQueryExp() (result TESQueryExp) {
	result.qeType = 0
	result.qeName = ""
	result.qeName2 = ""
	result.qeValue = ""
	result.qeValue2 = ""
	return result
}

func intDataExp() (result TESDataExp) {
	result.qeName = ""
	result.qeData = nil
	return result
}

func createEmptyJson() *gjson.Json {
	result := gjson.New(nil)
	result.SetSplitChar(levChar[0])
	return result
}

func ESInit(esHost, esIndex string) *gjson.Json {
	result := createEmptyJson()
	result.Set(nodeHost, esHost)
	result.Set(nodeIndex, esIndex)
	result.Set(nodeURL, myfunc.EndWithSlash(esHost)+myfunc.EndWithSlash(esIndex))
	return result
}

func ReturnFields(json *gjson.Json, fields ...string) {
	json.Set("_source", fields)
}

func Exp_Match(field, value string) (result TESQueryExp) {
	result = intQueryExp()
	result.qeName = "match" + levChar + field
	result.qeValue = value
	return result
}
func Exp_Term(field, value string, iskeyword ...bool) (result TESQueryExp) {
	result = intQueryExp()
	result.qeName = "term" + levChar + field
	if len(iskeyword) > 0 {
		if iskeyword[0] {
			result.qeName = result.qeName + ".keyword"
		}
	}
	result.qeValue = value
	return result
}
func Exp_MatchP(field, value string) (result TESQueryExp) {
	result = intQueryExp()
	result.qeName = "match_phrase" + levChar + field
	result.qeValue = value
	return result
}
func Exp_Wildcard(field, value string) (result TESQueryExp) {
	result = intQueryExp()
	result.qeName = "wildcard" + levChar + field
	result.qeValue = value
	return result
}
func Exp_Regexp(field, value string) (result TESQueryExp) {
	result = intQueryExp()
	result.qeName = "regexp" + levChar + field
	result.qeValue = value
	return result
}
func Exp_Prefix(field, value string) (result TESQueryExp) {
	result = intQueryExp()
	result.qeName = "prefix" + levChar + field
	result.qeValue = value
	return result
}

func Exp_Range(field, value1, value2 string) (result TESQueryExp) {
	result = intQueryExp()
	result.qeType = 1
	result.qeName = "range" + levChar + field + levChar + "gte"
	result.qeName2 = "range" + levChar + field + levChar + "lte"
	result.qeValue = value1
	result.qeValue2 = value2
	return result
}

func Exp_RangeEx1(field, op, value string) (result TESQueryExp) {
	result = intQueryExp()
	result.qeType = 1
	result.qeName = "range" + levChar + field + levChar + opConvert(op)
	result.qeValue = value
	return result
}

func Exp_RangeEx2(field, op1, value1, op2, value2 string) (result TESQueryExp) {
	result = intQueryExp()
	result.qeType = 1
	result.qeName = "range" + levChar + field + levChar + opConvert(op1)
	result.qeName2 = "range" + levChar + field + levChar + opConvert(op2)
	result.qeValue = value1
	result.qeValue2 = value2
	return result
}

func Exp_Exists(field string) (result TESQueryExp) {
	result = intQueryExp()
	result.qeName = "exists" + levChar + "field"
	result.qeValue = field
	return result
}

func Exp_QString(value string) (result TESQueryExp) {
	result = intQueryExp()
	result.qeName = "query_string" + levChar + "query"
	result.qeValue = value
	return result
}

func logicEXP(prefixStr, errTip string, json *gjson.Json, exps ...TESQueryExp) {
	n := len(json.GetArray(prefixStr))
	for _, k := range exps {
		ItemName := prefixStr + levChar + gftool.IntToStr(n) + levChar
		switch k.qeType {
		case 0:
			json.Set(ItemName+k.qeName, k.qeValue)
		case 1:
			json.Set(ItemName+k.qeName, k.qeValue)
			if k.qeName2 != "" {
				json.Set(ItemName+k.qeName2, k.qeValue2)
			}
		default:
			panic("Unknow qeType: " + gftool.IntToStr(k.qeType))
		}

		n = n + 1
	}
	if n == 0 {
		panic(errTip + " cannot be empty!")
	}

}

func Exp_Sort(field, order string) (result TESQueryExp) {
	result.qeName = field + levChar + "order"
	result.qeValue = strings.ToLower(order)
	switch result.qeValue {
	case "asc":
	case "desc":
	default:
		panic("Unknow order: " + result.qeValue)
	}
	return result
}

func Exp_Data(field string, value interface{}) (result TESDataExp) {
	result = intDataExp()
	result.qeName = field
	result.qeData = value
	return result
}

func SetSort(json *gjson.Json, exps ...TESQueryExp) {
	prefixStr := "sort"
	n := len(json.GetArray(prefixStr))
	for _, k := range exps {
		ItemName := prefixStr + levChar + gftool.IntToStr(n) + levChar
		json.Set(ItemName+k.qeName, k.qeValue)
		n = n + 1
	}
}

func Query_AND(json *gjson.Json, exps ...TESQueryExp) {
	logicEXP("query"+levChar+"bool"+levChar+"must", "Query_AND", json, exps...)
}

func Query_OR(json *gjson.Json, exps ...TESQueryExp) {
	logicEXP("query"+levChar+"bool"+levChar+"should", "Query_OR", json, exps...)
}

func Query_NOT(json *gjson.Json, exps ...TESQueryExp) {
	logicEXP("query"+levChar+"bool"+levChar+"must_not", "Query_NOT", json, exps...)
}

func Filter_AND(json *gjson.Json, exps ...TESQueryExp) {
	logicEXP("filter"+levChar+"bool"+levChar+"must", "Filter_AND", json, exps...)
}

func Filter_OR(json *gjson.Json, exps ...TESQueryExp) {
	logicEXP("filter"+levChar+"bool"+levChar+"should", "Filter_OR", json, exps...)
}

func Filter_NOT(json *gjson.Json, exps ...TESQueryExp) {
	logicEXP("filter"+levChar+"bool"+levChar+"must_not", "Filter_NOT", json, exps...)
}

func DoMSearch(taskid string, jsons ...*gjson.Json) (*gjson.Json, error) {
	if len(jsons) == 0 {
		return nil, myfunc.NewError("Err_Empty_jsons")
	}
	url := jsons[0].GetString(nodeHost) + "/_msearch"
	postData := ""
	for _, json := range jsons {
		data := createEmptyJson()
		//指定返回的字段
		if json.Contains("_source") {
			data.Set("_source", json.Get("_source"))
		}
		//指定查询条件
		if json.Contains("query") {
			data.Set("query", json.Get("query"))
		}
		if json.Contains("filter") {
			data.Set("query"+levChar+"bool"+levChar+"filter", json.Get("filter"))
		}
		//指定排序
		if json.Contains("sort") {
			data.Set("sort", json.Get("sort"))
		}
		data.Set("track_total_hits", true)
		postData = postData + "{\"index\":\"" + json.GetString(nodeIndex) + "\"}" + "\n"
		postData = postData + data.MustToJsonString() + "\n"
	}

	r, err := DoPost(taskid, url, postData, true)
	return r, err
}

func DoQuery(taskid string, json *gjson.Json, page_index, page_size int) (*gjson.Json, error) {
	url := json.GetString(nodeURL) + "_search"
	data := createEmptyJson()
	//指定返回的字段
	if json.Contains("_source") {
		data.Set("_source", json.Get("_source"))
	}
	//指定查询条件
	if json.Contains("query") {
		data.Set("query", json.Get("query"))
	}
	if json.Contains("filter") {
		data.Set("query"+levChar+"bool"+levChar+"filter", json.Get("filter"))
	}
	//指定排序
	if json.Contains("sort") {
		data.Set("sort", json.Get("sort"))
	}
	//返回所有记录数
	data.Set("track_total_hits", true)
	PageIndex := page_index
	PageCount := 0
	if PageIndex >= 1 {
		RecordIndex := (PageIndex - 1) * page_size
		data.Set("from", RecordIndex)
	}
	data.Set("size", page_size)
	if es_debug {
		mylog.Info("DoQuery Debug: " + data.MustToJsonString())
		return nil, errors.New("ERR_InDebug")
	}
	r, err := DoPost(taskid, url, data.MustToJsonString(), true)
	if err == nil {
		RecordCount := r.GetInt("hits.total.value", 0)
		Rows := r.GetJsons("hits.hits")
		RowCount := len(Rows)
		if RecordCount == 0 {
			PageIndex = 0
			PageCount = 0
		} else {
			PageCount = myfunc.Div(RecordCount, page_size)
			n := myfunc.Mod(RecordCount, page_size)
			if n > 0 {
				PageCount = PageCount + 1
			}
			if PageIndex > PageCount {
				//PageIndex = PageCount
			}
		}
		r.Set("page_index", PageIndex)
		r.Set("page_size", page_size)
		r.Set("page_count", PageCount)
		r.Set("record_count", RecordCount)
		r.Set("row_count", RowCount)
	}
	return r, err
}

func DoQuery_ScrollData(taskid string, json *gjson.Json, page_size int) (*gjson.Json, error) {
	url := json.GetString(nodeURL) + "_search"
	data := createEmptyJson()
	//指定返回的字段
	if json.Contains("_source") {
		data.Set("_source", json.Get("_source"))
	}
	//指定查询条件
	if json.Contains("query") {
		data.Set("query", json.Get("query"))
	}
	if json.Contains("filter") {
		data.Set("query"+levChar+"bool"+levChar+"filter", json.Get("filter"))
	}
	//指定排序
	if json.Contains("sort") {
		data.Set("sort", json.Get("sort"))
	}
	//指定排序
	if json.Contains("search_after") {
		data.Set("search_after", json.Get("search_after"))
	}
	//Search after查询模式需要指定_id字段
	chksort := data.GetJsons("sort")
	n := len(chksort)
	//mylog.Println(n, chksort)
	if n > 0 {
		has_id := false
		for i := 0; i < n; i++ {
			if chksort[i].Contains("_id") {
				has_id = true
				break
			}
		}
		if has_id {
			n = -1
		}

	}
	if n >= 0 {
		//data.Set("sort"+levChar+gftool.IntToStr(n)+levChar+"_id"+levChar+"order", "desc")
	}
	//返回所有记录数
	data.Set("track_total_hits", true)
	data.Set("size", page_size)
	r, err := DoPost(taskid, url, data.MustToJsonString(), true)
	PageCount := 0
	if err == nil {
		RecordCount := r.GetInt("hits.total.value", 0)
		Rows := r.GetJsons("hits.hits")
		RowCount := len(Rows)
		if RecordCount == 0 {
			PageCount = 0
		} else {
			PageCount = myfunc.Div(RecordCount, page_size)
			n := myfunc.Mod(RecordCount, page_size)
			if n > 0 {
				PageCount = PageCount + 1
			}
		}
		if RowCount > 0 {
			r.Set("page_next", Rows[RowCount-1].Get("sort"))
		}
		r.Set("page_size", page_size)
		r.Set("page_count", PageCount)
		r.Set("record_count", RecordCount)
		r.Set("row_count", RowCount)
	}
	return r, err
}

func ScrollData_CanNext(json, last_data *gjson.Json) bool {
	if last_data.GetInt("row_count") < last_data.GetInt("page_size") {
		return false
	}
	if last_data.Contains("page_next") {
		json.Set("search_after", last_data.Get("page_next"))
		return true
	} else {
		return false
	}
}

func DoDelete(json *gjson.Json) (*gjson.Json, error) {
	url := json.GetString(nodeURL) + "_delete_by_query"
	data := createEmptyJson()
	FoundQuery := false
	if json.Contains("query") {
		data.Set("query", json.Get("query"))
		FoundQuery = true
	}
	if json.Contains("filter") {
		data.Set("query"+levChar+"bool"+levChar+"filter", json.Get("filter"))
		FoundQuery = true
	}
	if !FoundQuery {
		return nil, errors.New("ERR_Delete_Query_Empty")
	}
	if es_debug {
		mylog.Info("DoDelete Debug: " + data.MustToJsonString())
		return nil, errors.New("ERR_InDebug")
	}
	r, err := DoPost("", url, data.MustToJsonString(), false)
	return r, err
}

func DoAdd(json *gjson.Json, exps ...TESDataExp) (*gjson.Json, error) {
	url := json.GetString(nodeURL) + "_doc"
	for _, k := range exps {
		json.Set("data"+levChar+k.qeName, k.qeData)
	}
	if !json.Contains("data") {
		return nil, errors.New("ERR_Add_NoData")
	}
	data := gjson.New(json.Get("data"))

	if es_debug {
		mylog.Info("DoAdd Debug: " + data.MustToJsonString())
		return nil, errors.New("ERR_InDebug")
	}

	r, err := DoPost("", url, data.MustToJsonString(), false)
	return r, err
}

func DoUpdateByID_all(taskid string, json *gjson.Json, id string, exps ...TESDataExp) (*gjson.Json, error) {
	url := json.GetString(nodeURL) + "_doc/" + id
	for _, k := range exps {
		json.Set("data"+levChar+k.qeName, k.qeData)
	}
	if !json.Contains("data") {
		return nil, errors.New("ERR_DoUpdateAsPUT_NoData")
	}
	data := gjson.New(json.Get("data"))

	if es_debug {
		mylog.Info("DoUpdateAsPUT Debug: " + data.MustToJsonString())
		return nil, errors.New("ERR_InDebug")
	}

	r, err := DoPost(taskid, url, data.MustToJsonString(), false)
	return r, err
}
func DoUpdateByID_part_with_add(taskid string, json *gjson.Json, id string, exps ...TESDataExp) (*gjson.Json, error) {
	rsp, err := DoUpdateByID_part(taskid, json, id, exps...)
	if err != nil {
		return nil, err
	}
	if !UpdateDone(rsp) {
		rsp, err = DoUpdateByID_all(taskid, json, id, exps...)
	}
	return rsp, err
}
func UpdateDone(rsp *gjson.Json) bool {
	return rsp.GetInt("items.0.update.status") == 200
}

func IsEmpty(rsp *gjson.Json) bool {
	return rsp.GetInt("record_count") == 0
}

func DoUpdateByID_part(taskid string, json *gjson.Json, id string, exps ...TESDataExp) (*gjson.Json, error) {
	if id == "" {
		return nil, errors.New("ERR_ID_Is_Empty")
	}
	url := json.GetString(nodeHost) + "/_bulk"
	actJS := createEmptyJson()
	actJS.Set("update"+levChar+"_index", json.GetString(nodeIndex))
	actJS.Set("update"+levChar+"_id", id)
	rowJS := createEmptyJson()
	for _, k := range exps {
		rowJS.Set("doc"+levChar+k.qeName, k.qeData)
	}
	data := actJS.MustToJsonString() + "\n" + rowJS.MustToJsonString() + "\n"
	if es_debug {
		mylog.Info("DoUpdateByID Debug: " + data)
		return nil, errors.New("ERR_InDebug")
	}
	r, err := DoPost(taskid, url, data, false)
	return r, err
}

func DoUpdateScript(taskid string, json *gjson.Json, exps ...TESDataExp) (*gjson.Json, error) {
	fields := ""
	for _, k := range exps {
		fields = fields + "ctx._source." + k.qeName + "=" + gftool.Any2Str(k.qeData) + ";"
	}
	json.Set("script", fields)
	url := json.GetString(nodeURL) + "_update_by_query"
	if !json.Contains("script") {
		return nil, errors.New("ERR_Update_NoData")
	}
	data := createEmptyJson()
	FoundQuery := false
	if json.Contains("query") {
		data.Set("query", json.Get("query"))
		FoundQuery = true
	}
	if json.Contains("filter") {
		data.Set("query"+levChar+"bool"+levChar+"filter", json.Get("filter"))
		FoundQuery = true
	}
	if !FoundQuery {
		return nil, errors.New("ERR_Update_Query_Empty")
	}
	data.Set("script", json.Get("script"))
	if es_debug {
		mylog.Info("DoUpdate Debug: " + data.MustToJsonString())
		return nil, errors.New("ERR_InDebug")
	}
	r, err := DoPost(taskid, url, data.MustToJsonString(), false)
	return r, err
}

func DoPost(taskid, url string, data string, saveasDetail bool) (*gjson.Json, error) {
	c := g.Client()
	if es_proxy != "" {
		c.SetProxy(es_proxy)
	}
	c.SetTimeout(time.Second * time.Duration(es_timeout))

	ftaskid := taskid
	if ftaskid == "" {
		ftaskid = mylog.GetGuid()
	}
	mylog.Info("[" + ftaskid + "]提交ES[" + url + "]: " + data)
	c.SetContentType("application/json")
	retry := 0
	retryMax := es_retry
	for {
		retry = retry + 1
		tip := "[" + ftaskid + "第" + mylog.String(retry) + "次]提交ES"
		r, err := c.Post(url, data)
		if err != nil {
			mylog.Info(tip + "出错: " + err.Error())
			if retryMax < 0 {
				myfunc.SleepSeconds(1)
				continue
			} else {
				if retry >= retryMax {
					return nil, err
				} else {
					myfunc.SleepSeconds(1)
					continue
				}
			}

		} else {
			defer r.Close()
			html := r.ReadAllString()
			tip = tip + "结果: "
			if saveasDetail {
				mylog.Info(tip, data+"\n\n"+html)
			} else {
				mylog.Info(tip + html)
			}
			result := gjson.New(html)
			return result, nil
		}
	}

}

func JsonToExps(js *gjson.Json, excludes ...string) (result []TESDataExp) {
	for key, value := range js.ToMap() {
		if myfunc.StrInArray(key, excludes) {
			continue
		}
		result = append(result, Exp_Data(key, value))
	}
	return result
}

func opConvert(op string) string {
	switch op {
	case ">":
		return "gt"
	case ">=":
		return "gte"
	case "<":
		return "lt"
	case "<=":
		return "lte"
	default:
		panic("Unknow OP: " + op)
	}
}
