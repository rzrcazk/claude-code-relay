package common

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	VerificationCodeExpiry = 10 * time.Minute
	VerificationCodeLength = 6
)

type VerificationCodeType string

const (
	EmailVerification VerificationCodeType = "email_verification"
	PasswordReset     VerificationCodeType = "password_reset"
	EmailChange       VerificationCodeType = "email_change"
	LoginVerification VerificationCodeType = "login"
)

func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

func getVerificationKey(email string, codeType VerificationCodeType) string {
	return fmt.Sprintf("verify:%s:%s", codeType, email)
}

func SendVerificationCode(email string, codeType VerificationCodeType) (string, error) {
	if email == "" {
		return "", fmt.Errorf("邮箱地址不能为空")
	}

	code := generateVerificationCode()
	key := getVerificationKey(email, codeType)

	if RDB != nil {
		ctx := context.Background()
		err := RDB.Set(ctx, key, code, VerificationCodeExpiry).Err()
		if err != nil {
			SysError("验证码存储失败: " + err.Error())
			return "", fmt.Errorf("验证码存储失败")
		}
	}

	var subject string
	var message string
	systemName := GetEmailConfig().SystemName

	switch codeType {
	case EmailVerification:
		subject = "邮箱验证码"
		message = fmt.Sprintf(`
			<h2>邮箱验证</h2>
			<p>您正在注册 %s 账号，验证码如下：</p>
			<div style="background-color: #f5f5f5; padding: 20px; text-align: center; margin: 20px 0;">
				<h1 style="color: #007cba; font-size: 32px; letter-spacing: 5px; margin: 0;">%s</h1>
			</div>
			<p>验证码有效期为 <strong>10分钟</strong>，请及时使用。</p>
			<p>如果您没有进行此操作，请忽略此邮件。</p>
		`, systemName, code)
	case LoginVerification:
		subject = "登录验证码"
		message = fmt.Sprintf(`
			<h2>登录验证</h2>
			<p>您正在登录 %s 账号，验证码如下：</p>
			<div style="background-color: #f5f5f5; padding: 20px; text-align: center; margin: 20px 0;">
				<h1 style="color: #007cba; font-size: 32px; letter-spacing: 5px; margin: 0;">%s</h1>
			</div>
			<p>验证码有效期为 <strong>10分钟</strong>，请及时使用。</p>
			<p>如果您没有进行此操作，请立即检查账号安全。</p>
		`, systemName, code)
	case PasswordReset:
		subject = "密码重置验证码"
		message = fmt.Sprintf(`
			<h2>密码重置</h2>
			<p>您正在重置 %s 账号密码，验证码如下：</p>
			<div style="background-color: #f5f5f5; padding: 20px; text-align: center; margin: 20px 0;">
				<h1 style="color: #007cba; font-size: 32px; letter-spacing: 5px; margin: 0;">%s</h1>
			</div>
			<p>验证码有效期为 <strong>10分钟</strong>，请及时使用。</p>
			<p>如果您没有进行此操作，请忽略此邮件并检查账号安全。</p>
		`, systemName, code)
	case EmailChange:
		subject = "邮箱变更验证码"
		message = fmt.Sprintf(`
			<h2>邮箱变更验证</h2>
			<p>您正在更换 %s 账号邮箱，验证码如下：</p>
			<div style="background-color: #f5f5f5; padding: 20px; text-align: center; margin: 20px 0;">
				<h1 style="color: #007cba; font-size: 32px; letter-spacing: 5px; margin: 0;">%s</h1>
			</div>
			<p>验证码有效期为 <strong>10分钟</strong>，请及时使用。</p>
			<p>如果您没有进行此操作，请立即检查账号安全。</p>
		`, systemName, code)
	default:
		return "", fmt.Errorf("未知的验证码类型")
	}

	htmlContent := fmt.Sprintf(`
	<html>
	<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
		%s
		<hr style="margin: 30px 0; border: none; border-top: 1px solid #eee;">
		<small style="color: #666;">
			此邮件由系统自动发送，请勿回复。<br>
			发送时间: %s
		</small>
	</body>
	</html>`, message, time.Now().Format("2006-01-02 15:04:05"))

	err := SendEmail(subject, email, htmlContent)
	if err != nil {
		return "", err
	}

	SysLog(fmt.Sprintf("验证码发送成功: %s, 类型: %s", email, codeType))
	return code, nil
}

func VerifyCode(email string, code string, codeType VerificationCodeType) error {
	if email == "" || code == "" {
		return fmt.Errorf("邮箱地址和验证码不能为空")
	}

	key := getVerificationKey(email, codeType)

	if RDB != nil {
		ctx := context.Background()
		storedCode, err := RDB.Get(ctx, key).Result()
		if err != nil {
			SysError("验证码获取失败: " + err.Error())
			return fmt.Errorf("验证码已过期或不存在")
		}

		if storedCode != code {
			SysLog(fmt.Sprintf("验证码验证失败: %s, 输入: %s, 期望: %s", email, code, storedCode))
			return fmt.Errorf("验证码错误")
		}

		RDB.Del(ctx, key)
		SysLog(fmt.Sprintf("验证码验证成功: %s, 类型: %s", email, codeType))
		return nil
	}

	return fmt.Errorf("Redis未配置，无法验证验证码")
}

func CheckVerificationCodeFrequency(email string, codeType VerificationCodeType) error {
	if RDB == nil {
		return nil
	}

	frequencyKey := fmt.Sprintf("verify_freq:%s:%s", codeType, email)
	ctx := context.Background()

	count, err := RDB.Get(ctx, frequencyKey).Result()
	if err == nil {
		if countInt, _ := strconv.Atoi(count); countInt >= 5 {
			return fmt.Errorf("发送过于频繁，请1小时后再试")
		}
	}

	RDB.Incr(ctx, frequencyKey)
	RDB.Expire(ctx, frequencyKey, time.Hour)

	return nil
}
