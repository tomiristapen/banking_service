package smtp

import (
    "fmt"
    "gopkg.in/gomail.v2"
    "os"
)

type Mailer interface {
    SendVerificationEmail(to string, code string) error
}

type smtpMailer struct {
    from     string
    password string
}

func NewMailer() Mailer {
    return &smtpMailer{
        from:     os.Getenv("SMTP_EMAIL"),
        password: os.Getenv("SMTP_PASSWORD"),
    }
}

func (m *smtpMailer) SendVerificationEmail(to string, code string) error {
    msg := gomail.NewMessage()
    msg.SetHeader("From", m.from)
    msg.SetHeader("To", to)
    msg.SetHeader("Subject", "Email Verification")
    msg.SetBody("text/plain", fmt.Sprintf("Your verification code is: %s", code))

    dialer := gomail.NewDialer("smtp.gmail.com", 587, m.from, m.password)
    return dialer.DialAndSend(msg)
}
