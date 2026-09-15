package compiler

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/blackzig/zigguard/internal/policy"
)

func TestJavaSpringOracleFixture(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	exampleDir := filepath.Join(filepath.Dir(currentFile), "..", "..", "examples", "java-spring-oracle")
	p, err := policy.Load(filepath.Join(exampleDir, "zigguard.yml"))
	if err != nil {
		t.Fatal(err)
	}

	artifacts, err := Compile(p)
	if err != nil {
		t.Fatal(err)
	}

	expectedByArtifact := map[string]string{
		"AGENTS.md":               filepath.Join(exampleDir, "generated", "AGENTS.md"),
		"CLAUDE.md":               filepath.Join(exampleDir, "generated", "CLAUDE.md"),
		".zigguard/manifest.json": filepath.Join(exampleDir, "generated", "manifest.json"),
	}

	if len(artifacts) != len(expectedByArtifact) {
		t.Fatalf("got %d artifacts, want %d", len(artifacts), len(expectedByArtifact))
	}

	for _, artifact := range artifacts {
		expectedPath, ok := expectedByArtifact[artifact.Path]
		if !ok {
			t.Fatalf("unexpected artifact %q", artifact.Path)
		}
		expected, err := os.ReadFile(expectedPath)
		if err != nil {
			t.Fatal(err)
		}
		if artifact.Content != string(expected) {
			t.Fatalf("artifact %s differs from fixture %s", artifact.Path, expectedPath)
		}
	}
}
