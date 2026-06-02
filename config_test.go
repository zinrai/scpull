package main

import (
	"strings"
	"testing"
)

const sampleConfig = `
profiles:
  switch-logs:
    user: admin
    paths:
      - /var/log/messages
      - /flash/syslog.txt
  ansible-applied:
    user: deploy
    paths:
      - /etc/nginx/conf.d/site.conf
`

func TestProfile(t *testing.T) {
	cfg, err := parseConfig(strings.NewReader(sampleConfig))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p, err := cfg.profile("switch-logs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.User != "admin" {
		t.Errorf("user = %q, want admin", p.User)
	}
	if len(p.Paths) != 2 {
		t.Errorf("paths length = %d, want 2", len(p.Paths))
	}
}

func TestProfileNotFound(t *testing.T) {
	cfg, _ := parseConfig(strings.NewReader(sampleConfig))
	if _, err := cfg.profile("missing"); err == nil {
		t.Fatal("expected error for unknown profile, got nil")
	}
}

func TestProfileEmptyUser(t *testing.T) {
	cfg, _ := parseConfig(strings.NewReader(`
profiles:
  bad:
    paths:
      - /etc/hosts
`))
	if _, err := cfg.profile("bad"); err == nil {
		t.Fatal("expected error for empty user, got nil")
	}
}

func TestProfileEmptyPaths(t *testing.T) {
	cfg, _ := parseConfig(strings.NewReader(`
profiles:
  bad:
    user: admin
`))
	if _, err := cfg.profile("bad"); err == nil {
		t.Fatal("expected error for empty paths, got nil")
	}
}
