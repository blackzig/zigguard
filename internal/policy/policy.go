package policy

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const Version = "0.1"

type Policy struct {
	Version string   `yaml:"version"`
	Project Project  `yaml:"project"`
	Targets []string `yaml:"targets"`
	Rules   Rules    `yaml:"rules"`
}

type Project struct {
	Name  string   `yaml:"name"`
	Stack []string `yaml:"stack"`
}

type Rules struct {
	Scope    ScopeRules    `yaml:"scope"`
	Security SecurityRules `yaml:"security"`
	Quality  QualityRules  `yaml:"quality"`
	Git      GitRules      `yaml:"git"`
}

type ScopeRules struct {
	RefactorOutsideRequest string `yaml:"refactor_outside_request"`
	PublicAPIChanges       string `yaml:"public_api_changes"`
}

type SecurityRules struct {
	Secrets    string `yaml:"secrets"`
	DynamicSQL string `yaml:"dynamic_sql"`
}

type QualityRules struct {
	TestsRequired bool `yaml:"tests_required"`
	LintRequired  bool `yaml:"lint_required"`
}

type GitRules struct {
	ForcePushProtectedBranches string `yaml:"force_push_protected_branches"`
}

var supportedTargets = map[string]struct{}{
	"claude":  {},
	"codex":   {},
	"cursor":  {},
	"copilot": {},
}

var supportedActions = map[string]struct{}{
	"allow":            {},
	"warn":             {},
	"require_approval": {},
	"block":            {},
}

const DefaultConfig = `version: "0.1"

project:
  name: my-project
  stack: []

targets:
  - claude
  - codex
  - cursor
  - copilot

rules:
  scope:
    refactor_outside_request: block
    public_api_changes: require_approval

  security:
    secrets: block
    dynamic_sql: block

  quality:
    tests_required: true
    lint_required: true

  git:
    force_push_protected_branches: block
`

func Load(path string) (Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Policy{}, fmt.Errorf("read policy: %w", err)
	}
	return Parse(data)
}

func Parse(data []byte) (Policy, error) {
	var p Policy
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	if err := decoder.Decode(&p); err != nil {
		return Policy{}, fmt.Errorf("parse policy: %w", err)
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return Policy{}, fmt.Errorf("parse policy: multiple YAML documents are not supported")
		}
		return Policy{}, fmt.Errorf("parse policy: %w", err)
	}

	if err := p.Validate(); err != nil {
		return Policy{}, err
	}
	return p, nil
}

func (p Policy) Validate() error {
	if p.Version != Version {
		return fmt.Errorf("validate policy: unsupported version %q; expected %q", p.Version, Version)
	}
	if strings.TrimSpace(p.Project.Name) == "" {
		return fmt.Errorf("validate policy: project.name is required")
	}
	if len(p.Targets) == 0 {
		return fmt.Errorf("validate policy: at least one target is required")
	}

	seenTargets := map[string]struct{}{}
	for _, target := range p.Targets {
		if _, ok := supportedTargets[target]; !ok {
			return fmt.Errorf("validate policy: unsupported target %q", target)
		}
		if _, exists := seenTargets[target]; exists {
			return fmt.Errorf("validate policy: duplicate target %q", target)
		}
		seenTargets[target] = struct{}{}
	}

	seenStack := map[string]struct{}{}
	for _, item := range p.Project.Stack {
		item = strings.TrimSpace(item)
		if item == "" {
			return fmt.Errorf("validate policy: project.stack cannot contain empty values")
		}
		if _, exists := seenStack[item]; exists {
			return fmt.Errorf("validate policy: duplicate project.stack value %q", item)
		}
		seenStack[item] = struct{}{}
	}

	actions := []struct {
		name  string
		value string
	}{
		{"rules.scope.refactor_outside_request", p.Rules.Scope.RefactorOutsideRequest},
		{"rules.scope.public_api_changes", p.Rules.Scope.PublicAPIChanges},
		{"rules.security.secrets", p.Rules.Security.Secrets},
		{"rules.security.dynamic_sql", p.Rules.Security.DynamicSQL},
		{"rules.git.force_push_protected_branches", p.Rules.Git.ForcePushProtectedBranches},
	}

	configured := p.Rules.Quality.TestsRequired || p.Rules.Quality.LintRequired
	for _, action := range actions {
		if action.value == "" {
			continue
		}
		configured = true
		if _, ok := supportedActions[action.value]; !ok {
			return fmt.Errorf("validate policy: unsupported action %q for %s", action.value, action.name)
		}
	}
	if !configured {
		return fmt.Errorf("validate policy: configure at least one rule")
	}

	return nil
}

func (p Policy) SortedTargets() []string {
	targets := append([]string(nil), p.Targets...)
	sort.Strings(targets)
	return targets
}
