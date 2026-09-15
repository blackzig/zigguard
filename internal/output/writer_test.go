package output

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/blackzig/zigguard/internal/compiler"
	"github.com/blackzig/zigguard/internal/managed"
)

func TestWriteArtifactsCreatesManagedSectionFile(t *testing.T) {
	root := t.TempDir()
	artifact := compiler.Artifact{
		Path:      "AGENTS.md",
		Ownership: compiler.ManagedSection,
		Content:   managed.Wrap("policy"),
	}

	if err := WriteArtifacts(root, []compiler.Artifact{artifact}, false); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != artifact.Content {
		t.Fatalf("content = %q, want %q", string(got), artifact.Content)
	}
}

func TestWriteArtifactsPreservesHumanContent(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "AGENTS.md")
	human := []byte("# Human instructions\n\nKeep this.\n")
	if err := os.WriteFile(path, human, 0o644); err != nil {
		t.Fatal(err)
	}

	artifact := compiler.Artifact{
		Path:      "AGENTS.md",
		Ownership: compiler.ManagedSection,
		Content:   managed.Wrap("policy"),
	}
	if err := WriteArtifacts(root, []compiler.Artifact{artifact}, false); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(got, human) {
		t.Fatalf("human content was not preserved: %q", got)
	}
	section, state, err := managed.Extract(got)
	if err != nil {
		t.Fatal(err)
	}
	if state != managed.Present || string(section) != artifact.Content {
		t.Fatalf("section = %q, want %q", section, artifact.Content)
	}
}

func TestWriteArtifactsUpdatesOnlyManagedSection(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "AGENTS.md")
	before := []byte("# Before\n\n")
	after := []byte("\n# After\n")
	old := []byte(managed.Wrap("old policy"))
	current := append(append(append([]byte(nil), before...), old...), after...)
	if err := os.WriteFile(path, current, 0o644); err != nil {
		t.Fatal(err)
	}

	artifact := compiler.Artifact{
		Path:      "AGENTS.md",
		Ownership: compiler.ManagedSection,
		Content:   managed.Wrap("new policy"),
	}
	if err := WriteArtifacts(root, []compiler.Artifact{artifact}, false); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(got, before) || !bytes.HasSuffix(got, after) {
		t.Fatalf("surrounding human content changed: %q", got)
	}
	if !bytes.Contains(got, []byte("new policy")) || bytes.Contains(got, []byte("old policy")) {
		t.Fatalf("managed section was not replaced: %q", got)
	}
}

func TestWriteArtifactsRejectsMalformedSectionEvenWithForce(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "AGENTS.md")
	current := []byte(managed.StartMarker + "\nfirst\n" + managed.StartMarker + "\n" + managed.EndMarker + "\n")
	if err := os.WriteFile(path, current, 0o644); err != nil {
		t.Fatal(err)
	}

	err := WriteArtifacts(root, []compiler.Artifact{{
		Path:      "AGENTS.md",
		Ownership: compiler.ManagedSection,
		Content:   managed.Wrap("new policy"),
	}}, true)
	if !errors.Is(err, managed.ErrMalformed) {
		t.Fatalf("WriteArtifacts() error = %v, want ErrMalformed", err)
	}

	got, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !bytes.Equal(got, current) {
		t.Fatal("malformed file was modified")
	}
}

func TestWriteArtifactsMigratesLegacyManagedFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "AGENTS.md")
	if err := os.WriteFile(path, []byte(managed.LegacyMarker+"\n\nlegacy generated policy\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	artifact := compiler.Artifact{
		Path:      "AGENTS.md",
		Ownership: compiler.ManagedSection,
		Content:   managed.Wrap("new policy"),
	}
	if err := WriteArtifacts(root, []compiler.Artifact{artifact}, false); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != artifact.Content {
		t.Fatalf("legacy migration = %q, want %q", got, artifact.Content)
	}
}

func TestWriteArtifactsUpdatesZigGuardMetadataWithoutForce(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".zigguard")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "manifest.json")
	if err := os.WriteFile(path, []byte("{\"old\":true}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	artifact := compiler.Artifact{
		Path:      ".zigguard/manifest.json",
		Ownership: compiler.ManagedFile,
		Content:   "{\"new\":true}\n",
	}
	if err := WriteArtifacts(root, []compiler.Artifact{artifact}, false); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != artifact.Content {
		t.Fatalf("content = %q, want %q", string(got), artifact.Content)
	}
}
