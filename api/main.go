package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gocolly/colly"
)

func main() {
	// 创建一个新的Collector
	c := colly.NewCollector()

	// 设置登录页面的URL
	loginURL := "https://sjdxykt.sztu.edu.cn/sso/doLogin"

	// 设置验证码图片的URL
	captchaURL := "https://sjdxykt.sztu.edu.cn/sso/captchaCode"

	// 下载验证码图片
	resp, err := http.Get(captchaURL)
	if err != nil {
		fmt.Println("Failed to download captcha:", err)
		return
	}
	defer resp.Body.Close()

	// 保存验证码图片到本地
	captchaFile, err := os.Create("captcha.jpg")
	if err != nil {
		fmt.Println("Failed to create captcha file:", err)
		return
	}
	defer captchaFile.Close()

	_, err = io.Copy(captchaFile, resp.Body)
	if err != nil {
		fmt.Println("Failed to read captcha response body:", err)
		return
	}

	// 提示用户输入验证码
	var captcha string
	fmt.Println("Please open captcha.jpg and enter the captcha:")
	fmt.Scanln(&captcha)

	// 设置登录表单数据
	loginData := map[string]string{
		"loginType": "rftSigner",
		"account":   "202200201079",
		"password":  "271530",
		"captcha":   captcha,
	}

	// 处理登录请求
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Login response received")
		fmt.Println(string(r.Body))
	})

	// 提交登录表单
	err = c.Post(loginURL, loginData)
	if err != nil {
		fmt.Println("Failed to submit login form:", err)
		return
	}

	// 开始爬取其他页面
	// c.Visit("https://example.com/protected_page")
}
