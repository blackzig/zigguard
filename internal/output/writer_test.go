package output

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/blackzig/zigguard/internal/compiler"
)

func TestWriteArtifactsCreatesManagedFile(t *testing.T) {
	root := t.TempDir()
	artifact := compiler.Artifact{
		Path:    "AGENTS.md",
		Content: compiler.ManagedMarker + "\npolicy\n",
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

func TestWriteArtifactsRefusesUnmanagedCollision(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "AGENTS.md")
	if err := os.WriteFile(path, []byte("# existing instructions\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := WriteArtifacts(root, []compiler.Artifact{{
		Path:    "AGENTS.md",
		Content: compiler.ManagedMarker + "\npolicy\n",
	}}, false)
	if err == nil {
		t.Fatal("WriteArtifacts() error = nil, want collision error")
	}
	if !strings.Contains(err.Error(), "unmanaged") {
		t.Fatalf("error = %q, want unmanaged collision", err)
	}
}

func TestWriteArtifactsAllowsExplicitForce(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "AGENTS.md")
	if err := os.WriteFile(path, []byte("# existing instructions\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	artifact := compiler.Artifact{
		Path:    "AGENTS.md",
		Content: compiler.ManagedMarker + "\npolicy\n",
	}
	if err := WriteArtifacts(root, []compiler.Artifact{artifact}, true); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != artifact.Content {
		t.Fatalf("content = %q, want forced content", string(got))
	}
}
