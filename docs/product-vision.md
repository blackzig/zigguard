# Product Vision

## One sentence

ZigGuard lets engineering teams define how AI coding agents are allowed to work once, then translates and enforces those rules across supported tools.

## Problem

AI coding agents are increasingly capable, but each product has its own instruction model. A team may end up maintaining several overlapping rule files while still lacking a consistent enforcement boundary.

This creates four problems:

1. duplicated configuration;
2. inconsistent agent behavior;
3. weak auditability;
4. policy drift between repositories and tools.

## Proposed solution

ZigGuard introduces a neutral policy file, `zigguard.yml`, and a compiler that generates tool-specific configuration.

The core focuses on deterministic policy compilation and local checks. Optional commercial services can add organization-wide synchronization, private policy packs, audit history, dashboards, and integrations.

## Initial users

- individual developers using multiple AI coding agents;
- freelancers who need repeatable project safety rules;
- software teams standardizing AI-assisted development;
- maintainers of Java, PHP, JavaScript and legacy systems;
- security-conscious teams that need explicit approval/block rules.

## Positioning

ZigGuard is not primarily an AI code reviewer.

Code review products inspect code after or during creation. ZigGuard's main responsibility is defining and enforcing the policy under which an AI coding agent is expected to operate.

## Non-goals for MVP 0.1

- replacing CI platforms;
- replacing SAST/secret scanners;
- hosting source code;
- building a general-purpose LLM gateway;
- guaranteeing that natural-language-only instructions are enforceable.

When a rule cannot be technically enforced, ZigGuard must describe that limitation explicitly.
