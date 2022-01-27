package stockx_billbee_bridge

import (
	"encoding/json"
	"fmt"
	"github.com/applike/gosoline/pkg/cfg"
	"github.com/applike/gosoline/pkg/mon"
	"net/smtp"
)

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

type Mailer struct {
	logger mon.Logger
	config smtpConfig
}

type smtpConfig struct {
	Server   string
	Port     int
	Email    string
	Password string
}

type Request struct {
	config  smtpConfig
	from    string
	to      string
	subject string
	body    string
}

func NewMailer(config cfg.Config, logger mon.Logger) Mailer {
	emailConfig := smtpConfig{
		Server:   config.GetString("smtp_server"),
		Port:     config.GetInt("smtp_port"),
		Email:    config.GetString("smtp_email"),
		Password: config.GetString("smtp_password"),
	}

	return Mailer{
		logger: logger,
		config: emailConfig,
	}
}

func (m *Mailer) NewRequest(to string, subject string) *Request {
	return &Request{
		config:  m.config,
		to:      to,
		subject: subject,
	}
}

func (r *Request) sendMail(content []byte) bool {
	body := "To: " + r.to + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + string(content)
	SMTP := fmt.Sprintf("%s:%d", r.config.Server, r.config.Port)
	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", r.config.Email, r.config.Password, r.config.Server), r.config.Email, []string{r.to}, []byte(body)); err != nil {
		return false
	}
	return true
}

func (r *Request) Send(billbeeSale BillbeeSale) error {
	emailContent, err := json.Marshal(billbeeSale)

	if err != err {
		return err
	}

	if ok := r.sendMail(emailContent); ok {
		return nil
	}

	return fmt.Errorf("could not send email")
}
