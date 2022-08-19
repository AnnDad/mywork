package myhtml

import (
	"io/ioutil"
	"net"
	"net/http"
	_ "net/http/httputil"
	"regexp"
	"strings"
	"time"

	"github.com/anndad/mywork/myfunc"
	"github.com/anndad/mywork/mylog"

	"encoding/base64"

	"github.com/axgle/mahonia"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"golang.org/x/net/html/charset"
)

var (
	client_timeout int = 10
	client_retry   int = 0
)

func SetDefaultHttpClientTimeout(timeout int) {
	client_timeout = timeout
}

//DialWithIP 根据指定Ip连接；
func DialWithIP(netw, addr, outip string, timeout int) (net.Conn, error) {
	TimeOut := time.Duration(timeout)
	//本地地址  ipaddr是本地外网IP
	lAddr, err := net.ResolveTCPAddr(netw, outip+":0")
	if err != nil {
		return nil, err
	}
	//被请求的地址
	rAddr, err := net.ResolveTCPAddr(netw, addr)
	if err != nil {
		return nil, err
	}
	connTimeOut, errTimeOut := net.DialTimeout(netw, addr, TimeOut*time.Second)
	if errTimeOut != nil {
		return nil, errTimeOut
	}
	defer connTimeOut.Close()

	conn, err := net.DialTCP(netw, lAddr, rAddr)
	if err != nil {
		return nil, err
	}
	deadline := time.Now().Add(TimeOut * time.Second)
	conn.SetDeadline(deadline)
	return conn, nil
}

func HttpGetByIP(url, outip string, timeout int) (string, error) {
	result := ""
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				myfunc.PrintlnTip("netw: ", netw)
				return DialWithIP(netw, addr, outip, timeout)
			},
		},
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_3) AppleWebKit/537.36 (KHTML,like Gecko) Chrome/27.0.1453.93 Safari/537.36")
	resp, err := client.Do(req)
	if err == nil {
		if resp.Body != nil {
			defer resp.Body.Close()
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			bodystr := myfunc.UnsafeBytesToStr(body)
			result = ConvertToUTF8(bodystr, "")
		}
	}

	return result, err
}

func ConvertToUTF8(content string, contentType string) string {
	var htmlEncode string

	if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
		htmlEncode = "gb18030"
	} else if strings.Contains(contentType, "big5") {
		htmlEncode = "big5"
	} else if strings.Contains(contentType, "utf-8") {
		htmlEncode = "utf-8"
	}
	if htmlEncode == "" {
		//先尝试读取charset
		reg := regexp.MustCompile(`(?is)<meta[^>]*charset\s*=["']?\s*([A-Za-z0-9\-]+)`)
		match := reg.FindStringSubmatch(content)
		if len(match) > 1 {
			contentType = strings.ToLower(match[1])
			if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
				htmlEncode = "gb18030"
			} else if strings.Contains(contentType, "big5") {
				htmlEncode = "big5"
			} else if strings.Contains(contentType, "utf-8") {
				htmlEncode = "utf-8"
			}
		}
		if htmlEncode == "" {
			reg = regexp.MustCompile(`(?is)<title[^>]*>(.*?)<\/title>`)
			match = reg.FindStringSubmatch(content)
			if len(match) > 1 {
				aa := match[1]
				_, contentType, _ = charset.DetermineEncoding([]byte(aa), "")

				if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
					htmlEncode = "gb18030"
				} else if strings.Contains(contentType, "big5") {
					htmlEncode = "big5"
				} else if strings.Contains(contentType, "utf-8") {
					htmlEncode = "utf-8"
				}
			}
		}
	}
	htmlEncode = strings.ToLower(htmlEncode)
	if htmlEncode != "" && htmlEncode != "utf-8" {
		content = Convert(content, htmlEncode, "utf-8")
	}

	return content
}

/**
 * 编码转换
 * 需要传入原始编码和输出编码，如果原始编码传入出错，则转换出来的文本会乱码
 */
func Convert(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func HttpClient() *ghttp.Client {
	c := g.Client()
	c.SetTimeout(myfunc.Seconds(client_timeout))
	c.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	return c
}

func bytes2base64(bytes []byte) string {
	result := ""
	contentType := http.DetectContentType(bytes)
	switch contentType {
	case "image/jpeg", "image/png", "image/gif":
		result = "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(bytes)
	}
	return result
}

func GetImg2base64(url string) string {
	result := ""
	c := HttpClient()
	rsp, err := c.Get(url)
	defer rsp.Close()
	if err == nil {
		if rsp.StatusCode == 200 {
			bytes, err := ioutil.ReadAll(rsp.Body)
			if err == nil {
				result = bytes2base64(bytes)
			}
		}
	}
	return result
}

func GetHTMLWithClient(client *ghttp.Client, url string) (string, error) {
	_client := client
	if _client == nil {
		_client = HttpClient()
	}
	n := 0
	for {
		rsp, err := _client.Get(url)
		if err == nil {
			defer rsp.Close()
			html := rsp.ReadAllString()
			//mylog.DebugInfo("[" + url + "]RequestHeader:\n" + rsp.RawRequest())
			//mylog.DebugInfo("[" + url + "]html:\n" + html)
			return html, nil
		} else {
			n = n + 1
			if n <= client_retry {
				mylog.Info("访问[" + url + "]出错了, 重试: " + mylog.String(n))
				continue
			} else {
				return "", err
			}
		}
	}
}

func GetBytesWithClient(client *ghttp.Client, url string) ([]byte, error) {
	_client := client
	if _client == nil {
		_client = HttpClient()
	}
	rsp, err := _client.Get(url)
	if err == nil {
		defer rsp.Close()
		return rsp.ReadAll(), nil
	} else {
		return myfunc.Empty_Bytes(), err
	}
}

func GetHTMLWithClientPost(client *ghttp.Client, url, data string) (string, error) {
	_client := client
	if _client == nil {
		_client = HttpClient()
	}
	rsp, err := _client.Post(url, data)
	if err == nil {
		defer rsp.Close()
		return rsp.ReadAllString(), nil
	} else {
		return "", err
	}
}
