package mydbsql

import (
	"database/sql"

	"errors"
	"sort"

	"github.com/anndad/mywork/gftool"
	"github.com/anndad/mywork/myfunc"
	"github.com/anndad/mywork/mylog"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
)

type paramitem struct {
	ParamName string
	FieldName string
	Length    int
}

type TransItem struct {
	sql   string
	param g.Map
}

type TMap = g.Map

type TFieldValue = [2]string

const Null_DB = "@@null_db@@"

func EmptyParams() g.Map {
	return g.Map{}
}

func InitParams(fieldName string, fieldValue interface{}) g.Map {
	return g.Map{fieldName: fieldValue}
}

func DB(name ...string) gdb.DB {
	return g.DB(name...)
}
func TI(sql string, param g.Map) TransItem {
	var result TransItem
	result.sql = sql
	result.param = param
	return result
}
func FV(fieldName, fieldValue string) TFieldValue {
	var result TFieldValue
	result[0] = fieldName
	result[1] = fieldValue
	return result
}

func GBKOrder_mysql(field string, isDesc ...bool) string {
	return " convert(" + field + " using gbk) " + myfunc.CaseWhenStr(myfunc.BooleanDef(false, isDesc...), "desc", "") + " "
}

func PrepareSQL(sql string, params g.Map) (string, []interface{}, string, error) {
	r1 := sql
	r2 := make([]interface{}, 0)
	var paramlist []paramitem
	match, _ := gregex.MatchAllString(`\:(\w*)`, sql)
	for _, value := range match {
		param_name := value[0]
		field_name := value[1]
		value, ok := params[field_name]
		if !ok {
			mylog.Info("PrepareSQL Error: \n" + sql)
			return "", nil, "", errors.New("params not exist: [" + field_name + "]")
		}
		paramlist = append(paramlist, paramitem{param_name, field_name, len(param_name)})
		if myfunc.String(value) != Null_DB {
			r2 = append(r2, value)
		}
	}
	//按照参数名称的长度倒序排列, 保证名称长的先处理
	sort.Slice(paramlist, func(i, j int) bool {
		return paramlist[i].Length > paramlist[j].Length
	})
	sqlLog := r1
	for _, item := range paramlist {
		value := myfunc.String(params[item.FieldName])
		if value == Null_DB {
			r1 = gstr.Replace(r1, item.ParamName, "null")
			sqlLog = gstr.Replace(sqlLog, item.ParamName, "null")
		} else {
			r1 = gstr.Replace(r1, item.ParamName, "?")
			sqlLog = gstr.Replace(sqlLog, item.ParamName, value)
		}
	}

	return r1, r2, sqlLog, nil
}

func SetRowField(rows *gjson.Json, index int, field string, value interface{}) {
	rows.Set("rows."+myfunc.String(index)+"."+field, value)
}
func GetRows(rows *gjson.Json) *gjson.Json {
	if rows.Contains("rows") {
		return rows.GetJson("rows")
	}
	return gjson.New("[]")
}

func MapFieldAsString(obj g.Map, field string) string {
	v, ok := obj[field]
	if ok {
		return myfunc.String(v)
	} else {
		return ""
	}
}

func ExecSQL(taskid string, DB gdb.DB, SQL string, params g.Map) (result sql.Result, err error) {
	r1, r2, rLog, err := PrepareSQL(SQL, params)
	if err != nil {
		return nil, err
	}
	result, err = DB.Exec(r1, r2)
	logExecSQLResult(taskid, rLog, result, err)
	return result, err
}

func transExecSQL(taskid string, DB *gdb.TX, SQL string, params g.Map) (result sql.Result, err error) {
	r1, r2, rLog, err := PrepareSQL(SQL, params)
	if err != nil {
		return nil, err
	}

	result, err = DB.Exec(r1, r2)
	logExecSQLResult(taskid, rLog, result, err)
	return result, err
}

func logExecSQLResult(taskid, rLog string, result sql.Result, err error) {
	log := "execSQL: " + rLog + "\n" + "execResult: "
	if err == nil {
		n, err2 := result.RowsAffected()
		log = log + "Successed, "
		if err2 == nil {
			log = log + "RowsAffected: " + gftool.Int64ToStr(n)
		} else {
			log = log + "RowsAffected Error: " + err2.Error()
		}
	} else {
		log = log + "Failed, Error: " + err.Error()
	}
	mylog.Info("[" + taskid + "]" + log)
}

func ExecSQLTrans(taskid string, DB gdb.DB, SQLs ...TransItem) (err error) {
	trans, err := DB.Begin()
	if err != nil {
		return errors.New("[" + taskid + "]err_trans_begin: " + err.Error())
	}
	mylog.Info("[" + taskid + "] BeginTrans;")
	defer func() {
		if err != nil {
			mylog.Info("[" + taskid + "] Rollback;")
			trans.Rollback()
		} else {
			mylog.Info("[" + taskid + "] Commit;")
			trans.Commit()
		}
	}()
	for i, _ := range SQLs {
		_, err = transExecSQL(taskid, trans, SQLs[i].sql, SQLs[i].param)
		if err != nil {
			return errors.New("err_trans_exec: " + err.Error())
		}
	}
	return nil
}

func GetDataPage(taskid string, DB gdb.DB, SQL string, params g.Map, pIndex, pSize int) (result *gjson.Json, err error) {
	pSQL := "select count(*) rcount from (" + SQL + ") vtbCount"
	pRs, err := GetData(taskid, DB, pSQL, params)
	if err != nil {
		return nil, err
	}
	if pIndex < 1 {
		pIndex = 1
	}
	recordcount := pRs.GetInt("rows.0.rcount")
	pSQL = "select * from (" + SQL + ") vtbdata "
	if pSize > 0 && pSize <= 100 {
		pSQL = pSQL + " limit " + myfunc.String((pIndex-1)*pSize) + "," + myfunc.String(pSize)
	}
	data, err := GetData(taskid, DB, pSQL, params)
	if err != nil {
		return nil, err
	}
	data.Set("recordcount", recordcount)
	data.Set("pIndex", pIndex)
	data.Set("pSize", pSize)
	return data, nil
}

func IsEmpty(rs *gjson.Json) bool {
	return RecordCount(rs) == 0
}

func RecordCount(rs *gjson.Json) int {
	return rs.GetInt("recordcount")
}

func RowCount(rs *gjson.Json) int {
	return rs.GetInt("rowcount")
}

func GetData(taskid string, DB gdb.DB, SQL string, params g.Map) (result *gjson.Json, err error) {
	r1, r2, rLog, err := PrepareSQL(SQL, params)
	if err != nil {
		return nil, err
	}
	mylog.Info("[" + taskid + "]dataSQL: " + rLog)
	rows, err := DB.Query(r1, r2)
	if err != nil {
		return nil, err
	} else {
		defer rows.Close()
		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		count := len(columns)
		tableData := make([]map[string]interface{}, 0)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)
		for rows.Next() {
			for i := 0; i < count; i++ {
				valuePtrs[i] = &values[i]
			}
			rows.Scan(valuePtrs...)
			entry := make(map[string]interface{})
			for i, col := range columns {
				var v interface{}
				val := values[i]
				b, ok := val.([]byte)
				if ok {
					v = string(b)
				} else {
					v = val
				}
				entry[col] = v
			}
			tableData = append(tableData, entry)
		}
		result = gjson.New(tableData)
		rowcount := len(tableData)
		//recordcount=记录总数, rowcount=分页记录数
		result.Set("recordcount", rowcount)
		result.Set("rowcount", rowcount)
		result.Set("fieldcount", len(columns))
		result.Set("fields", columns)
		result.Set("rows", tableData)
		result.Set("isempty", rowcount == 0)
		return result, nil
	}
}

func CaseWhen(exp bool, value1, value2 interface{}) (result interface{}) {

	if exp {
		result = value1
	} else {
		result = value2
	}
	return result
}

func GetResultRows_bak(rows *sql.Rows) map[int]map[string]interface{} {

	//返回所有列
	columns, _ := rows.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(columns))
	//这里表示一行填充数据
	scans := make([]interface{}, len(columns))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	result := make(map[int]map[string]interface{})
	for rows.Next() {
		//填充数据
		rows.Scan(scans...)
		//每行数据
		row := make(map[string]interface{})
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := columns[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result[i] = row
		i++
	}
	return result
}

func GetIDx_mysql(keyname string, cnt int) (result []int64) {
	sql := "CALL getidx('" + keyname + "'," + myfunc.String(cnt) + ");"
	rs, _ := GetData("", DB(), sql, EmptyParams())
	maxid := rs.GetInt64("rows.0.result")
	n := maxid - gftool.Str2Int64(gftool.IntToStr(cnt)) + 1
	for i := n; i <= maxid; i++ {
		result = append(result, i)
	}
	return result
}
