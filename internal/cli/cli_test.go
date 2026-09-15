package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestCompileRefusesUnmanagedAgentsFile(t *testing.T) {
	root := t.TempDir()
	policyPath := filepath.Join(root, "zigguard.yml")
	if err := os.WriteFile(policyPath, []byte(`version: "0.1"
project:
  name: collision
targets:
  - codex
rules:
  security:
    secrets: block
`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("# existing\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	var errOut bytes.Buffer
	code := Run([]string{"compile", "--file", policyPath, "--root", root}, &out, &errOut)
	if code == 0 {
		t.Fatal("compile code = 0, want failure")
	}
	if !strings.Contains(errOut.String(), "unmanaged") {
		t.Fatalf("stderr = %q, want unmanaged collision", errOut.String())
	}
}
