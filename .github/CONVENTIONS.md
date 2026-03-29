# CRDT Tutorial Conventions

This document captures the architectural decisions and patterns established for creating tutorial modules in this project.

## Directory Structure

```
crdt-tutorial/
├── go.mod                      # Single Go module at root
├── module-XX-name/             # Working directory (student implementation)
│   ├── README.md              # Module instructions & learning objectives
│   ├── *.go                   # Student's working code
│   └── *_test.go              # Test suite (serves as specification)
├── solutions/
│   └── module-XX-name/        # Reference implementation
│       ├── *_solution.go      # Reference code
│       └── *_solution_test.go # Reference tests
└── starters/
    └── module-XX-name/        # Template files with TODO markers
        ├── *.go.tmpl          # Starter code template (NO package declaration)
        └── *_test.go.tmpl     # Starter test template (NO package declaration)
```

## Critical Rules for gopls Compatibility

### 1. Starter Templates: NO Package Declarations

**MUST DO:**
```go
// starters/module-XX-name/pncounter.go.tmpl
// TODO: Add package declaration here (package pncounter)

type PNCounter struct {
    // ...
}
```

**NEVER DO:**
```go
// starters/module-XX-name/pncounter.go.tmpl
package pncounter  // ❌ This causes gopls to parse the file!

type PNCounter struct {
    // ...
}
```

**Why:** gopls only parses valid Go files. Without a package declaration, the .tmpl file is invalid Go syntax and gopls ignores it completely. This prevents "duplicate declaration" errors.

### 2. Solutions: Different Package Name

**Pattern:** `package <modulename>solution`

**Example:**
- Module: `module-02-pncounter`
- Working code: `package pncounter` 
- Solution code: `package pncountersolution`

**Files:**
```go
// solutions/module-02-pncounter/pncounter_solution.go
package pncountersolution

// solutions/module-02-pncounter/pncounter_solution_test.go
package pncountersolution
```

**Why:** Different package names prevent gopls from flagging duplicate test names like `TestPNCounterIncrement` appearing in both the working directory and solutions directory.

### 3. Single go.mod at Repository Root

**MUST DO:**
- One `go.mod` at `/Users/zshainky/Projects/crdt-tutorial/go.mod`
- Module path: `github.com/zshainsky/crdt-tutorial`

**NEVER DO:**
- Multiple `go.mod` files in subdirectories
- Separate modules for solutions/ or starters/

**Why:** Single module simplifies dependency management and is more user-friendly. Students don't need to understand Go workspaces or multi-module repositories.

## File Naming Patterns

### Starter Templates
- **Extension:** `.tmpl` (not `.go`)
- **Naming:** `<name>.go.tmpl` and `<name>_test.go.tmpl`
- **Examples:**
  - `pncounter.go.tmpl`
  - `pncounter_test.go.tmpl`

### Solution Files
- **Suffix:** `_solution` before extension
- **Package:** `<modulename>solution`
- **Examples:**
  - `pncounter_solution.go`
  - `pncounter_solution_test.go`

### Working Files (Student Code)
- **Standard Go naming:** `<name>.go` and `<name>_test.go`
- **Package:** `<modulename>`
- **Examples:**
  - `pncounter.go`
  - `pncounter_test.go`

## Module Creation Workflow

### Step 1: Create Working Directory
```bash
mkdir module-02-pncounter
```

Create `module-02-pncounter/README.md` with:
- Learning objectives
- Concept explanation
- API specification
- Step-by-step workflow
- Hints section
- Link to next module

### Step 2: Create Starter Templates
```bash
mkdir -p starters/module-02-pncounter
```

Create `.tmpl` files with:
- **NO** package declaration (gopls must ignore)
- TODO comments guiding implementation
- Function stubs with comments
- Test function signatures with TODO bodies

Example structure:
```go
// starters/module-02-pncounter/pncounter.go.tmpl
// TODO: Add package declaration here (package pncounter)

// PNCounter is a positive-negative counter that can increment and decrement.
type PNCounter struct {
    // TODO: Implement using two G-Counters
}

// NewPNCounter creates a new PN-Counter for the given replica.
func NewPNCounter(replicaID string) *PNCounter {
    // TODO: Initialize both P and N counters
    return nil
}
```

### Step 3: Create Solution Implementation
```bash
mkdir -p solutions/module-02-pncounter
```

Create solution files with:
- Package name: `package <modulename>solution`
- Full working implementation
- Complete test suite
- All tests must pass

Example:
```go
// solutions/module-02-pncounter/pncounter_solution.go
package pncountersolution

import "github.com/zshainsky/crdt-tutorial/module-01-gcounter"

type PNCounter struct {
    replicaID string
    p         *gcounter.GCounter
    n         *gcounter.GCounter
}
```

### Step 4: Verify No gopls Conflicts
```bash
# All three packages should be distinct in gopls:
# - github.com/zshainsky/crdt-tutorial/module-02-pncounter
# - github.com/zshainsky/crdt-tutorial/solutions/module-02-pncounter
# - starters/ .tmpl files are ignored (no package declaration)

# Test solution compiles and passes
cd solutions/module-02-pncounter
go test .

# Verify starter templates are syntactically incomplete (expected)
# They should NOT compile due to missing package declaration
```

### Step 5: Test End-to-End Workflow
Simulate student experience:
```bash
cd module-02-pncounter

# Copy starters
cp ../starters/module-02-pncounter/*.tmpl .

# Rename .tmpl to .go
for f in *.tmpl; do mv "$f" "${f%.tmpl}"; done

# Student adds package declaration and implements
# Verify tests can run
go test .
```

## What We Tried That DIDN'T Work

### ❌ Nested solutions/ and starters/ Inside Module Directory
```
module-01-gcounter/
├── gcounter.go
├── solution/
│   └── gcounter_solution.go  # Same package as parent → gopls conflicts
└── starter/
    └── gcounter_starter.go   # Same package as parent → gopls conflicts
```
**Problem:** All use `package gcounter`, gopls sees duplicate declarations.

### ❌ Multiple go.mod Files
```
go.mod
solutions/module-01-gcounter/go.mod
starters/module-01-gcounter/go.mod
```
**Problem:** Complexity. Students need to understand Go workspaces. Rejected for simplicity.

### ❌ .vscode/settings.json for gopls Exclusions
```json
{
  "gopls": {
    "directoryFilters": ["-solutions", "-starters"]
  }
}
```
**Problem:** Every user who clones the repo needs this config. Not portable. Rejected.

### ❌ Starter Files Named `*_starter.go`
```
starters/module-01-gcounter/
└── gcounter_starter.go  # Still has .go extension
```
**Problem:** gopls still indexes `.go` files even in starters/ directory. Must use `.tmpl` extension.

### ❌ Starter Templates WITH Package Declarations
```go
// starters/module-01-gcounter/gcounter.go.tmpl
package gcounter  // ❌ Makes it valid Go!

type GCounter struct {}
```
**Problem:** gopls parses valid Go files, even `.tmpl` ones. Must make them syntactically invalid by removing package declaration.

## What WORKS ✅

### ✅ Top-Level solutions/ and starters/ Directories
Clean separation, no nesting, different packages.

### ✅ .tmpl Extension WITHOUT Package Declaration
gopls cannot parse files without package declarations.

### ✅ Different Package Names for Solutions
`package pncountersolution` vs `package pncounter` eliminates namespace collisions.

### ✅ Single go.mod at Root
Simple, standard Go project structure. No workspace complexity.

### ✅ Industry Pattern (exercism.io)
Top-level directories for exercises, solutions, starters is proven pattern.

## Module README Template

Each module's README.md should include:

```markdown
# Module XX — Topic Name
## One-Line Description

**Concept:** Key CRDT properties taught  
**Challenge:** What students will build  
**Time:** Estimated completion time

---

## Introduction
[Explain the CRDT concept and motivation]

## The Problem
[Concrete example showing why naive approaches fail]

## How It Works
[Code examples showing the CRDT approach]

## Getting Started

Copy the starter templates:

\`\`\`bash
# From the module-XX-name directory
cp ../starters/module-XX-name/*.tmpl .
for f in *.tmpl; do mv "$f" "${f%.tmpl}"; done
\`\`\`

This gives you:
- `name.go` — struct definition and function stubs
- `name_test.go` — test suite with TODOs

**Why .tmpl?** Template files let you reset to the starting state anytime.

### Need to Reset?

\`\`\`bash
rm *.go
cp ../starters/module-XX-name/*.tmpl .
for f in *.tmpl; do mv "$f" "${f%.tmpl}"; done
\`\`\`

## The Challenge
[API specification with type definitions]

## Workflow

### Step 1: Set Up Your Files
[Instructions to copy templates]

### Step 2: Write Your Tests
[Guide to filling in test TODOs]

### Step 3: Implement
[Guide to implementation]

### Step 4: Run Your Tests
\`\`\`bash
go test .
\`\`\`

### Step 5: Compare with Solution
\`\`\`bash
cat ../solutions/module-XX-name/name_solution.go
cd ../solutions/module-XX-name && go test .
\`\`\`

## Hints
[Expandable hint sections]

## Next Module
[Link to next module]
```

## Go-Specific Conventions

### Imports Between Modules
When Module 02 needs Module 01's GCounter:

```go
// solutions/module-02-pncounter/pncounter_solution.go
package pncountersolution

import "github.com/zshainsky/crdt-tutorial/module-01-gcounter"

type PNCounter struct {
    p *gcounter.GCounter  // Use module-01-gcounter package
    n *gcounter.GCounter
}
```

**Key:** Import the working directory package, NOT the solution package.

### Starter Template Content Rules

TODO comments in `.go.tmpl` files must use **plain English or pseudocode only** — never actual Go code.

**Bad (gives away the answer):**
```go
// TODO: Merge both P and N counters
// pn.p.Merge(other.p)
// pn.n.Merge(other.n)
```

**Good (guides without spoiling):**
```go
// TODO: Merge both the positive and negative sub-counters
```

This applies to:
- Implementation stubs (`.go.tmpl`)
- Test stubs (`_test.go.tmpl`)
- Module README workflow and hint sections — describe *what* to do, not *how* in code

Hints in README `<details>` blocks may use pseudocode when the logic is genuinely hard to describe in English, but must never contain valid, compilable Go.

### Test File Conventions
- Always `*_test.go` suffix (Go convention)
- Never `starter_test.go` (bad naming)
- Test function names must start with `Test`
- Use table-driven tests where appropriate

### Package Naming
- Working directory: lowercase, no underscores (e.g., `package pncounter`)
- Solution: append `solution` (e.g., `package pncountersolution`)
- Avoid stuttering: `pncounter.PNCounter` not `pncounter.PNCounterStruct`

## Quality Checklist

Before committing a new module:

- [ ] README.md explains concept clearly with examples
- [ ] Starter templates have `.tmpl` extension
- [ ] Starter templates have NO package declaration
- [ ] Solution uses `package <name>solution`
- [ ] Solution tests pass: `cd solutions/module-XX && go test .`
- [ ] No gopls errors in VS Code with all three locations present
- [ ] Working directory initially empty except README.md
- [ ] Imports use correct package paths
- [ ] Module linked from main README.md
- [ ] Committed to feature branch before merging

## Summary: The Golden Rules

1. **Starters:** `.tmpl` extension + NO package declaration = gopls ignores completely
2. **Solutions:** Different package name = no namespace collision
3. **Single go.mod:** Simplicity over multi-module complexity
4. **Top-level directories:** solutions/ and starters/ at repo root
5. **Never commit working code:** Only README.md in module-XX-name/ initially

---

**Established:** March 22, 2026  
**Based on:** Module 01 implementation and gopls conflict resolution  
**Status:** Active conventions for all future modules
