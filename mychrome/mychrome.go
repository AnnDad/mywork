package mychrome

import (
	//"anndad/gftool"
	"anndad/myfunc"
	"anndad/mylog"
	"context"
	"io/ioutil"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	//"github.com/gogf/gf/container/glist"
	//"github.com/gogf/gf/container/gtype"
)

type TMyChromeOption struct {
	ShowBrowser         bool   //是否可视化
	ShowImage           bool   //是否加载图片
	ShowMaximized       bool   //是否最大化显示
	ShowScrollbar       bool   //是否显示滚动条
	ShowAutomationLabel bool   //是否显示自动化标签
	PlaySound           bool   //是否播放声音
	SandboxMode         bool   //沙盒模式
	UserAgent           string //浏览器标识
	ChromePath          string
	ProxyAddr           string //代理设置
	Width               int    //
	Height              int    //
}

type TMyChrome struct {
	ctxBrowser     context.Context
	ctxFirstTab    context.Context
	cancelFirstTab context.CancelFunc
	Options        TMyChromeOption
}

func NewMyChrome() *TMyChrome {
	result := new(TMyChrome)
	result.SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")
	result.SetPlaySound(false)
	result.SetShowBrowser(false)
	result.SetShowImage(true)
	result.SetSandBoxMode(false)
	result.SetShowScrollbar(true)
	result.SetShowMaximized(false)
	result.SetSize(1024, 768)
	result.SetShowAutomationLabel(false)
	result.SetChromePath("")
	result.SetProxyAddr("")
	result.ctxBrowser = nil
	return result
}

func (this *TMyChrome) Debug_TabStatus() {
	status, err := chromedp.Targets(this.ctxFirstTab)
	if err != nil {
		SaveLog("获取状态出错: " + err.Error())
	}
	n := len(status)
	for i := 0; i < n; i++ {
		item := "窗口" + myfunc.String(i)
		SaveLog(item + ".TargetID: " + status[i].TargetID.String())
		SaveLog(item + ".Type: " + status[i].Type)
		SaveLog(item + ".Title: " + status[i].Title)
		SaveLog(item + ".URL: " + status[i].URL)
		SaveLog(item + ".Attached: " + myfunc.Bool2Str(status[i].Attached))
		SaveLog(item + ".OpenerID: " + status[i].OpenerID.String())
		SaveLog(item + ".CanAccessOpener: " + myfunc.Bool2Str(status[i].CanAccessOpener))
		SaveLog(item + ".OpenerFrameID: " + status[i].OpenerFrameID.String())
		SaveLog(item + ".BrowserContextID: " + status[i].BrowserContextID.String())
	}
}

func (this *TMyChrome) OpenBrowser() {
	options := chromedp.DefaultExecAllocatorOptions[:]
	options = append(options, chromedp.Flag("headless", !this.Options.ShowBrowser))
	options = append(options, chromedp.UserAgent(this.Options.UserAgent))
	if this.Options.ShowMaximized {
		options = append(options, chromedp.Flag("start-maximized", true))
	} else {
		options = append(options, chromedp.WindowSize(this.Options.Width, this.Options.Height))
	}
	options = append(options, chromedp.Flag("enable-automation", this.Options.ShowAutomationLabel))
	//options = append(options, chromedp.Flag("incognito", true)) //无痕模式, 如果该模式为true,chromedp将不会以多标签页的方式工作.
	options = append(options, chromedp.Flag("blink-settings", "imagesEnabled="+myfunc.Bool2Str(this.Options.ShowImage)))
	options = append(options, chromedp.DisableGPU)
	if this.Options.ChromePath != "" {
		options = append(options, chromedp.ExecPath(this.Options.ChromePath))
	}
	if this.Options.ProxyAddr != "" {
		options = append(options, chromedp.ProxyServer(this.Options.ProxyAddr))
	}
	options = append(options, chromedp.NoDefaultBrowserCheck)
	//options = append(options, chromedp.Flag("restore-last-session", false))
	this.ctxBrowser, _ = chromedp.NewExecAllocator(context.Background(), options...)
	this.ctxFirstTab, this.cancelFirstTab = chromedp.NewContext(this.ctxBrowser)
	if err := chromedp.Run(this.ctxFirstTab); err != nil {
		SaveLog("RunBrowser: " + err.Error())
		panic("RunBrowser: " + err.Error())
	} else {
		SaveLog("OpenBrowser Done!")
	}
}

func (this *TMyChrome) CloseBrowser() {
	if this.cancelFirstTab != nil {
		SaveLog("CloseBrowser...")
		this.cancelFirstTab()
		SaveLog("CloseBrowser Done!")
	} else {
		SaveLog("Err_CloseBrowser: cancelFirstTab is nil")
	}
}

func ScreenshotAsBytes(ctxTab context.Context) ([]byte, error) {
	var buf []byte
	cxtdo, _ := context.WithTimeout(ctxTab, time.Duration(15)*time.Second)
	err := chromedp.Run(
		cxtdo,
		chromedp.FullScreenshot(&buf, 100),
	)
	return buf, err
}

func ScreenshotAsFile(ctxTab context.Context, Path string) error {
	buf, err := ScreenshotAsBytes(ctxTab)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(Path, buf, 0o644); err != nil {
		return err
	}
	return nil
}

func (this *TMyChrome) OpenTab(timeout_sec int, actions ...chromedp.Action) (context.Context, context.CancelFunc, error) {
	if this.ctxBrowser == nil {
		panic("Err_Call_OpenBrowser_First")
	}
	cxtTab, cancelTab := chromedp.NewContext(this.ctxFirstTab)
	cxt, _ := context.WithTimeout(cxtTab, time.Duration(timeout_sec)*time.Second)
	//this.tabOpening.PushBack(cxt)
	err := chromedp.Run(cxt, actions...)
	if err != nil {
		SaveLog("Err_OpenTab: " + err.Error())
	} else {
		SaveLog("OpenTab Done!")
	}
	return cxtTab, cancelTab, err
}

func (this *TMyChrome) DisableNotify() {
	this.DoWithTab(60, this.ctxFirstTab,
		chromedp.Navigate("chrome://settings/content/notifications"),
		SleepSecond(1),
		ExecJavascript(`document.querySelector("body > settings-ui").shadowRoot.querySelector("#main").shadowRoot.querySelector("settings-basic-page").shadowRoot.querySelector("#basicPage > settings-section.expanded > settings-privacy-page").shadowRoot.querySelector("#notificationRadioGroup > settings-radio-group > settings-collapse-radio-button:nth-child(4)").shadowRoot.querySelector("#button > div.disc").click();`),
	)
}

func (this *TMyChrome) DoWithTab(timeout_sec int, ctxTab context.Context, actions ...chromedp.Action) error {
	//mylog.Println(ctxTab.Deadline())
	ctxDo, _ := context.WithTimeout(ctxTab, time.Duration(timeout_sec)*time.Second)
	//mylog.Println(ctxDo.Deadline())
	err := chromedp.Run(ctxDo, actions...)
	if err != nil {
		SaveLog("Err_DoWithTab: " + err.Error())
	} else {
		SaveLog("DoWithTab Done!")
	}
	return err
}

func (this *TMyChrome) CloseTab(ctxTab context.Context, cancel context.CancelFunc, closebrowser ...bool) {
	if cancel != nil {
		SaveLog("CloseTab...")
		cancel()
		SaveLog("CloseTab Done!")
	} else {
		SaveLog("Err_CloseTab: cancel is nil")
	}
	if len(closebrowser) > 0 {
		if closebrowser[0] {
			this.CloseBrowser()
		}
	}
}

func SaveLog(info string) {
	mylog.Info(info)
}

func SleepMillisecond(Millisecond int) chromedp.Action {
	return chromedp.Sleep(time.Duration(Millisecond) * time.Millisecond)
}

func SleepSecond(second int) chromedp.Action {
	return chromedp.Sleep(time.Duration(second) * time.Second)
}

func ExecJavascript(script string) chromedp.Action {
	return chromedp.EvaluateAsDevTools(script, nil)
}

func ExecJavascriptResult(script string, result interface{}) chromedp.Action {
	return chromedp.EvaluateAsDevTools(script, &result)
}

func SleepRandom(min, max int) chromedp.Action {
	second := myfunc.Random(min, max)
	return chromedp.Sleep(time.Duration(second) * time.Second)
}

func SaveCookiesToFile(path string) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		// cookies的获取对应是在devTools的network面板中
		// 1. 获取cookies
		cookies, err := network.GetAllCookies().Do(ctx)
		if err != nil {
			return err
		}
		// 2. 序列化
		cookiesData, err := network.GetAllCookiesReturns{Cookies: cookies}.MarshalJSON()
		if err != nil {
			return err
		}
		myfunc.SaveFileBytes(true, path, cookiesData)
		return nil
	}
}

func SaveCookies(cookies *string) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		// cookies的获取对应是在devTools的network面板中
		// 1. 获取cookies
		cookies_t, err := network.GetAllCookies().Do(ctx)
		if err != nil {
			myfunc.Println("err1: ", err.Error())
			return err
		}
		//myfunc.Println("cookies_t: ", cookies_t)
		// 2. 序列化
		cookiesData, err := network.GetAllCookiesReturns{Cookies: cookies_t}.MarshalJSON()
		if err != nil {
			myfunc.Println("err1: ", err.Error())
			return err
		}
		*cookies = myfunc.Bytes2str(cookiesData)
		return nil
	}
}

func LoadCookies(path string) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		// 如果存在则读取cookies的数据
		cookiesData, err := myfunc.ReadFileAsBytes(path)
		if err != nil {
			return err
		}

		// 反序列化
		cookiesParams := network.SetCookiesParams{}
		if err = cookiesParams.UnmarshalJSON(cookiesData); err != nil {
			return err
		}

		// 设置cookies
		return network.SetCookies(cookiesParams.Cookies).Do(ctx)
	}
}

func (this *TMyChrome) SetShowAutomationLabel(value bool) {
	this.Options.ShowAutomationLabel = value
}
func (this *TMyChrome) SetShowBrowser(value bool) {
	this.Options.ShowBrowser = value
}
func (this *TMyChrome) SetShowImage(value bool) {
	this.Options.ShowImage = value
}
func (this *TMyChrome) SetShowScrollbar(value bool) {
	this.Options.ShowScrollbar = value
}
func (this *TMyChrome) SetSandBoxMode(value bool) {
	this.Options.SandboxMode = value
}
func (this *TMyChrome) SetPlaySound(value bool) {
	this.Options.PlaySound = value
}
func (this *TMyChrome) SetShowMaximized(value bool) {
	this.Options.ShowMaximized = value
}
func (this *TMyChrome) SetUserAgent(value string) {
	this.Options.UserAgent = value
}
func (this *TMyChrome) SetChromePath(value string) {
	this.Options.ChromePath = value
}

func (this *TMyChrome) SetProxyAddr(value string) {
	this.Options.ProxyAddr = value
}
func (this *TMyChrome) SetSize(width, height int) {
	this.Options.Width = width
	this.Options.Height = height
}
