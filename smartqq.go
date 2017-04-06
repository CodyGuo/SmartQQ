package smartqq

import (
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/CodyGuo/logs"
)

const (
	QQ_CAPTCHA_PNG = "qq.png"
)

type QQClient struct {
	*Client
	Ptqrtoken       string
	IsLogin         bool
	oncaptchaChange func(*QQClient, []byte)
	onLogined       func(*QQClient)

	captchaPublisher EventPublisher
	loginedPublisher EventPublisher
}

func NewQQClient() *QQClient {
	qqClient := new(QQClient)
	qqClient.Client = newClient()

	return qqClient
}

func (this *QQClient) OnCaptchaChange(fn func(*QQClient, []byte)) {
	this.oncaptchaChange = fn
}

func (this *QQClient) OnLogined(fn func(*QQClient)) {
	this.onLogined = fn
}

func (this *QQClient) CaptchaChange() *Event {
	return this.captchaPublisher.Event()
}

func (this *QQClient) Logined() *Event {
	return this.loginedPublisher.Event()
}

func (this *QQClient) SaveCaptach(f string, data []byte) error {
	wr, err := os.Create(f)
	if err != nil {
		return err
	}
	_, err = wr.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (this *QQClient) getCaptcha() ([]byte, error) {
	urlStr := `https://ssl.ptlogin2.qq.com/ptqrshow?appid=501004106&e=0&l=M&s=5&d=72&v=4&t=0.8`

	resp, err := this.get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	this.getPtqrtoken(resp)
	if err := this.updateCookie(); err != nil {
		return nil, err
	}
	return data, nil
}

func (this *QQClient) checkCaptach() (string, error) {
	urlStr := "https://ssl.ptlogin2.qq.com/ptqrlogin?ptqrtoken=" + this.Ptqrtoken + "&webqq_type=10&remember_uin=1&login2qq=1&aid=501004106&u1=http%3A%2F%2Fw.qq.com%2Fproxy.html%3Flogin2qq%3D1%26webqq_type%3D10&ptredirect=0&ptlang=2052&daid=164&from_ui=1&pttype=1&dumy=&fp=loginerroralert&action=0-0-263174&mibao_css=m_webqq&t=undefined&g=1&js_type=0&js_ver=10197&login_sig=&pt_randsalt=0"

	// logs.Noticef("ptqrtoken --> %s", this.Ptqrtoken)
	// logs.Noticef("urlStr --> %s", urlStr)
	resp, err := this.get(urlStr)
	if err != nil {
		return "-2", err
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	// logs.Noticef("data --> %s", data)
	regexp_image_status := regexp.MustCompile(`\d+`)
	code := regexp_image_status.FindAllString(string(data), 1)[0]

	return code, nil
}

func (this *QQClient) updateCookie() error {
	urlStr := "https://ui.ptlogin2.qq.com/cgi-bin/login?daid=164&target=self&style=16&mibao_css=m_webqq&appid=501004106&enable_qlogin=0&no_verifyimg=1&s_url=http%3A%2F%2Fw.qq.com%2Fproxy.html&f_url=loginerroralert&strong_login=1&login_state=10&t=20131024001"
	return this.Client.updateCookie(urlStr)
}

func (this *QQClient) getPtqrtoken(resp *http.Response) {
	cookies := this.cookies(resp)
	for _, cookie := range cookies {
		if cookie.Name == "qrsig" {
			this.Ptqrtoken = hash33(cookie.Value)
		}
	}
}

func (this *QQClient) SetTimeout(t time.Duration) {
	this.timeout = t
}

func (this *QQClient) Run() {
	// 验证码回调
	data, err := this.getCaptcha()
	if err != nil {
		logs.Fatal(err)
	}
	this.oncaptchaChange(this, data)

	n := 1
	for {
		logs.Noticef("二维码接口验证 [%d] 次.", n)
		code, err := this.checkCaptach()
		if err != nil {
			logs.Noticef("Run CheckCaptach : %v\n", err)
			return
		}

		switch code {
		case "65":
			logs.Warning("二维码已失效.")
			data, err := this.getCaptcha()
			if err != nil {
				logs.Fatal(err)
			}
			this.oncaptchaChange(this, data)
			n = 0
		case "66":
			logs.Notice("二维码未失效.")
		case "67":
			logs.Notice("二维码正在验证中...")
		case "0":
			logs.Notice("二维码验证成功.")
			this.loginedPublisher.Publish(this)
			this.onLogined(this)
			return
		case "403":
			logs.Errorf("可能是接口错误,请检查接口. (%s)\n", code)
			return
		default:
			logs.Errorf("未知状态码 (%s)", code)
			return
		}
		this.captchaPublisher.Publish(this)
		time.Sleep(1 * time.Second)
		n++
	}

}
