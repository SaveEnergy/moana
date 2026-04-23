package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const defaultSendGridEndpoint = "https://api.sendgrid.com/v3/mail/send"

// sendGridSender implements [PasswordResetSender] via the SendGrid v3 mail API (dynamic template).
type sendGridSender struct {
	apiKey     string
	from       string
	templateID string
	endpoint   string
	client     *http.Client
}

// NewSendGridSender returns a sender that posts to SendGrid using a [dynamic
// template](https://docs.sendgrid.com/ui/sending-email/how-to-send-an-email-with-dynamic-templates).
// apiKey, from, and templateID must be non-empty (from must be a verified sender; templateID
// is the d-… id, with a Handlebars `{{reset_url}}` in the design). If any are empty, returns nil.
func NewSendGridSender(apiKey, from, templateID string) PasswordResetSender {
	if strings.TrimSpace(apiKey) == "" || strings.TrimSpace(from) == "" || strings.TrimSpace(templateID) == "" {
		return nil
	}
	return &sendGridSender{
		apiKey:     strings.TrimSpace(apiKey),
		from:       strings.TrimSpace(from),
		templateID: strings.TrimSpace(templateID),
		endpoint:   defaultSendGridEndpoint,
		client:     &http.Client{Timeout: 30 * time.Second},
	}
}

// SendPasswordReset implements [PasswordResetSender] using the configured dynamic template
// and `reset_url` in [dynamic template data](https://docs.sendgrid.com/docs/for-developers/sending-email/using-handlebars).
func (s *sendGridSender) SendPasswordReset(ctx context.Context, toEmail, resetURL string) error {
	if s == nil {
		return fmt.Errorf("mail: nil sender")
	}
	payload := map[string]any{
		"template_id": s.templateID,
		"from": map[string]string{
			"email": s.from,
		},
		"personalizations": []map[string]any{
			{
				"to": []map[string]string{
					{"email": toEmail},
				},
				"dynamic_template_data": map[string]string{
					"reset_url": resetURL,
				},
			},
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
		// 202 Accepted: queued for send — correlate in SendGrid “Email activity” (search this id).
		if id := strings.TrimSpace(resp.Header.Get("X-Message-Id")); id != "" {
			slog.Info("sendgrid", "op", "password_reset_accepted", "x_message_id", id, "http_status", resp.StatusCode)
		} else {
			slog.Info("sendgrid", "op", "password_reset_accepted", "http_status", resp.StatusCode)
		}
		return nil
	}
	lim := io.LimitReader(resp.Body, 4096)
	b, _ := io.ReadAll(lim)
	return fmt.Errorf("sendgrid: %s: %s", resp.Status, strings.TrimSpace(string(b)))
}
