# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A hands-on Go tutorial for building CRDTs (Conflict-Free Replicated Data Types) from first principles. Students implement each CRDT themselves, culminating in a real-time collaborative todo list with WebSocket sync.

**Module path:** `github.com/zshainsky/crdt-tutorial`
**Go version:** 1.25+
**Single `go.mod` at repo root** — no sub-modules, no Go workspaces.

## Common Commands

```bash
# Run tests for a specific module (student's working code)
cd module-01-gcounter && go test .

# Run tests for a solution
cd solutions/module-01-gcounter && go test .

# Run all tests across the repo
go test ./...
```

## Repository Structure

```
crdt-tutorial/
├── go.mod                          # Single module at root
├── module-XX-name/                 # Student working directory
│   ├── README.md                   # Module instructions
│   ├── *.go                        # Student implementation
│   └── *_test.go                   # Test suite (serves as spec)
├── solutions/
│   └── module-XX-name/             # Reference implementation
│       ├── *_solution.go           # Package: <name>solution
│       └── *_solution_test.go
└── starters/
    └── module-XX-name/             # Starter templates
        ├── *.go.tmpl               # NO package declaration
        └── *_test.go.tmpl
```

### Module Dependencies

Later modules import earlier ones directly:
```go
import gcounter "github.com/zshainsky/crdt-tutorial/module-01-gcounter"
```
Solutions import the working-directory package (not the solution package).

## gopls Compatibility Rules

These rules prevent duplicate declaration errors in VS Code:

1. **Starter templates use `.tmpl` extension with NO package declaration** — gopls ignores syntactically invalid files entirely.

2. **Solution packages use `<name>solution` suffix** (e.g., `package pncountersolution`) — prevents test name collisions with student code.

3. **Single `go.mod`** — never add sub-module go.mod files.

## Creating a New Module

See `.github/CONVENTIONS.md` for the full workflow. Key checklist:

- `module-XX-name/` contains only `README.md` initially (no committed working code)
- `starters/module-XX-name/*.go.tmpl` — NO package declaration at top
- `solutions/module-XX-name/` — package named `<modulename>solution`
- Solution tests pass: `cd solutions/module-XX-name && go test .`
- Module linked from root `README.md`

### Module README Template

Each module's README should include: concept explanation, API spec, starter copy instructions, workflow steps, hints, and link to next module. See `.github/CONVENTIONS.md` for the full template.

### Student Starter Workflow

```bash
# From module-XX-name directory:
cp ../starters/module-XX-name/*.tmpl .
for f in *.tmpl; do mv "$f" "${f%.tmpl}"; done
# Student adds package declaration and implements
go test .
```
