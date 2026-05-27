# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added
- `--suggest-version` and `--risk-report` promoted to global persistent flags — available on all commands
- `--risk-report` and `--suggest-version` support added to `get product` command
- Risk level color coding in table output — CRITICAL (bold red), HIGH (red), MEDIUM (yellow), LOW (green)
- `cmd/helpers.go` — shared `printRiskNarrative`, `printUpgradeSuggestions`, and `renderRichRow` helpers used by all commands
- Release workflow auto-extracts `[Unreleased]` changelog section and attaches it as GitHub release notes
- Release workflow bumps CHANGELOG.md on publish — moves `[Unreleased]` to the tagged version and inserts a fresh `[Unreleased]` section

### Fixed
- `SuggestUpgradePath` prompt updated to enforce plain text — prevents Claude returning markdown tables that misalign in terminal output

---

## [1.3.0] - 2026-05-27

### Added
- AI-powered stack detection via `pkg/ai/detector.go` — replaces brittle regex scanner; uses Claude Haiku to identify language and version from project files
- Monorepo support — `DetectStack` returns `[]StackInfo`, detecting all languages in a single directory walk
- Risk scoring — `pkg/helpers/risk.go` classifies each component as CRITICAL / HIGH / MEDIUM / LOW / UNKNOWN based on days until EOL
- `--risk-report` flag on `scan project` — AI-generated risk narrative using Claude Opus
- Vendor directory exclusions in stack detection — skips `node_modules`, `vendor`, `.git`, `.venv`, `__pycache__`
- `--suggest-version` flag on `scan project` — AI-generated upgrade recommendations per component

### Fixed
- URL-encode product and version values before building API request URLs
- Version range comparison now uses numeric segment comparison instead of lexicographic `strings.Compare` (fixes `"9"` > `"10"` ordering)
- `GetStringValue` now handles boolean EOL values (returns `"true"` / `"false"`)
- `CheckProductEOL` bool-false case no longer falls through to a zero-time comparison
- `homeDir` error in cache initialization is now properly handled instead of silently ignored
- Product existence check uses exact JSON array match instead of substring search (fixes false positives like `"java"` matching `"javascript"`)
- Markdown fence stripping in Claude responses uses index-based JSON extraction instead of fragile `TrimPrefix` chaining
- Replaced `log` package in `cmd/root.go` with `fmt.Fprintf` for consistent output handling
- `cmd/project.go` argument validation replaced with `cobra.ExactArgs(1)` to prevent panic on missing argument

---

## [1.2.4] - 2026-05-26

### Fixed
- Various typos and error handling improvements across multiple files

---

## [1.2.3] - 2025-01-06

### Fixed
- Local cache stability fix

---

## [1.2.2] - 2024-12-23

### Added
- `CheckProductEOL` function improvements

---

## [1.2.1] - 2024-12-23

### Added
- Local cache for `available-products` API responses to reduce repeated network calls
- `CheckProductEOL` helper function for programmatic EOL status checks
- Helpers moved to `pkg/helpers` package

### Fixed
- Output rendering fix for table display

---

## [1.2.0] - 2024-12-19

### Added
- Multi-project scanning — `scan project` now walks subdirectories and reports per-project results

---

## [1.1.4] - 2024-12-05

### Fixed
- Config file not found error now handled gracefully instead of fatally

---

## [1.1.3] - 2024-12-03

### Added
- Version column added to `scan project` table output

---

## [1.1.2] - 2024-12-02

### Added
- Dynamic version string embedded at build time via ldflags

---

## [1.1.1] - 2024-11-15

### Fixed
- Default output format set to `table` when `--output` flag is not specified

---

## [1.1.0] - 2024-10-29

### Added
- Python project scanning support (`Pipfile`, `pyproject.toml`, `setup.py`)

---

## [1.0.3] - 2024-10-28

### Fixed
- Table output fixed for `get product` and `list` commands
- Print table alignment fixes

---

## [1.0.2] - 2024-10-20

### Fixed
- CI/CD workflow and release pipeline fixes

---

## [1.0.1] - 2024-10-20

### Fixed
- Workflow fixes

---

## [1.0.0] - 2024-10-20

### Added
- Initial stable release
- `get product` command — query EOL status for a specific product and version
- `get list` command — list all available products from endoflife.date
- Custom version range filtering with `--min` / `--max`
- `scan project` command — detect Go and JavaScript project language/version and fetch EOL status
- Table and JSON output formats
- Export results to JSON file with `--output-path`
- Structured logging with configurable log level via `--log-level`
- Version identification from `go.mod`, `package.json`, `requirements.txt`
- GitHub Actions CI workflow with lint and staticcheck
- Homebrew tap support
- MIT License
