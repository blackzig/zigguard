package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/blackzig/zigguard/internal/checker"
)

func TestInitValidateAndCompileDryRun(t *testing.T) {
	root := t.TempDir()
	policyPath := filepath.Join(root, "zigguard.yml")

	var out bytes.Buffer
	var errOut bytes.Buffer

	if code := Run([]string{"init", "--file", policyPath}, &out, &errOut); code != 0 {
		t.Fatalf("init code = %d, stderr = %s", code, errOut.String())
	}
	if _, err := os.Stat(policyPath); err != nil {
		t.Fatal(err)
	}

	out.Reset()
	errOut.Reset()
	if code := Run([]string{"validate", "--file", policyPath}, &out, &errOut); code != 0 {
		t.Fatalf("validate code = %d, stderr = %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "valid:") {
		t.Fatalf("validate output = %q", out.String())
	}

	out.Reset()
	errOut.Reset()
	if code := Run([]string{"compile", "--file", policyPath, "--root", root, "--dry-run"}, &out, &errOut); code != 0 {
		t.Fatalf("compile code = %d, stderr = %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "AGENTS.md") || !strings.Contains(out.String(), "CLAUDE.md") {
		t.Fatalf("compile output = %q", out.String())
	}
}

func TestCompilePreservesExistingHumanAgentsFile(t *testing.T) {
	root := t.TempDir()
	policyPath := filepath.Join(root, "zigguard.yml")
	if err := os.WriteFile(policyPath, []byte(`version: "0.1"
project:
  name: coexistence
targets:
  - codex
rules:
  security:
    secrets: block
`), 0o644); err != nil {
		t.Fatal(err)
	}

	agentsPath := filepath.Join(root, "AGENTS.md")
	human := "# Existing human instructions\n\nKeep this text.\n"
	if err := os.WriteFile(agentsPath, []byte(human), 0o644); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	var errOut bytes.Buffer
	code := Run([]string{"compile", "--file", policyPath, "--root", root}, &out, &errOut)
	if code != 0 {
		t.Fatalf("compile code = %d, stderr = %s", code, errOut.String())
	}

	got, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(got), human) {
		t.Fatalf("human instructions were not preserved: %q", got)
	}
	if !strings.Contains(string(got), "<!-- zigguard:managed:start -->") {
		t.Fatalf("managed section was not appended: %q", got)
	}
}

func TestCheckPassesAfterCompile(t *testing.T) {
	root := t.TempDir()
	policyPath := filepath.Join(root, "zigguard.yml")

	var out bytes.Buffer
	var errOut bytes.Buffer
	if code := Run([]string{"init", "--file", policyPath}, &out, &errOut); code != 0 {
		t.Fatalf("init code = %d, stderr = %s", code, errOut.String())
	}

	out.Reset()
	errOut.Reset()
	if code := Run([]string{"compile", "--file", policyPath, "--root", root}, &out, &errOut); code != 0 {
		t.Fatalf("compile code = %d, stderr = %s", code, errOut.String())
	}

	out.Reset()
	errOut.Reset()
	if code := Run([]string{"check", "--file", policyPath, "--root", root}, &out, &errOut); code != 0 {
		t.Fatalf("check code = %d, stdout = %s, stderr = %s", code, out.String(), errOut.String())
	}
	if !strings.Contains(out.String(), "check: PASS") {
		t.Fatalf("check output = %q", out.String())
	}
}

func TestCheckJSONReportsManagedSectionDrift(t *testing.T) {
	root := t.TempDir()
	policyPath := filepath.Join(root, "zigguard.yml")

	var out bytes.Buffer
	var errOut bytes.Buffer
	if code := Run([]string{"init", "--file", policyPath}, &out, &errOut); code != 0 {
		t.Fatalf("init code = %d, stderr = %s", code, errOut.String())
	}
	out.Reset()
	errOut.Reset()
	if code := Run([]string{"compile", "--file", policyPath, "--root", root}, &out, &errOut); code != 0 {
		t.Fatalf("compile code = %d, stderr = %s", code, errOut.String())
	}

	path := filepath.Join(root, "AGENTS.md")
	current, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	changed := bytes.Replace(current,
		[]byte("Do not refactor code outside the user's requested scope."),
		[]byte("Changed managed policy text."),
		1,
	)
	if bytes.Equal(changed, current) {
		t.Fatal("test setup failed to modify managed policy")
	}
	if err := os.WriteFile(path, changed, 0o644); err != nil {
		t.Fatal(err)
	}

	out.Reset()
	errOut.Reset()
	code := Run([]string{"check", "--file", policyPath, "--root", root, "--format", "json"}, &out, &errOut)
	if code != 1 {
		t.Fatalf("check code = %d, want 1; stderr = %s", code, errOut.String())
	}

	var report checker.Report
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.OK {
		t.Fatal("report.OK = true, want false")
	}
	found := false
	for _, violation := range report.Violations {
		if violation.ID == checker.DriftedArtifact && violation.Path == "AGENTS.md" {
			found = true
		}
	}
	if !found {
		t.Fatalf("violations = %#v, want drift for AGENTS.md", report.Violations)
	}
}
