package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/blackzig/zigguard/internal/checker"
	"github.com/blackzig/zigguard/internal/compiler"
	"github.com/blackzig/zigguard/internal/output"
	"github.com/blackzig/zigguard/internal/policy"
)

var Version = "0.1.0-dev"

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	switch args[0] {
	case "init":
		return runInit(args[1:], stdout, stderr)
	case "validate":
		return runValidate(args[1:], stdout, stderr)
	case "compile":
		return runCompile(args[1:], stdout, stderr)
	case "check":
		return runCheck(args[1:], stdout, stderr)
	case "version", "--version", "-version":
		fmt.Fprintf(stdout, "zigguard %s\n", Version)
		return 0
	case "help", "--help", "-h":
		printUsage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n", args[0])
		printUsage(stderr)
		return 2
	}
}

func runInit(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(stderr)
	file := fs.String("file", "zigguard.yml", "policy file to create")
	force := fs.Bool("force", false, "overwrite an existing policy file")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(stderr, "init does not accept positional arguments")
		return 2
	}

	if !*force {
		if _, err := os.Stat(*file); err == nil {
			fmt.Fprintf(stderr, "refusing to overwrite existing %s; use --force to replace it\n", *file)
			return 1
		} else if !os.IsNotExist(err) {
			fmt.Fprintf(stderr, "inspect %s: %v\n", *file, err)
			return 1
		}
	}

	dir := filepath.Dir(*file)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintf(stderr, "create policy directory: %v\n", err)
			return 1
		}
	}
	if err := os.WriteFile(*file, []byte(policy.DefaultConfig), 0o644); err != nil {
		fmt.Fprintf(stderr, "write %s: %v\n", *file, err)
		return 1
	}

	fmt.Fprintf(stdout, "created %s\n", *file)
	fmt.Fprintln(stdout, "next: edit the project name, stack, targets, and rules, then run zigguard validate")
	return 0
}

func runValidate(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	file := fs.String("file", "zigguard.yml", "policy file to validate")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(stderr, "validate does not accept positional arguments")
		return 2
	}

	p, err := policy.Load(*file)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	fmt.Fprintf(stdout, "valid: %s (policy %s, project %s, %d target(s))\n",
		*file, p.Version, p.Project.Name, len(p.Targets))
	return 0
}

func runCompile(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("compile", flag.ContinueOnError)
	fs.SetOutput(stderr)
	file := fs.String("file", "zigguard.yml", "policy file to compile")
	root := fs.String("root", ".", "repository root where generated files are written")
	force := fs.Bool("force", false, "overwrite unmanaged target files")
	dryRun := fs.Bool("dry-run", false, "show generated paths without writing files")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(stderr, "compile does not accept positional arguments")
		return 2
	}

	p, err := policy.Load(*file)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	artifacts, err := compiler.Compile(p)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if *dryRun {
		fmt.Fprintf(stdout, "compile plan for %s:\n", p.Project.Name)
		for _, artifact := range artifacts {
			fmt.Fprintf(stdout, "  %s [%s]\n", artifact.Path, artifact.Capability)
		}
		fmt.Fprintln(stdout, "note: instruction-context surfaces guide agents; they are not hard security enforcement")
		return 0
	}

	if err := output.WriteArtifacts(*root, artifacts, *force); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	fmt.Fprintf(stdout, "compiled %d artifact(s) for %s:\n", len(artifacts), p.Project.Name)
	for _, artifact := range artifacts {
		fmt.Fprintf(stdout, "  %s [%s]\n", artifact.Path, artifact.Capability)
	}
	fmt.Fprintln(stdout, "note: instruction-context surfaces guide agents; they are not hard security enforcement")
	return 0
}

func runCheck(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	fs.SetOutput(stderr)
	file := fs.String("file", "zigguard.yml", "policy file to check")
	root := fs.String("root", ".", "repository root whose generated artifacts are checked")
	format := fs.String("format", "text", "output format: text or json")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(stderr, "check does not accept positional arguments")
		return 2
	}
	if *format != "text" && *format != "json" {
		fmt.Fprintf(stderr, "unsupported check format %q; expected text or json\n", *format)
		return 2
	}

	p, err := policy.Load(*file)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	report, err := checker.Check(*root, p)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if *format == "json" {
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(report); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	} else {
		if report.OK {
			fmt.Fprintf(stdout, "check: PASS (%s, %d artifact(s) in sync)\n",
				report.Project, report.CheckedArtifacts)
		} else {
			fmt.Fprintf(stdout, "check: FAIL (%d violation(s))\n", len(report.Violations))
			for _, violation := range report.Violations {
				fmt.Fprintf(stdout, "  %s %s: %s\n",
					violation.ID, violation.Path, violation.Message)
			}
		}
	}

	if !report.OK {
		return 1
	}
	return 0
}

func printUsage(w io.Writer) {
	fmt.Fprint(w, `ZigGuard - universal governance layer for AI coding agents

Usage:
  zigguard init [--file zigguard.yml] [--force]
  zigguard validate [--file zigguard.yml]
  zigguard compile [--file zigguard.yml] [--root .] [--dry-run] [--force]
  zigguard check [--file zigguard.yml] [--root .] [--format text|json]
  zigguard version

Commands:
  init      create a starter ZigGuard policy
  validate  parse and validate a policy
  compile   generate managed agent instruction surfaces and a manifest
  check     verify generated artifacts are present and in sync with policy
  version   print the CLI version

Safety:
  compile refuses to overwrite user-managed instruction files unless --force is explicit.
  check is read-only and does not execute project scripts.
`)
}
