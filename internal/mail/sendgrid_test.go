package mail

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendGridSender_SendPasswordReset_roundTrip(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method %s want POST", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer testkey" {
			t.Errorf("Authorization header missing or wrong: %q", r.Header.Get("Authorization"))
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type: %q want application/json", ct)
		}
		b, _ := io.ReadAll(r.Body)
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			t.Fatal(err)
		}
		from, ok := m["from"].(map[string]any)
		if !ok {
			t.Fatalf("from missing: %#v", m)
		}
		if from["email"] != "from@moana.test" {
			t.Errorf("from email: want from@moana.test got %v", from["email"])
		}
		if pers, _ := m["personalizations"].([]any); len(pers) < 1 {
			t.Fatalf("personalizations: %#v", m["personalizations"])
		} else {
			p0, _ := pers[0].(map[string]any)
			if to, _ := p0["to"].([]any); len(to) < 1 {
				t.Fatalf("to: %#v", p0["to"])
			} else {
				t0, _ := to[0].(map[string]any)
				if t0["email"] != "u@u.test" {
					t.Errorf("to email: got %v", t0["email"])
				}
			}
		}
		if subj, _ := m["subject"].(string); subj == "" {
			t.Fatal("empty subject")
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	t.Cleanup(srv.Close)

	sg := &sendGridSender{
		apiKey:   "testkey",
		from:     "from@moana.test",
		endpoint: srv.URL,
		client:   srv.Client(),
	}
	ctx := context.Background()
	if err := sg.SendPasswordReset(ctx, "u@u.test", "https://app/reset?token=1"); err != nil {
		t.Fatal(err)
	}
}

func TestNewSendGridSender_emptyReturnsNil(t *testing.T) {
	t.Parallel()
	if NewSendGridSender("", "a@b.com") != nil {
		t.Fatal("empty key")
	}
	if NewSendGridSender("key", "") != nil {
		t.Fatal("empty from")
	}
}
