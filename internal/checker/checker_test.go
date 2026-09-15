package checker

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/blackzig/zigguard/internal/compiler"
	"github.com/blackzig/zigguard/internal/managed"
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

func TestCheckPassesWithHumanContentAroundManagedSection(t *testing.T) {
	root := t.TempDir()
	p, err := policy.Parse([]byte(policy.DefaultConfig))
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := compiler.Compile(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("# Human instructions\n"), 0o644); err != nil {
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
	if err := os.WriteFile(path, append(current, []byte("\n# Human suffix\n")...), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Check(root, p)
	if err != nil {
		t.Fatal(err)
	}
	if !report.OK {
		t.Fatalf("human content outside managed section caused violations: %#v", report.Violations)
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

func TestCheckReportsMissingSectionInHumanFile(t *testing.T) {
	root := t.TempDir()
	p, err := policy.Parse([]byte(policy.DefaultConfig))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("# Human instructions\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Check(root, p)
	if err != nil {
		t.Fatal(err)
	}
	if !containsViolation(report, MissingArtifact, "AGENTS.md") {
		t.Fatalf("violations = %#v, want ZG001 for missing managed section", report.Violations)
	}
}

func TestCheckReportsManagedSectionDrift(t *testing.T) {
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
	current = bytes.Replace(current, []byte("Do not refactor code outside"), []byte("CHANGED refactor rule outside"), 1)
	if err := os.WriteFile(path, current, 0o644); err != nil {
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

func TestCheckReportsMalformedManagedMarkers(t *testing.T) {
	root := t.TempDir()
	p, err := policy.Parse([]byte(policy.DefaultConfig))
	if err != nil {
		t.Fatal(err)
	}
	content := []byte(managed.StartMarker + "\nfirst\n" + managed.StartMarker + "\n" + managed.EndMarker + "\n")
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), content, 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Check(root, p)
	if err != nil {
		t.Fatal(err)
	}
	if !containsViolation(report, UnmanagedConflict, "AGENTS.md") {
		t.Fatalf("violations = %#v, want ZG003 for malformed section", report.Violations)
	}
}

func TestCheckReportsLegacyManagedFileAsDrift(t *testing.T) {
	root := t.TempDir()
	p, err := policy.Parse([]byte(policy.DefaultConfig))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte(managed.LegacyMarker+"\nlegacy\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Check(root, p)
	if err != nil {
		t.Fatal(err)
	}
	if !containsViolation(report, DriftedArtifact, "AGENTS.md") {
		t.Fatalf("violations = %#v, want ZG002 for legacy managed file", report.Violations)
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
