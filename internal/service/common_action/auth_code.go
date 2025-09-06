package common_action

import (
	"ai_jianli_go/component"
	"ai_jianli_go/config"
	"ai_jianli_go/logs"
	"ai_jianli_go/pkg/utils"
	"ai_jianli_go/types/resp/common"
	"context"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	e "github.com/jordan-wright/email"
)

type CommonActionService struct {
}

func NewCommonActionService() *CommonActionService {
	return &CommonActionService{}
}

// 发送验证码
func (s *CommonActionService) SendAuthCode(em string) int64 {
	err := sendAuthCode(em)
	if err != nil {
		// 检查是否是 "short response" 错误，这通常表示邮件已发送但连接异常
		if strings.Contains(err.Error(), "short response") {
			logs.SugarLogger.Warn("邮件可能已发送成功，但SMTP连接异常:", err.Error())
			return common.CodeSuccess// 认为发送成功
		}
		logs.SugarLogger.Error("发送邮箱失败:", err)
		return common.CodeSendEmailFail
	}
	return common.CodeSuccess
}

// ---------------------发送验证码----------------------------------
func sendAuthCode(to string) error {
	code, err := createAuthCode(to)
	if err != nil {
		return err
	}
	subject := "【Easy Offer】邮箱验证"
	html := fmt.Sprintf(`<div style="text-align: center;">
		<h2 style="color: #333;">欢迎使用，你的验证码为：</h2>
		<h1 style="margin: 1.2em 0;">%s</h1>
		<p style="font-size: 12px; color: #666;">请在5分钟内完成验证，过期失效，请勿告知他人，以防个人信息泄露</p>
	</div>`, code)
	em := e.NewEmail()
	data := config.GetEmail()
	// 设置 sender 发送方 的邮箱 ， 此处可以填写自己的邮箱
	em.From = data.From

	// 设置 receiver 接收方 的邮箱  此处也可以填写自己的邮箱， 就是自己发邮件给自己
	em.To = []string{to}

	// 设置主题
	em.Subject = subject

	// 简单设置文件发送的内容，暂时设置成纯文本
	em.HTML = []byte(html)
	//fmt.Println(config.EmailAddr, config.Email, config.EmailAuth, config.EmailHost, config.EmailFrom)
	//设置服务器相关的配置
	err = em.Send(data.Addr, smtp.PlainAuth("", data.Email, data.Auth, data.Host))
	return err
}

// 创建随机验证码
func createAuthCode(em string) (string, error) {
	code := fmt.Sprintf("%d", utils.RandNum(900000)+100000)
	rcli := component.GetRedisDB()
	err := rcli.Set(context.Background(), em, code, 300*time.Second).Err()
	if err != nil {
		return "", err
	}
	return code, nil
}

// 校验验证码
func IdentifyCode(em string, authCode string) (int, string) {
	rcli := component.GetRedisDB()
	res, err := rcli.Get(context.Background(), em).Result()
	if err != nil {
		logs.SugarLogger.Debug("获取验证码失败:", err)
		return http.StatusUnauthorized, "验证码失效"
	}
	if res != authCode {
		return http.StatusUnauthorized, "验证码错误"
	}
	return http.StatusOK, "验证成功"
}
