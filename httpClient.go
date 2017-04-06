package smartqq

import (
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

const (
	DEFAULT_TIMEOUT = 30
)

type Client struct {
	Client  http.Client
	timeout time.Duration
}

func newClient() *Client {
	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar, Timeout: DEFAULT_TIMEOUT * time.Second}

	return &Client{Client: client}
}

func (this *Client) newRequest(method, urlStr, body string) (req *http.Request, err error) {
	req, err = http.NewRequest(method, urlStr, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	if this.timeout != 0 {
		this.Client.Timeout = this.timeout
	}

	return req, nil
}

func (this *Client) get(urlStr string) (resp *http.Response, err error) {
	req, err := this.newRequest("GET", urlStr, "")
	if err != nil {
		return nil, err
	}

	return this.Client.Do(req)
}

func (this *Client) post(urlStr, data string) (resp *http.Response, err error) {
	req, err := this.newRequest("POST", urlStr, data)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	return this.Client.Do(req)
}

func (this *Client) cookies(resp *http.Response) []*http.Cookie {
	cookies := this.Client.Jar.Cookies(resp.Request.URL)

	return cookies
}

func (this *Client) updateCookie(urlStr string) error {
	resp, err := this.get(urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
