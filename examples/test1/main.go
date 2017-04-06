package main

import (
	"fmt"

	"github.com/CodyGuo/logs"
	"github.com/codyguo/SmartQQ"
)

func main() {
	qq := smartqq.NewQQClient()
	qq.OnCaptchaChange(func(this *smartqq.QQClient, data []byte) {
		this.SaveCaptach(smartqq.QQ_CAPTCHA_PNG, data)
	})

	// qq.CaptchaChange().Attach(func(this *smartqq.QQClient) {
	// 	fmt.Printf("[CaptchaChange] Prt: %s\n", this.Ptqrtoken)
	// })

	qq.OnLogined(func(this *smartqq.QQClient) {
		fmt.Printf("[OnLogined] 登录成功, --> %s\n", this.Ptqrtoken)
	})

	qq.Logined().Attach(func(this *smartqq.QQClient) {
		fmt.Printf("[Logined] 登录成功, --> %s\n", this.Ptqrtoken)
	})

	logs.Notice("开始登录QQ...")
	// faygo.Print("faygo...")
	qq.Run()
	logs.Notice("QQ登录成功...")
}
