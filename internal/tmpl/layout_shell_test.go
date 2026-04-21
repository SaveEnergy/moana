package tmpl

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	"moana/internal/render"
	"moana/internal/store"
)

func TestLayoutShell_embedded_notificationBadge(t *testing.T) {
	t.Parallel()
	root, err := Parse()
	if err != nil {
		t.Fatal(err)
	}
	u := &store.User{Email: "shell@moana.test"}
	base := render.LayoutData{
		Title:  "T",
		User:   u,
		Year:   2026,
		Active: "dashboard",
		RepoURL: "https://example.com/repo",
		Body:   template.HTML("<p>x</p>"),
	}

	run := func(ld render.LayoutData) string {
		t.Helper()
		var buf bytes.Buffer
		if err := root.ExecuteTemplate(&buf, "layout.html", ld); err != nil {
			t.Fatal(err)
		}
		return buf.String()
	}

	t.Run("no badge when zero", func(t *testing.T) {
		t.Parallel()
		ld := base
		ld.UnreadNotificationCount = 0
		s := run(ld)
		if strings.Contains(s, `class="app-notif-badge"`) {
			t.Fatalf("unexpected badge in %q", s[:min(800, len(s))])
		}
	})

	t.Run("badge with count", func(t *testing.T) {
		t.Parallel()
		ld := base
		ld.UnreadNotificationCount = 3
		s := run(ld)
		if !strings.Contains(s, `class="app-notif-badge"`) || !strings.Contains(s, ">3<") {
			t.Fatalf("expected badge 3, got prefix %q", s[:min(1200, len(s))])
		}
	})

	t.Run("badge caps at 99+", func(t *testing.T) {
		t.Parallel()
		ld := base
		ld.UnreadNotificationCount = 150
		s := run(ld)
		if !strings.Contains(s, `class="app-notif-badge"`) || !strings.Contains(s, "99+") {
			t.Fatalf("expected 99+ cap, got prefix %q", s[:min(1200, len(s))])
		}
	})
}
