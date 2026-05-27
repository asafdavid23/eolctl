# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added
- `--suggest-version` flag on `scan project` — AI-generated upgrade recommendations per component, powered by Claude Opus
- `SuggestUpgradePath` function in `pkg/ai/advisor.go` with dedicated `UpgradeItem` struct

---

## [1.3.0] - 2026-05-27

### Added
- AI-powered stack detection via `pkg/ai/detector.go` — replaces brittle file scanner; uses Claude Haiku to identify language and version from project files
- Monorepo support — `DetectStack` now returns `[]StackInfo`, scanning all languages in a single pass
- Risk scoring — `pkg/helpers/risk.go` classifies each component as CRITICAL / HIGH / MEDIUM / LOW / UNKNOWN based on days until EOL
- `--risk-report` flag on `scan project` — AI-generated risk narrative using Claude Opus
- Vendor directory exclusions in `DetectStack` — skips `node_modules`, `vendor`, `.git`, `.venv`, `__pycache__`

### Fixed
- URL-encode product and version values before building API request URLs
- Version range comparison now uses numeric segment comparison instead of lexicographic `strings.Compare` (fixes `"9"` > `"10"` ordering)
- `GetStringValue` now handles boolean EOL values (returns `"true"` / `"false"`)
- `CheckProductEOL` bool-false case no longer falls through to a zero-time comparison
- `homeDir` error in cache initialization is now properly handled instead of silently ignored
- Product existence check uses exact JSON array match instead of substring search (fixes false positives like `"java"` matching `"javascript"`)
- Markdown fence stripping in Claude responses uses index-based JSON extraction instead of fragile `TrimPrefix` chaining
- Removed `log` package from `cmd/root.go`; replaced with `fmt.Fprintf(os.Stderr)` for consistency
- `cmd/project.go` argument validation replaced with `cobra.ExactArgs(1)` to prevent `args[0]` panic

---

## [1.2.4] - 2026-05-26

### Fixed
- Various typos and error handling improvements across multiple files

---

## [1.2.3] - 2025-01-06

### Fixed
- `CheckProductEOL` return value corrected for bool EOL responses

---

## [1.2.2] - 2024-12-23

### Fixed
- Local cache encoding/decoding stability fixes

---

## [1.2.1] - 2024-12-23

### Fixed
- Cache output handling

---

## [1.2.0] - 2024-12-19

### Added
- Local cache for `available-products` API responses to reduce repeated network calls
- `CheckProductEOL` helper function for programmatic EOL checks
- Helpers moved to `pkg/helpers` package

---

## [1.1.4] - 2024-12-05

### Fixed
- Multi-project scan stability improvements

---

## [1.1.3] - 2024-12-03

### Fixed
- Config file not found error now handled gracefully instead of fatally

---

## [1.1.2] - 2024-12-02

### Fixed
- Default output format set to `table`

---

## [1.1.1] - 2024-11-15

### Fixed
- Table rendering fixes for `get products` and `list products` commands

---

## [1.1.0] - 2024-10-29

### Added
- Python project scanning support (`Pipfile`, `pyproject.toml`, `setup.py`)
- Version column in scan output
- Dynamic version embedded at build time

---

## [1.0.3] - 2024-10-28

### Fixed
- Multi-project scanning — scan now walks subdirectories and reports per-project results

---

## [1.0.2] - 2024-10-20

### Fixed
- Workflow and release pipeline fixes

---

## [1.0.1] - 2024-10-20

### Fixed
- Workflow fixes

---

## [1.0.0] - 2024-10-20

### Added
- Initial stable release
- `get product` command — query EOL status for a specific product and version
- `get list` command — list all available products
- Custom version range filtering with `--min` / `--max`
- `scan project` command — detect Go and JavaScript project language/version and fetch EOL status
- Table and JSON output formats
- Export to JSON file with `--output-path`
- Structured logging with configurable log level
- GitHub Actions CI workflow
- Homebrew tap support
