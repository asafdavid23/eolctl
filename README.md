![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

# End of Life Control (eolctl)

`eolctl` is a command-line tool that helps developers and system administrators manage and monitor the end-of-life (EOL) status of various software products. It queries the [endoflife.date](https://endoflife.date/) API for real-time version lifecycle information, and uses AI (powered by Claude) to automatically detect project stacks and generate risk narratives.

## Features

- Check the EOL status of various programming languages and frameworks.
- AI-powered project scanning — automatically detects language and version from your codebase using Claude.
- Monorepo support — detects multiple languages (e.g. Go backend + Node.js frontend) in a single scan.
- Risk scoring — classifies each component as CRITICAL / HIGH / MEDIUM / LOW based on days until EOL.
- AI risk narrative — generates a concise, prioritized risk summary using Claude.
- Custom version range filtering.
- Export results to JSON file.

## Prerequisites

The `scan project` command and the `--risk-report` flag require an Anthropic API key:

```bash
export ANTHROPIC_API_KEY=your_api_key_here
```

## Installation

You can install eolctl by downloading the latest release from the [releases page](https://github.com/asafdavid23/eolctl/releases)

```bash
curl -LO https://github.com/asafdavid23/eolctl/releases/latest/download/eolctl
chmod +x eolctl
sudo mv eolctl /usr/local/bin/
```

For Windows, download the binary and add it to your system PATH.

Alternatively, you can use brew (macOS / Linux):

```bash
brew tap asafdavid23/tap
brew update
brew install eolctl
```

## Usage

### Look up a specific product version

```bash
eolctl get product --name go --version 1.23
```

```
+--------+-------------------+-------------+-----+-----+---------+
| LATEST | LATESTRELEASEDATE | RELEASEDATE | LTS | EOL | SUPPORT |
+--------+-------------------+-------------+-----+-----+---------+
| 1.23.2 | 2024-10-01        | 2024-08-13  |     |     |         |
+--------+-------------------+-------------+-----+-----+---------+
```

### Filter a version range

```bash
eolctl get product --name python --min 3.10 --max 3.12
```

```
+-------+---------+-------------------+-------------+-----+------------+------------+
| CYCLE | LATEST  | LATESTRELEASEDATE | RELEASEDATE | LTS |    EOL     |  SUPPORT   |
+-------+---------+-------------------+-------------+-----+------------+------------+
|  3.12 | 3.12.7  | 2024-10-01        | 2023-10-02  |     | 2028-10-31 | 2025-04-02 |
|  3.11 | 3.11.10 | 2024-09-07        | 2022-10-24  |     | 2027-10-31 | 2024-04-01 |
|  3.10 | 3.10.15 | 2024-09-07        | 2021-10-04  |     | 2026-10-31 | 2023-04-05 |
+-------+---------+-------------------+-------------+-----+------------+------------+
```

### Scan a project (AI-powered)

`eolctl` uses Claude to read your project files (`go.mod`, `package.json`, `requirements.txt`, etc.) and automatically detect the language and version — no manual configuration needed.

```bash
eolctl scan project ./myapp --output table
```

```
+---------+---------+------------+--------+
| PRODUCT | VERSION |    EOL     |  RISK  |
+---------+---------+------------+--------+
| go      | 1.22    | 2025-08-01 | MEDIUM |
+---------+---------+------------+--------+
```

### Scan a monorepo

For projects with multiple languages (e.g. a Go backend and a Node.js frontend), all stacks are detected and reported in a single scan:

```bash
eolctl scan project ./monorepo --output table
```

```
+---------+---------+------------+----------+
| PRODUCT | VERSION |    EOL     |   RISK   |
+---------+---------+------------+----------+
| go      | 1.22    | 2025-08-01 | MEDIUM   |
| nodejs  | 18      | 2025-04-30 | CRITICAL |
+---------+---------+------------+----------+
```

### AI risk narrative

Add `--risk-report` to get a concise, AI-generated summary that prioritises the most critical issues:

```bash
eolctl scan project ./monorepo --output table --risk-report
```

```
--- AI Risk Summary ---
Your Node.js 18 dependency reached end-of-life on April 30, 2025 and should be
upgraded to Node.js 20 or 22 immediately to avoid security exposure. Go 1.22
will reach EOL in August 2025, leaving a narrow window for an upgrade to 1.23
before support ends. Prioritise the Node.js upgrade given its critical status.
```

## Risk levels

| Level    | Condition                        |
|----------|----------------------------------|
| CRITICAL | Already EOL, or EOL is a boolean `true` |
| HIGH     | EOL within 90 days               |
| MEDIUM   | EOL within 180 days              |
| LOW      | EOL more than 180 days away      |
| UNKNOWN  | EOL date could not be determined |

## CI Integration

`eolctl` is well-suited for CI/CD pipelines. Here's an example GitHub Actions workflow:

```yaml
name: Check EOL

on: [push]

jobs:
  check-eol:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Run eolctl
      env:
        ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
      run: |
        curl -LO https://github.com/asafdavid23/eolctl/releases/latest/download/eolctl
        chmod +x eolctl
        ./eolctl scan project . --output table
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
