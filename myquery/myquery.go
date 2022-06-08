package myquery

import (
	"strings"

	"github.com/anndad/mywork/myfunc"

	"github.com/PuerkitoBio/goquery"
	"github.com/gogf/gf/text/gregex"
)

type TSelection = goquery.Selection

func Html2Doc(html string) (*TSelection, error) {
	doc, err := goquery.NewDocumentFromReader(myfunc.Str2Reader(html))
	if err != nil {
		return nil, err
	}
	return doc.Find("html"), nil
}
func GetOuterHTML(obj *TSelection) string {
	html, _ := goquery.OuterHtml(obj)
	return html
}

func ExecuteExp(source, exp string) string {
	items, _ := gregex.MatchAllString(exp, source)
	result := ""
	//myfunc.Println("source: ", source)
	//myfunc.Println("exp: ", exp)
	//myfunc.Println("items: ", items, " >>>>")
	if len(items) > 0 {
		if len(items[0]) > 0 {
			result = items[0][1]
		}
	}
	result = strings.TrimLeft(result, " ")
	result = strings.TrimRight(result, " ")
	return result
}

// 获取当前HTML标签的某个属性
func GetHTMLAttrib(html, attribEXP, attribADD string) (string, string) {
	result1 := ExecuteExp(html, attribEXP)
	result2 := result1
	if attribADD != "" {
		result2 = myfunc.Replace(attribADD, "@@this@@", result2)
	}
	return result2, result1
}

//获取当前对象的属性
func GetHTMLAttrib_Obj(obj *goquery.Selection, attribEXP, attribADD string) (string, string) {
	html, _ := goquery.OuterHtml(obj)
	result1, result2 := GetHTMLAttrib(html, attribEXP, attribADD)
	return result1, result2
}

//查找当前对象的子对象并获取该对象的属性
func GetHTMLAttrib_SubObj(obj *goquery.Selection, attribSTR, attribEXP, attribADD string) (string, string) {
	result1 := ""
	result2 := ""
	subObj := obj.Find(attribSTR)
	if subObj != nil {
		result1, result2 = GetHTMLAttrib_Obj(subObj, attribEXP, attribADD)
	}
	return result1, result2
}
