![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

# End of Life Control (eolctl)

`eolctl` is a command-line tool that helps developers and system administrators manage and monitor the end-of-life (EOL) status of various software products. This tool queries the [endoflife.date](https://endoflife.date/) API to provide real-time information about software versions and their support lifecycle.

## Features

- Check the EOL status of various programming languages and frameworks.
- Scan your code project (GO and JS only supported for now.)
- Get a custom range for versions.
- Export to JSON file


## Installation

You can install eolctl by downloading the latest release from the [releases page](https://github.com/asafdavid23/eolctl/releases)

```bash
curl -LO https://github.com/asafdavid23/eolctl/releases/latest/download/eolctl
chmod +x eolctl
sudo mv eolctl /usr/local/bin/
```
For Windows, download the binary and add it to your system PATH.

Alternatively, you use brew for (MacOS / Linux)
```bash
brew tap asafdavid23/tap
brew update
brew install eolctl
```


## Usage

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

```bash
eolctl get prodeuct --name python --min 3.10 --max 3.12
```

```
+-------+---------+-------------------+-------------+-----+------------+------------+
| CYCLE | LATEST  | LATESTRELEASEDATE | RELEASEDATE | LTS |    EOL     |  SUPPORT   |
+-------+---------+-------------------+-------------+-----+------------+------------+
|  3.13 | 3.13.0  | 2024-10-07        | 2024-10-07  |     | 2029-10-31 | 2026-10-01 |
|  3.12 | 3.12.7  | 2024-10-01        | 2023-10-02  |     | 2028-10-31 | 2025-04-02 |
|  3.11 | 3.11.10 | 2024-09-07        | 2022-10-24  |     | 2027-10-31 | 2024-04-01 |
|  3.10 | 3.10.15 | 2024-09-07        | 2021-10-04  |     | 2026-10-31 | 2023-04-05 |
+-------+---------+-------------------+-------------+-----+------------+------------+
```

```bash
eolctl scan project /tmp/testproj --output table
```

```
+---------+---------+------------+
| PRODUCT | VERSION |    EOL     |
+---------+---------+------------+
| go      |    1.20 | 2024-02-06 |
| python  |    3.10 | 2026-10-31 |
+---------+---------+------------+
```

## CI Integration

`eolctl` is perfect for use in CI/CD pipelines to ensure the languages and versions in your project are not deprecated. Here's an example GitHub Action workflow:

```yaml
name: Check EOL

on: [push]

jobs:
  check-eol:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Run eolctl
      run: |
        curl -LO https://github.com/asafdavid23/eolctl/releases/latest/download/eolctl
        chmod +x eolctl
        ./eolctl scan proejct .
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes. Refer to the Contributing Guide for more information.

## License

This project is licensed under the MIT License - see the LICENSE file for details.