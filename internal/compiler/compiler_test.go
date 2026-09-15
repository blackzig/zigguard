package compiler

import (
	"reflect"
	"strings"
	"testing"

	"github.com/blackzig/zigguard/internal/managed"
	"github.com/blackzig/zigguard/internal/policy"
)

func TestCompileUsesSharedAgentsSurface(t *testing.T) {
	p, err := policy.Parse([]byte(policy.DefaultConfig))
	if err != nil {
		t.Fatal(err)
	}

	artifacts, err := Compile(p)
	if err != nil {
		t.Fatal(err)
	}

	paths := artifactPaths(artifacts)
	want := []string{".zigguard/manifest.json", "AGENTS.md", "CLAUDE.md"}
	if !reflect.DeepEqual(paths, want) {
		t.Fatalf("paths = %#v, want %#v", paths, want)
	}

	agents := findArtifact(t, artifacts, "AGENTS.md")
	if agents.Ownership != ManagedSection {
		t.Fatalf("AGENTS.md ownership = %q, want %q", agents.Ownership, ManagedSection)
	}
	if !strings.Contains(agents.Content, "**BLOCK:**") {
		t.Fatal("AGENTS.md does not contain BLOCK semantics")
	}
	if !strings.Contains(agents.Content, managed.StartMarker) || !strings.Contains(agents.Content, managed.EndMarker) {
		t.Fatal("AGENTS.md does not contain managed section markers")
	}

	claude := findArtifact(t, artifacts, "CLAUDE.md")
	if claude.Ownership != ManagedSection {
		t.Fatalf("CLAUDE.md ownership = %q, want %q", claude.Ownership, ManagedSection)
	}
	if !strings.Contains(claude.Content, "@AGENTS.md") {
		t.Fatal("CLAUDE.md should import AGENTS.md when shared targets are selected")
	}

	manifest := findArtifact(t, artifacts, ".zigguard/manifest.json")
	if manifest.Ownership != ManagedFile {
		t.Fatalf("manifest ownership = %q, want %q", manifest.Ownership, ManagedFile)
	}
}

func TestCompileClaudeOnlyEmbedsFullPolicy(t *testing.T) {
	input := strings.Replace(policy.DefaultConfig,
		"  - claude\n  - codex\n  - cursor\n  - copilot\n",
		"  - claude\n", 1)

	p, err := policy.Parse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := Compile(p)
	if err != nil {
		t.Fatal(err)
	}

	paths := artifactPaths(artifacts)
	want := []string{".zigguard/manifest.json", "CLAUDE.md"}
	if !reflect.DeepEqual(paths, want) {
		t.Fatalf("paths = %#v, want %#v", paths, want)
	}

	claude := findArtifact(t, artifacts, "CLAUDE.md")
	if strings.Contains(claude.Content, "@AGENTS.md") {
		t.Fatal("Claude-only output must not import missing AGENTS.md")
	}
	if !strings.Contains(claude.Content, "Governance rules") {
		t.Fatal("Claude-only output does not contain full policy")
	}
	if claude.Ownership != ManagedSection {
		t.Fatalf("Claude-only ownership = %q, want %q", claude.Ownership, ManagedSection)
	}
}

func TestCompileIsDeterministic(t *testing.T) {
	p, err := policy.Parse([]byte(policy.DefaultConfig))
	if err != nil {
		t.Fatal(err)
	}

	first, err := Compile(p)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Compile(p)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatal("Compile() is not deterministic")
	}
}

func artifactPaths(artifacts []Artifact) []string {
	paths := make([]string, len(artifacts))
	for i, artifact := range artifacts {
		paths[i] = artifact.Path
	}
	return paths
}

func findArtifact(t *testing.T, artifacts []Artifact, path string) Artifact {
	t.Helper()
	for _, artifact := range artifacts {
		if artifact.Path == path {
			return artifact
		}
	}
	t.Fatalf("artifact %q not found", path)
	return Artifact{}
}
