package service

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"

	"go-gateway/internal/config"
)

// EmailService 邮件服务
type EmailService struct {
	cfg *config.EmailConfig
}

// NewEmailService 创建邮件服务
func NewEmailService(cfg *config.EmailConfig) *EmailService {
	return &EmailService{cfg: cfg}
}

// IsEnabled 检查邮件服务是否启用
func (s *EmailService) IsEnabled() bool {
	return s.cfg.Enabled
}

// SendPasswordResetEmail 发送密码重置邮件
func (s *EmailService) SendPasswordResetEmail(toEmail, resetToken, resetURL string) error {
	if !s.cfg.Enabled {
		return nil // 未启用时静默返回
	}

	subject := "密码重置请求 - NovelAI"

	// 构建重置链接
	if resetURL == "" {
		resetURL = fmt.Sprintf("http://localhost:3000/reset-password?token=%s", resetToken)
	} else {
		resetURL = fmt.Sprintf("%s?token=%s", resetURL, resetToken)
	}

	// HTML 邮件模板
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #1890ff; color: white; padding: 20px; text-align: center; }
        .content { padding: 30px; background: #f9f9f9; }
        .button { display: inline-block; padding: 12px 30px; background: #1890ff; color: white; text-decoration: none; border-radius: 4px; margin: 20px 0; }
        .footer { padding: 20px; text-align: center; color: #666; font-size: 12px; }
        .warning { color: #ff4d4f; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>NovelAI</h1>
        </div>
        <div class="content">
            <h2>密码重置请求</h2>
            <p>您好，</p>
            <p>我们收到了您的密码重置请求。请点击下方按钮重置您的密码：</p>
            <p style="text-align: center;">
                <a href="{{.ResetURL}}" class="button">重置密码</a>
            </p>
            <p>或者复制以下链接到浏览器：</p>
            <p style="word-break: break-all; background: #eee; padding: 10px;">{{.ResetURL}}</p>
            <p class="warning">此链接将在 1 小时后失效。如果您没有请求重置密码，请忽略此邮件。</p>
        </div>
        <div class="footer">
            <p>此邮件由系统自动发送，请勿回复。</p>
            <p>&copy; 2026 NovelAI. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

	tmpl, err := template.New("email").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	data := struct {
		ResetURL string
	}{
		ResetURL: resetURL,
	}
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	return s.sendHTML(toEmail, subject, body.String())
}

// SendWelcomeEmail 发送欢迎邮件
func (s *EmailService) SendWelcomeEmail(toEmail, username string) error {
	if !s.cfg.Enabled {
		return nil
	}

	subject := "欢迎加入 NovelAI"

	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #1890ff; color: white; padding: 20px; text-align: center; }
        .content { padding: 30px; background: #f9f9f9; }
        .footer { padding: 20px; text-align: center; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>NovelAI</h1>
        </div>
        <div class="content">
            <h2>欢迎加入 NovelAI！</h2>
            <p>亲爱的 {{.Username}}，</p>
            <p>感谢您注册 NovelAI - 智能小说创作助手！</p>
            <p>现在您可以：</p>
            <ul>
                <li>创建和管理您的小说项目</li>
                <li>使用 AI 辅助生成章节内容</li>
                <li>构建丰富的世界观和角色设定</li>
                <li>管理伏笔和剧情线索</li>
            </ul>
            <p>开始您的创作之旅吧！</p>
        </div>
        <div class="footer">
            <p>&copy; 2026 NovelAI. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

	tmpl, err := template.New("email").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	data := struct {
		Username string
	}{
		Username: username,
	}
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	return s.sendHTML(toEmail, subject, body.String())
}

// sendHTML 发送 HTML 邮件
func (s *EmailService) sendHTML(to, subject, htmlBody string) error {
	from := s.cfg.FromAddress
	if s.cfg.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.cfg.FromName, s.cfg.FromAddress)
	}

	// 构建邮件头
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var message bytes.Buffer
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.WriteString(htmlBody)

	// SMTP 认证
	auth := smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost)

	// 连接 SMTP 服务器
	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)

	// 使用 TLS
	if s.cfg.SMTPPort == 465 {
		// SSL/TLS 直接连接
		return s.sendWithTLS(addr, auth, s.cfg.FromAddress, to, message.Bytes())
	}

	// STARTTLS
	return smtp.SendMail(addr, auth, s.cfg.FromAddress, []string{to}, message.Bytes())
}

// sendWithTLS 使用 TLS 发送邮件 (端口 465)
func (s *EmailService) sendWithTLS(addr string, auth smtp.Auth, from, to string, msg []byte) error {
	tlsConfig := &tls.Config{
		ServerName: s.cfg.SMTPHost,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.cfg.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	if err := client.Mail(from); err != nil {
		return fmt.Errorf("SMTP MAIL command failed: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("SMTP RCPT command failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA command failed: %w", err)
	}

	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close email writer: %w", err)
	}

	return client.Quit()
}
