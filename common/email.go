package common

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

type EmailConfig struct {
	SMTPServer      string
	SMTPPort        int
	SMTPAccount     string
	SMTPPassword    string
	SMTPFrom        string
	SMTPSSLEnabled  bool
	SystemName      string
	CacheEnabled    bool
	CacheExpireTime int
}

var emailConfig *EmailConfig

func InitEmailConfig() {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if port == 0 {
		port = 587 // 默认端口
	}

	sslEnabled, _ := strconv.ParseBool(os.Getenv("SMTP_SSL_ENABLED"))
	cacheEnabled, _ := strconv.ParseBool(os.Getenv("EMAIL_CACHE_ENABLED"))
	cacheExpire, _ := strconv.Atoi(os.Getenv("EMAIL_CACHE_EXPIRE_TIME"))
	if cacheExpire == 0 {
		cacheExpire = 300 // 默认5分钟
	}

	smtpFrom := os.Getenv("SMTP_FROM")
	if smtpFrom == "" {
		smtpFrom = os.Getenv("SMTP_ACCOUNT") // 兼容性处理
	}

	systemName := os.Getenv("SYSTEM_NAME")
	if systemName == "" {
		systemName = "Claude Code Relay"
	}

	emailConfig = &EmailConfig{
		SMTPServer:      os.Getenv("SMTP_SERVER"),
		SMTPPort:        port,
		SMTPAccount:     os.Getenv("SMTP_ACCOUNT"),
		SMTPPassword:    os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:        smtpFrom,
		SMTPSSLEnabled:  sslEnabled,
		SystemName:      systemName,
		CacheEnabled:    cacheEnabled,
		CacheExpireTime: cacheExpire,
	}
}

func GetEmailConfig() *EmailConfig {
	if emailConfig == nil {
		InitEmailConfig()
	}
	return emailConfig
}

func generateMessageID() (string, error) {
	config := GetEmailConfig()
	split := strings.Split(config.SMTPFrom, "@")
	if len(split) < 2 {
		return "", fmt.Errorf("invalid SMTP account")
	}
	domain := split[1]
	return fmt.Sprintf("<%d.%s@%s>", time.Now().UnixNano(), GenerateRandomString(12), domain), nil
}

func isOutlookServer(account string) bool {
	return strings.Contains(account, "@outlook.") || strings.Contains(account, "@hotmail.") ||
		strings.Contains(account, "@live.") || strings.Contains(account, "@msn.")
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unknown fromServer")
		}
	}
	return nil, nil
}

func SendEmail(subject string, receiver string, content string) error {
	config := GetEmailConfig()

	if config.SMTPServer == "" || config.SMTPAccount == "" || config.SMTPPassword == "" {
		return fmt.Errorf("SMTP服务器配置不完整，请检查SMTP_SERVER、SMTP_ACCOUNT、SMTP_PASSWORD环境变量")
	}

	if receiver == "" {
		return fmt.Errorf("收件人地址不能为空")
	}

	if config.SMTPFrom == "" {
		config.SMTPFrom = config.SMTPAccount
	}

	// 检查缓存（防止短时间内重复发送相同邮件）
	if config.CacheEnabled {
		contentLen := len(content)
		if contentLen > 36 {
			contentLen = 36
		}
		truncatedContent := content[:contentLen]
		hash := md5.Sum([]byte(truncatedContent + receiver))
		cacheMD5Key := "email_cache:" + hex.EncodeToString(hash[:])

		if RDB != nil {
			ctx := context.Background()
			redisValue, _ := RDB.Get(ctx, cacheMD5Key).Result()
			if redisValue != "" {
				SysLog("邮件发送被缓存跳过: " + receiver)
				return nil
			}
			defer func() {
				_ = RDB.Set(ctx, cacheMD5Key, "1", time.Duration(config.CacheExpireTime)*time.Second).Err()
			}()
		}
	}

	messageID, err := generateMessageID()
	if err != nil {
		return fmt.Errorf("生成Message-ID失败: %w", err)
	}

	// 编码主题（支持中文）
	encodedSubject := fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(subject)))

	// 构造邮件内容
	mail := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s<%s>\r\n"+
		"Subject: %s\r\n"+
		"Date: %s\r\n"+
		"Message-ID: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n\r\n%s\r\n",
		receiver, config.SystemName, config.SMTPFrom, encodedSubject,
		time.Now().Format(time.RFC1123Z), messageID, content))

	auth := smtp.PlainAuth("", config.SMTPAccount, config.SMTPPassword, config.SMTPServer)
	addr := fmt.Sprintf("%s:%d", config.SMTPServer, config.SMTPPort)
	to := strings.Split(receiver, ";")

	var sendErr error

	// 根据端口和SSL配置选择发送方式
	if config.SMTPPort == 465 || config.SMTPSSLEnabled {
		// SSL/TLS方式
		sendErr = sendEmailWithTLS(addr, auth, config, mail, to, receiver)
	} else if isOutlookServer(config.SMTPAccount) || config.SMTPServer == "smtp.azurecomm.net" {
		// Outlook需要特殊的LOGIN认证
		auth = LoginAuth(config.SMTPAccount, config.SMTPPassword)
		sendErr = smtp.SendMail(addr, auth, config.SMTPFrom, to, mail)
	} else {
		// 标准SMTP方式
		sendErr = smtp.SendMail(addr, auth, config.SMTPFrom, to, mail)
	}

	if sendErr != nil {
		SysError("邮件发送失败: " + sendErr.Error())
		return fmt.Errorf("邮件发送失败: %w", sendErr)
	}

	SysLog("邮件发送成功: " + receiver)
	return nil
}

func sendEmailWithTLS(addr string, auth smtp.Auth, config *EmailConfig, mail []byte, to []string, receiver string) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         config.SMTPServer,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS连接失败: %w", err)
	}

	client, err := smtp.NewClient(conn, config.SMTPServer)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %w", err)
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP认证失败: %w", err)
	}

	if err = client.Mail(config.SMTPFrom); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}

	receiverEmails := strings.Split(receiver, ";")
	for _, rcpt := range receiverEmails {
		if err = client.Rcpt(rcpt); err != nil {
			return fmt.Errorf("设置收件人失败 %s: %w", rcpt, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("获取数据写入器失败: %w", err)
	}

	_, err = w.Write(mail)
	if err != nil {
		return fmt.Errorf("写入邮件数据失败: %w", err)
	}

	return w.Close()
}

func SendPlainTextEmail(subject string, receiver string, content string) error {
	htmlContent := fmt.Sprintf("<html><body><pre style=\"white-space: pre-wrap; word-wrap: break-word;\">%s</pre></body></html>", content)
	return SendEmail(subject, receiver, htmlContent)
}

func SendHTMLEmail(subject string, receiver string, htmlContent string) error {
	return SendEmail(subject, receiver, htmlContent)
}

func SendSystemNotificationEmail(receiver string, title string, message string) error {
	subject := fmt.Sprintf("[%s] %s", GetEmailConfig().SystemName, title)
	htmlContent := fmt.Sprintf(`
<html>
<body>
<h2>系统通知</h2>
<h3>%s</h3>
<div style="padding: 10px; border-left: 4px solid #007cba; background-color: #f8f9fa;">
<p>%s</p>
</div>
<hr>
<small>发送时间: %s</small>
</body>
</html>`, title, strings.ReplaceAll(message, "\n", "<br>"), time.Now().Format("2006-01-02 15:04:05"))

	return SendEmail(subject, receiver, htmlContent)
}
