package output

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/blackzig/zigguard/internal/compiler"
	"github.com/blackzig/zigguard/internal/managed"
)

func WriteArtifacts(root string, artifacts []compiler.Artifact, force bool) error {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolve output root: %w", err)
	}

	for _, artifact := range artifacts {
		if err := writeArtifact(rootAbs, artifact, force); err != nil {
			return err
		}
	}
	return nil
}

func writeArtifact(rootAbs string, artifact compiler.Artifact, force bool) error {
	relative := filepath.FromSlash(artifact.Path)
	if filepath.IsAbs(relative) {
		return fmt.Errorf("refusing absolute artifact path %q", artifact.Path)
	}

	target := filepath.Join(rootAbs, relative)
	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return fmt.Errorf("resolve artifact %q: %w", artifact.Path, err)
	}

	rel, err := filepath.Rel(rootAbs, targetAbs)
	if err != nil {
		return fmt.Errorf("validate artifact path %q: %w", artifact.Path, err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("artifact path escapes output root: %q", artifact.Path)
	}

	content := []byte(artifact.Content)
	if current, err := os.ReadFile(targetAbs); err == nil {
		next, err := mergeExisting(artifact, current, force)
		if err != nil {
			return fmt.Errorf("prepare artifact %q: %w", artifact.Path, err)
		}
		content = next
		if bytes.Equal(current, content) {
			return nil
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("read existing artifact %q: %w", artifact.Path, err)
	}

	if err := os.MkdirAll(filepath.Dir(targetAbs), 0o755); err != nil {
		return fmt.Errorf("create artifact directory for %q: %w", artifact.Path, err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(targetAbs), ".zigguard-*")
	if err != nil {
		return fmt.Errorf("create temporary artifact for %q: %w", artifact.Path, err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if err := tmp.Chmod(0o644); err != nil {
		tmp.Close()
		return fmt.Errorf("set artifact permissions for %q: %w", artifact.Path, err)
	}
	if _, err := tmp.Write(content); err != nil {
		tmp.Close()
		return fmt.Errorf("write artifact %q: %w", artifact.Path, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close artifact %q: %w", artifact.Path, err)
	}
	if err := os.Rename(tmpName, targetAbs); err != nil {
		return fmt.Errorf("replace artifact %q: %w", artifact.Path, err)
	}
	return nil
}

func mergeExisting(artifact compiler.Artifact, current []byte, force bool) ([]byte, error) {
	switch ownership(artifact) {
	case compiler.ManagedSection:
		return managed.Merge(current, []byte(artifact.Content))
	case compiler.ManagedFile:
		if bytes.Equal(current, []byte(artifact.Content)) {
			return current, nil
		}
		if force || isZigGuardOwnedFile(artifact.Path, current) {
			return []byte(artifact.Content), nil
		}
		return nil, fmt.Errorf("refusing to overwrite unmanaged file; rerun with --force only if full replacement is intended")
	default:
		return nil, fmt.Errorf("unsupported ownership mode %q", artifact.Ownership)
	}
}

func ownership(artifact compiler.Artifact) compiler.Ownership {
	if artifact.Ownership == "" {
		return compiler.ManagedFile
	}
	return artifact.Ownership
}

func isZigGuardOwnedFile(path string, current []byte) bool {
	normalized := filepath.ToSlash(filepath.Clean(path))
	if strings.HasPrefix(normalized, ".zigguard/") {
		return true
	}
	return bytes.HasPrefix(current, []byte(managed.LegacyMarker))
}
