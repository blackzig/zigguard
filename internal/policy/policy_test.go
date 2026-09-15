package policy

import (
	"strings"
	"testing"
)

func TestParseValidPolicy(t *testing.T) {
	p, err := Parse([]byte(DefaultConfig))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if p.Version != "0.1" {
		t.Fatalf("Version = %q, want 0.1", p.Version)
	}
	if p.Project.Name != "my-project" {
		t.Fatalf("Project.Name = %q, want my-project", p.Project.Name)
	}
	if len(p.Targets) != 4 {
		t.Fatalf("len(Targets) = %d, want 4", len(p.Targets))
	}
}

func TestParseRejectsUnknownField(t *testing.T) {
	input := strings.Replace(DefaultConfig, "project:\n", "project:\n  unexpected: true\n", 1)
	if _, err := Parse([]byte(input)); err == nil {
		t.Fatal("Parse() error = nil, want unknown-field error")
	}
}

func TestParseRejectsUnsupportedTarget(t *testing.T) {
	input := strings.Replace(DefaultConfig, "  - codex\n", "  - unknown-agent\n", 1)
	if _, err := Parse([]byte(input)); err == nil {
		t.Fatal("Parse() error = nil, want unsupported-target error")
	}
}

func TestParseRejectsUnsupportedAction(t *testing.T) {
	input := strings.Replace(DefaultConfig, "refactor_outside_request: block", "refactor_outside_request: maybe", 1)
	if _, err := Parse([]byte(input)); err == nil {
		t.Fatal("Parse() error = nil, want unsupported-action error")
	}
}

func TestSortedTargetsIsDeterministic(t *testing.T) {
	p, err := Parse([]byte(DefaultConfig))
	if err != nil {
		t.Fatal(err)
	}
	got := p.SortedTargets()
	want := []string{"claude", "codex", "copilot", "cursor"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("SortedTargets()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
