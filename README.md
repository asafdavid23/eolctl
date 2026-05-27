![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

# End of Life Control (eolctl)

`eolctl` is a command-line tool that helps developers and system administrators manage and monitor the end-of-life (EOL) status of various software products. It queries the [endoflife.date](https://endoflife.date/) API for real-time version lifecycle information, and uses AI (powered by Claude) to automatically detect project stacks and generate risk narratives.

## Features

- Check the EOL status of various programming languages and frameworks.
- AI-powered project scanning — automatically detects language and version from your codebase using Claude.
- Monorepo support — detects multiple languages (e.g. Go backend + Node.js frontend) in a single scan.
- **Kubernetes cluster scanning** — lists all Helm releases across namespaces and checks each chart's app version for EOL status.
- **ArtifactHub fallback** — for charts not tracked by endoflife.date, falls back to ArtifactHub to derive risk from version staleness and deprecation status.
- Risk scoring — classifies each component as CRITICAL / HIGH / MEDIUM / LOW based on days until EOL.
- AI risk narrative — generates a concise, prioritized risk summary using Claude.
- AI upgrade suggestions — recommends specific versions to upgrade to for each EOL component.
- Custom version range filtering.
- Export results to JSON file.

## Prerequisites

The `scan project` and `scan cluster` commands, as well as the `--risk-report` and `--suggest-version` flags, require an Anthropic API key:

```bash
export ANTHROPIC_API_KEY=your_api_key_here
```

The `scan cluster` command additionally requires:

- [Helm](https://helm.sh/docs/intro/install/) v3+ installed and available on your `$PATH`.
- A valid `kubeconfig` pointing at the target cluster (the default `~/.kube/config` is used automatically).

```bash
# verify Helm can reach your cluster
helm list --all-namespaces
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

### Scan a Kubernetes cluster

`eolctl` lists every Helm release in your cluster (across all namespaces), uses Claude to map each chart to its endoflife.date product slug, and then checks the app version for EOL status and risk level. For charts that endoflife.date does not track, it falls back to [ArtifactHub](https://artifacthub.io/) and derives risk from version staleness.

```bash
eolctl scan cluster --output table
```

```
+------------------+-------------+--------------+---------+--------------------+----------+
|     RELEASE      |  NAMESPACE  |   PRODUCT    | VERSION |        EOL         |   RISK   |
+------------------+-------------+--------------+---------+--------------------+----------+
| nginx-ingress    | ingress      | nginx        | 1.23    | 2025-04-01         | CRITICAL |
| cert-manager     | cert-manager | cert-manager | 1.11    | latest: 1.14.2     | HIGH     |
| prometheus       | monitoring   | prometheus   | 2.44    | latest: 2.52.0     | MEDIUM   |
| redis            | default      | redis        | 7.0     | 2027-01-01         | LOW      |
+------------------+-------------+--------------+---------+--------------------+----------+
```

Output as JSON:

```bash
eolctl scan cluster --output json
```

```json
[
  {
    "release": "nginx-ingress",
    "namespace": "ingress",
    "chart": "ingress-nginx-4.7.1",
    "product": "nginx",
    "version": "1.23",
    "eol": "2025-04-01",
    "risk": "CRITICAL",
    "days_until_eol": -57
  }
]
```

#### AI risk report and upgrade suggestions for cluster

The `--risk-report` and `--suggest-version` flags work the same way as on `scan project`:

```bash
eolctl scan cluster --output table --risk-report --suggest-version
```

```
--- AI Risk Summary ---
Your nginx 1.23 deployment reached end-of-life on April 1, 2025 and is
actively accumulating unpatched CVEs. cert-manager 1.11 is three major versions
behind 1.14; upgrade to unblock critical TLS policy fixes. Prometheus and Redis
are healthy but should be scheduled for minor-version bumps in the next cycle.

--- AI Upgrade Suggestions ---
For nginx 1.23 (EOL: 2025-04-01, CRITICAL): upgrade to 1.27 (current stable).
For cert-manager 1.11 (HIGH): upgrade to 1.14.2 via Helm chart bump.
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

### AI upgrade suggestions

Add `--suggest-version` to get specific, actionable upgrade recommendations for each component:

```bash
eolctl scan project ./monorepo --output table --suggest-version
```

```
--- AI Upgrade Suggestions ---
For Go 1.22 (EOL: 2025-08-01, MEDIUM risk): upgrade to Go 1.23, which is the
current stable release and supported until mid-2026. For Node.js 20 (EOL:
2026-04-30, LOW risk): consider planning a migration to Node.js 22 (LTS) ahead
of the EOL date to avoid a future critical-risk window.
```

Both flags can be combined in a single run:

```bash
eolctl scan project ./monorepo --output table --risk-report --suggest-version
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
