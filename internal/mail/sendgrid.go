package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultSendGridEndpoint = "https://api.sendgrid.com/v3/mail/send"

// sendGridSender implements [PasswordResetSender] via the SendGrid v3 mail API.
type sendGridSender struct {
	apiKey   string
	from     string
	endpoint string
	client   *http.Client
}

// NewSendGridSender returns a sender that posts to SendGrid. apiKey and from must both be
// non-empty (from must be a verified sender in SendGrid). If either is empty, returns nil.
func NewSendGridSender(apiKey, from string) PasswordResetSender {
	if strings.TrimSpace(apiKey) == "" || strings.TrimSpace(from) == "" {
		return nil
	}
	return &sendGridSender{
		apiKey:   strings.TrimSpace(apiKey),
		from:     strings.TrimSpace(from),
		endpoint: defaultSendGridEndpoint,
		client:   &http.Client{Timeout: 30 * time.Second},
	}
}

// SendPasswordReset implements [PasswordResetSender] with a plain-text body.
func (s *sendGridSender) SendPasswordReset(ctx context.Context, toEmail, resetURL string) error {
	if s == nil {
		return fmt.Errorf("mail: nil sender")
	}
	subj := "Reset your Moana password"
	body := "Someone requested a password reset for your Moana account.\n\n" +
		"Open this link to choose a new password (it expires soon):\n" + resetURL + "\n\n" +
		"If you did not request this, you can ignore this email.\n"
	payload := map[string]any{
		"personalizations": []map[string]any{
			{
				"to": []map[string]string{
					{"email": toEmail},
				},
			},
		},
		"from": map[string]string{
			"email": s.from,
		},
		"subject": subj,
		"content": []map[string]string{
			{"type": "text/plain", "value": body},
		},
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); _ = resp.Body.Close() }()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	lim := io.LimitReader(resp.Body, 4096)
	b, _ := io.ReadAll(lim)
	return fmt.Errorf("sendgrid: %s: %s", resp.Status, strings.TrimSpace(string(b)))
}
