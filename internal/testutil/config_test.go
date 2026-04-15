package testutil

import (
	"testing"
)

func TestDefaultTestConfig_repoURLMatchesDefaultRepoConstant(t *testing.T) {
	t.Parallel()
	c := DefaultTestConfig()
	if c.RepoURL != DefaultTestRepoURL {
		t.Fatalf("RepoURL %q want %q", c.RepoURL, DefaultTestRepoURL)
	}
}
