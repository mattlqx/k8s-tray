# ADR-002: Code Hygiene with Pre-commit Hooks

## Status

Accepted - Implemented

## Date

2025-07-13

## Context

To maintain high code quality and consistency across the k8s-tray project, we need to
establish automated code hygiene practices. Manual code reviews alone are insufficient to
catch formatting issues, linting violations, and ensure adherence to commit message
standards. Pre-commit hooks provide an automated way to enforce these standards before
code enters the repository.

## Decision

We will implement a comprehensive pre-commit hook system using the `pre-commit` framework
to enforce code quality standards across multiple dimensions:

## Pre-commit Hook Configuration

### 1. Go Code Quality

- **Linting**: golangci-lint for comprehensive Go code analysis
- **Formatting**: gofmt and goimports for consistent code formatting
- **Security**: gosec for security vulnerability detection
- **Complexity**: gocyclo for cyclomatic complexity analysis

### 2. Commit Message Standards

- **Conventional Commits**: Enforce conventional commit format
- **Message Length**: Limit subject line to 50 characters
- **Body Wrapping**: Wrap commit message body at 72 characters
- **Reference Links**: Validate issue/PR references

### 3. Test Requirements

- **Unit Tests**: Run all unit tests before commit
- **Coverage**: Ensure minimum test coverage threshold
- **Test Naming**: Validate test function naming conventions
- **Benchmark Tests**: Run performance benchmarks

### 4. Documentation Standards

- **Markdown Formatting**: markdownlint for consistent markdown
- **Link Validation**: Check for broken links in documentation
- **Spelling**: Automated spell checking
- **TOC Generation**: Auto-generate table of contents

## Implementation Details

### Pre-commit Configuration File

```yaml
# .pre-commit-config.yaml
repos:
  # Commit message validation
  - repo: https://github.com/compilerla/conventional-pre-commit
    rev: v2.4.0
    hooks:
      - id: conventional-pre-commit
        name: Conventional Commit
        description: Validates commit message format
        stages: [commit-msg]

  # General code quality
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.4.0
    hooks:
      - id: trailing-whitespace
        name: Trim Trailing Whitespace
        description: Removes trailing whitespace
      - id: end-of-file-fixer
        name: Fix End of Files
        description: Ensures files end with newline
      - id: check-yaml
        name: Check YAML
        description: Validates YAML files
      - id: check-json
        name: Check JSON
        description: Validates JSON files
      - id: check-merge-conflict
        name: Check Merge Conflicts
        description: Checks for merge conflict markers
      - id: check-added-large-files
        name: Check Large Files
        description: Prevents large files from being committed
        args: [--maxkb=500]

  # Markdown formatting
  - repo: https://github.com/igorshubovych/markdownlint-cli
    rev: v0.35.0
    hooks:
      - id: markdownlint
        name: Markdown Lint
        description: Lints markdown files
        args: [--fix]
        files: \.(md|markdown)$

  # Spell checking
  - repo: https://github.com/crate-ci/typos
    rev: v1.16.1
    hooks:
      - id: typos
        name: Spell Check
        description: Checks for typos in code and documentation
```

### golangci-lint Configuration

```yaml
# .golangci.yml
run:
  timeout: 5m
  issues-exit-code: 1
  tests: true
  modules-download-mode: readonly

output:
  format: colored-line-number
  print-issued-lines: true
  print-linter-name: true

linters-settings:
  govet:
    check-shadowing: true
  golint:
    min-confidence: 0.8
  gocyclo:
    min-complexity: 15
  goimports:
    local-prefixes: github.com/k8s-tray
  goconst:
    min-len: 3
    min-occurrences: 3
  misspell:
    locale: US
  unused:
    check-exported: false
  unparam:
    check-exported: false
  funlen:
    lines: 80
    statements: 40

linters:
  enable:
    - bodyclose
    - deadcode
    - depguard
    - dogsled
    - dupl
    - errcheck
    - exportloopref
    - funlen
    - gochecknoinits
    - goconst
    - gocritic
    - gocyclo
    - gofmt
    - goimports
    - golint
    - gomnd
    - goprintffuncname
    - gosec
    - gosimple
    - govet
    - ineffassign
    - interfacer
    - lll
    - maligned
    - misspell
    - nakedret
    - noctx
    - nolintlint
    - rowserrcheck
    - scopelint
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - varcheck
    - whitespace

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - gomnd
        - funlen
        - goconst
```

### Conventional Commits Configuration

```yaml
# .conventional-commits.yaml
types:
  - feat      # New feature
  - fix       # Bug fix
  - docs      # Documentation changes
  - style     # Code style changes (formatting, etc.)
  - refactor  # Code refactoring
  - perf      # Performance improvements
  - test      # Adding or modifying tests
  - build     # Build system changes
  - ci        # CI/CD changes
  - chore     # Maintenance tasks
  - revert    # Revert changes

scopes:
  - api       # API related changes
  - ui        # User interface changes
  - config    # Configuration changes
  - docs      # Documentation
  - test      # Testing
  - build     # Build system
  - deploy    # Deployment

subject-min-length: 10
subject-max-length: 50
body-max-line-length: 72
footer-max-line-length: 72
```

### Markdownlint Configuration

```json
{
  "MD001": true,
  "MD002": true,
  "MD003": { "style": "atx" },
  "MD004": { "style": "dash" },
  "MD005": true,
  "MD006": true,
  "MD007": { "indent": 2 },
  "MD009": { "br_spaces": 2 },
  "MD010": { "code_blocks": true },
  "MD011": true,
  "MD012": { "maximum": 2 },
  "MD013": { "line_length": 100 },
  "MD014": true,
  "MD018": true,
  "MD019": true,
  "MD020": true,
  "MD021": true,
  "MD022": true,
  "MD023": true,
  "MD024": { "allow_different_nesting": true },
  "MD025": true,
  "MD026": { "punctuation": ".,;:!" },
  "MD027": true,
  "MD028": true,
  "MD029": { "style": "ordered" },
  "MD030": true,
  "MD031": true,
  "MD032": true,
  "MD033": { "allowed_elements": ["br", "pre", "code"] },
  "MD034": true,
  "MD035": { "style": "---" },
  "MD036": true,
  "MD037": true,
  "MD038": true,
  "MD039": true,
  "MD040": true,
  "MD041": true,
  "MD042": true,
  "MD043": true,
  "MD044": true,
  "MD045": true,
  "MD046": { "style": "fenced" },
  "MD047": true,
  "MD048": { "style": "backtick" },
  "MD049": { "style": "asterisk" },
  "MD050": { "style": "asterisk" }
}
```

## Setup and Installation

### 1. Install Pre-commit

```bash
# Install pre-commit
pip install pre-commit

# Or using homebrew on macOS
brew install pre-commit

# Or using package manager
curl -sSL https://install.python-poetry.org | python3 -
```

### 2. Install Pre-commit Hooks

```bash
# Install hooks in repository
pre-commit install

# Install commit message hooks
pre-commit install --hook-type commit-msg

# Run against all files (initial setup)
pre-commit run --all-files
```

### 3. Development Workflow Integration

```makefile
# Makefile targets
.PHONY: setup-dev lint test format pre-commit-install

setup-dev: pre-commit-install
 go mod download
 go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

pre-commit-install:
 pre-commit install
 pre-commit install --hook-type commit-msg

lint:
 golangci-lint run

test:
 go test -race -cover ./...

format:
 gofmt -w .
 goimports -w .

pre-commit-run:
 pre-commit run --all-files
```

## Hook Execution Order

### Pre-commit Hooks (Before Commit)

1. **File Formatting**: trailing whitespace, end-of-file fixes
2. **Go Formatting**: gofmt, goimports
3. **Go Linting**: golangci-lint
4. **Go Security**: gosec
5. **Go Testing**: unit tests with coverage
6. **Markdown Formatting**: markdownlint
7. **Spell Checking**: typos
8. **File Validation**: YAML, JSON syntax

### Commit Message Hooks (After Commit Message)

1. **Conventional Commits**: format validation
2. **Message Length**: subject and body length checks
3. **Reference Validation**: issue/PR link validation

## Quality Gates

### Test Coverage Requirements

- **Minimum Coverage**: 80% for all packages
- **Coverage Report**: Generated on each test run
- **Coverage Exclusions**: Test files, generated code
- **Benchmark Tests**: Performance regression detection

### Linting Standards

- **Cyclomatic Complexity**: Maximum 15 per function
- **Function Length**: Maximum 80 lines, 40 statements
- **Security**: No security vulnerabilities allowed
- **Import Organization**: Grouped and sorted imports

### Documentation Standards

- **API Documentation**: All public functions documented
- **README Updates**: Updated for feature changes
- **ADR Updates**: New ADRs for significant decisions
- **Changelog**: Updated for user-facing changes

## CI/CD Integration

### GitHub Actions Workflow

```yaml
name: Code Quality
on: [push, pull_request]

jobs:
  pre-commit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-python@v4
        with:
          python-version: '3.x'
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - uses: pre-commit/action@v3.0.0
```

### Local Development

- **Pre-commit installation**: Required for all developers
- **IDE Integration**: EditorConfig and Go plugin settings
- **Make targets**: Convenient commands for common tasks
- **Documentation**: Clear setup instructions in README

## Performance Considerations

### Hook Optimization

- **Parallel Execution**: Independent hooks run in parallel
- **File Filtering**: Only run hooks on relevant file types
- **Caching**: Cache dependencies and build artifacts
- **Incremental**: Only check changed files where possible

### Developer Experience

- **Fast Feedback**: Hooks complete in <30 seconds
- **Clear Messages**: Descriptive error messages
- **Auto-fixing**: Automatic fixes where possible
- **Skip Options**: Emergency skip for urgent fixes

## Monitoring and Metrics

### Quality Metrics

- **Pre-commit Success Rate**: Track hook failure rates
- **Code Coverage Trends**: Monitor coverage over time
- **Linting Violations**: Track violation types and frequency
- **Commit Message Quality**: Adherence to conventional commits

### Developer Adoption

- **Hook Installation**: Track developer setup completion
- **Bypass Frequency**: Monitor emergency skips
- **Feedback Collection**: Regular developer experience surveys
- **Training Needs**: Identify areas for additional training

## Exception Handling

### Emergency Bypasses

```bash
# Skip all pre-commit hooks (emergency only)
git commit --no-verify -m "emergency fix"

# Skip specific hooks
SKIP=golangci-lint git commit -m "docs: update README"
```

### Legitimate Exceptions

- **Third-party Code**: Exclude vendor directories
- **Generated Code**: Exclude auto-generated files
- **Legacy Code**: Gradual migration strategy
- **Performance Critical**: Exempt specific optimized code

## Maintenance and Updates

### Regular Updates

- **Monthly Reviews**: Update hook versions and configurations
- **Dependency Updates**: Keep tools and frameworks current
- **Rule Adjustments**: Refine rules based on team feedback
- **Documentation**: Keep setup instructions current

### Continuous Improvement

- **Feedback Integration**: Regular team retrospectives
- **New Tools**: Evaluate and integrate new quality tools
- **Performance Optimization**: Continuously improve hook performance
- **Best Practices**: Stay current with Go community standards

## Success Metrics

### Code Quality Improvements

- **Reduced Bug Reports**: Fewer production issues
- **Faster Code Reviews**: Less time spent on formatting issues
- **Consistent Style**: Uniform code appearance
- **Better Documentation**: Improved API documentation

### Development Efficiency

- **Faster Onboarding**: New developers productive sooner
- **Reduced Rework**: Fewer commits fixing quality issues
- **Automated Quality**: Less manual quality checking
- **Confident Releases**: Higher confidence in code quality

## Decision Outcome

This ADR establishes a comprehensive code hygiene system using pre-commit hooks to
automatically enforce code quality standards, conventional commit messages, test
requirements, and documentation formatting. This will improve code quality, reduce review
overhead, and ensure consistent standards across the k8s-tray project.
