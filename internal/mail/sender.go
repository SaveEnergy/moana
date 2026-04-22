package mail

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

// PasswordResetSender is implemented by [Sender] for dependency injection in [handlers.App].
type PasswordResetSender interface {
	SendPasswordReset(ctx context.Context, toEmail, resetURL string) error
}

// Sender sends email via SMTP (STARTTLS when advertised).
type Sender struct {
	host       string
	port       int
	user       string
	password   string
	from       string
	serverName string
}

// NewSender returns an SMTP client. from must be a valid envelope sender (e.g. noreply@example.com).
func NewSender(host string, port int, user, password, from string) *Sender {
	if host == "" || from == "" {
		return nil
	}
	if port == 0 {
		port = 587
	}
	return &Sender{
		host:       host,
		port:       port,
		user:       user,
		password:   password,
		from:       from,
		serverName: host,
	}
}

// SendPasswordReset implements [PasswordResetSender] with a plain-text body.
func (s *Sender) SendPasswordReset(_ context.Context, toEmail, resetURL string) error {
	if s == nil {
		return fmt.Errorf("mail: nil sender")
	}
	subj := "Reset your Moana password"
	body := "Someone requested a password reset for your Moana account.\n\n" +
		"Open this link to choose a new password (it expires soon):\n" + resetURL + "\n\n" +
		"If you did not request this, you can ignore this email.\n"
	msg := []byte(fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		toEmail, s.from, subj, body))
	addr := net.JoinHostPort(s.host, fmt.Sprint(s.port))
	return smtpSend(addr, s.serverName, s.user, s.password, s.from, []string{toEmail}, msg)
}

func smtpSend(addr, serverName, user, password, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer func() { _ = c.Close() }()

	if ok, _ := c.Extension("STARTTLS"); ok {
		tcfg := &tls.Config{ServerName: serverName}
		if err := c.StartTLS(tcfg); err != nil {
			return err
		}
	}
	if user != "" || password != "" {
		auth := smtp.PlainAuth("", user, password, serverName)
		if err := c.Auth(auth); err != nil {
			return err
		}
	}
	if err := c.Mail(from); err != nil {
		return err
	}
	for i := range to {
		if err := c.Rcpt(to[i]); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return c.Quit()
}

// PublicResetURL builds an absolute https? URL for a path starting with /.
func PublicResetURL(publicBase, path string) string {
	return strings.TrimRight(publicBase, "/") + path
}
