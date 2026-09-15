package checker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/blackzig/zigguard/internal/compiler"
	"github.com/blackzig/zigguard/internal/output"
	"github.com/blackzig/zigguard/internal/policy"
)

func TestCheckPassesWhenArtifactsMatch(t *testing.T) {
	root := t.TempDir()
	p, err := policy.Parse([]byte(policy.DefaultConfig))
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := compiler.Compile(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := output.WriteArtifacts(root, artifacts, false); err != nil {
		t.Fatal(err)
	}

	report, err := Check(root, p)
	if err != nil {
		t.Fatal(err)
	}
	if !report.OK {
		t.Fatalf("report.OK = false, violations = %#v", report.Violations)
	}
	if report.CheckedArtifacts != len(artifacts) {
		t.Fatalf("CheckedArtifacts = %d, want %d", report.CheckedArtifacts, len(artifacts))
	}
}

func TestCheckReportsMissingArtifact(t *testing.T) {
	root := t.TempDir()
	p, err := policy.Parse([]byte(policy.DefaultConfig))
	if err != nil {
		t.Fatal(err)
	}

	report, err := Check(root, p)
	if err != nil {
		t.Fatal(err)
	}
	if report.OK {
		t.Fatal("report.OK = true, want false")
	}
	if !containsViolation(report, MissingArtifact, "AGENTS.md") {
		t.Fatalf("violations = %#v, want ZG001 for AGENTS.md", report.Violations)
	}
}

func TestCheckReportsManagedDrift(t *testing.T) {
	root := t.TempDir()
	p, err := policy.Parse([]byte(policy.DefaultConfig))
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := compiler.Compile(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := output.WriteArtifacts(root, artifacts, false); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(root, "AGENTS.md")
	current, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(current, []byte("\nmanual drift\n")...), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Check(root, p)
	if err != nil {
		t.Fatal(err)
	}
	if !containsViolation(report, DriftedArtifact, "AGENTS.md") {
		t.Fatalf("violations = %#v, want ZG002 for AGENTS.md", report.Violations)
	}
}

func TestCheckReportsUnmanagedConflict(t *testing.T) {
	root := t.TempDir()
	p, err := policy.Parse([]byte(policy.DefaultConfig))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("# existing human instructions\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Check(root, p)
	if err != nil {
		t.Fatal(err)
	}
	if !containsViolation(report, UnmanagedConflict, "AGENTS.md") {
		t.Fatalf("violations = %#v, want ZG003 for AGENTS.md", report.Violations)
	}
}

func containsViolation(report Report, id, path string) bool {
	for _, violation := range report.Violations {
		if violation.ID == id && violation.Path == path {
			return true
		}
	}
	return false
}
