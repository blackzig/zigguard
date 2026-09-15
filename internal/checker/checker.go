package checker

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/blackzig/zigguard/internal/compiler"
	"github.com/blackzig/zigguard/internal/managed"
	"github.com/blackzig/zigguard/internal/policy"
)

const (
	MissingArtifact   = "ZG001"
	DriftedArtifact   = "ZG002"
	UnmanagedConflict = "ZG003"
)

type Violation struct {
	ID      string `json:"id"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

type Report struct {
	Project          string      `json:"project"`
	PolicyVersion    string      `json:"policy_version"`
	OK               bool        `json:"ok"`
	CheckedArtifacts int         `json:"checked_artifacts"`
	Violations       []Violation `json:"violations"`
}

func Check(root string, p policy.Policy) (Report, error) {
	artifacts, err := compiler.Compile(p)
	if err != nil {
		return Report{}, err
	}

	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return Report{}, fmt.Errorf("resolve check root: %w", err)
	}

	report := Report{
		Project:          p.Project.Name,
		PolicyVersion:    p.Version,
		OK:               true,
		CheckedArtifacts: len(artifacts),
		Violations:       []Violation{},
	}

	for _, artifact := range artifacts {
		target, err := artifactPath(rootAbs, artifact.Path)
		if err != nil {
			return Report{}, err
		}

		current, err := os.ReadFile(target)
		if os.IsNotExist(err) {
			report.Violations = append(report.Violations, Violation{
				ID:      MissingArtifact,
				Path:    artifact.Path,
				Message: "generated artifact is missing",
			})
			continue
		}
		if err != nil {
			return Report{}, fmt.Errorf("read artifact %q: %w", artifact.Path, err)
		}

		violation := checkArtifact(artifact, current)
		if violation != nil {
			report.Violations = append(report.Violations, *violation)
		}
	}

	sort.Slice(report.Violations, func(i, j int) bool {
		if report.Violations[i].Path == report.Violations[j].Path {
			return report.Violations[i].ID < report.Violations[j].ID
		}
		return report.Violations[i].Path < report.Violations[j].Path
	})
	report.OK = len(report.Violations) == 0
	return report, nil
}

func checkArtifact(artifact compiler.Artifact, current []byte) *Violation {
	if artifact.Ownership == compiler.ManagedSection {
		section, state, err := managed.Extract(current)
		if err != nil {
			return &Violation{
				ID:      UnmanagedConflict,
				Path:    artifact.Path,
				Message: "ZigGuard managed section markers are malformed or ambiguous",
			}
		}
		if state == managed.Missing {
			if bytes.HasPrefix(current, []byte(managed.LegacyMarker)) {
				return &Violation{
					ID:      DriftedArtifact,
					Path:    artifact.Path,
					Message: "legacy ZigGuard-managed file must be recompiled into a managed section",
				}
			}
			if bytes.Contains(current, []byte(managed.LegacyMarker)) {
				return &Violation{
					ID:      UnmanagedConflict,
					Path:    artifact.Path,
					Message: "legacy ZigGuard marker is in an ambiguous location",
				}
			}
			return &Violation{
				ID:      MissingArtifact,
				Path:    artifact.Path,
				Message: "ZigGuard managed section is missing",
			}
		}
		if !bytes.Equal(section, []byte(artifact.Content)) {
			return &Violation{
				ID:      DriftedArtifact,
				Path:    artifact.Path,
				Message: "ZigGuard managed section differs from the current policy",
			}
		}
		return nil
	}

	if bytes.Equal(current, []byte(artifact.Content)) {
		return nil
	}
	if isZigGuardOwnedFile(artifact.Path, current) {
		return &Violation{
			ID:      DriftedArtifact,
			Path:    artifact.Path,
			Message: "generated artifact differs from the current policy",
		}
	}
	return &Violation{
		ID:      UnmanagedConflict,
		Path:    artifact.Path,
		Message: "expected fully managed artifact path is occupied by an unmanaged file",
	}
}

func artifactPath(rootAbs, artifactPath string) (string, error) {
	relative := filepath.FromSlash(artifactPath)
	if filepath.IsAbs(relative) {
		return "", fmt.Errorf("refusing absolute artifact path %q", artifactPath)
	}

	targetAbs, err := filepath.Abs(filepath.Join(rootAbs, relative))
	if err != nil {
		return "", fmt.Errorf("resolve artifact %q: %w", artifactPath, err)
	}

	rel, err := filepath.Rel(rootAbs, targetAbs)
	if err != nil {
		return "", fmt.Errorf("validate artifact path %q: %w", artifactPath, err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("artifact path escapes check root: %q", artifactPath)
	}
	return targetAbs, nil
}

func isZigGuardOwnedFile(path string, current []byte) bool {
	normalized := filepath.ToSlash(filepath.Clean(path))
	if strings.HasPrefix(normalized, ".zigguard/") {
		return true
	}
	return bytes.HasPrefix(current, []byte(managed.LegacyMarker))
}
