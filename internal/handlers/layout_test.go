package handlers

import (
	"testing"

	"moana/internal/store"
)

func TestLayoutShell_mainClassEmptyWhileShellMainSets(t *testing.T) {
	t.Parallel()
	u := &store.User{ID: 1, Email: "x@y.z"}
	shell := layoutShell("Hi", "dash", u)
	main := layoutShellMain("Hi", "dash", "settings-shell", u)
	if shell.MainClass != "" {
		t.Fatalf("layoutShell MainClass: %q", shell.MainClass)
	}
	if main.MainClass != "settings-shell" {
		t.Fatalf("layoutShellMain MainClass: %q", main.MainClass)
	}
	if shell.Title != main.Title || shell.Active != main.Active || shell.User != main.User {
		t.Fatalf("common fields differ: %#v vs %#v", shell, main)
	}
	if shell.UnreadNotificationCount != 0 || main.UnreadNotificationCount != 0 {
		t.Fatalf("layoutShell helpers expect zero unread count: %d %d", shell.UnreadNotificationCount, main.UnreadNotificationCount)
	}
}
