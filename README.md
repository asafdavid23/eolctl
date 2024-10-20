![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

# End of Life Control (eolctl)

`eolctl` is a command-line tool that helps developers and system administrators manage and monitor the end-of-life (EOL) status of various software products. This tool queries the [endoflife.date](https://endoflife.date/) API to provide real-time information about software versions and their support lifecycle.

## Features

- Check the EOL status of various programming languages and frameworks.
- Scan your code project (GO and JS only supported for now.)
- Get a custom range for versions.
- Export to JSON file


## Installation

### Prerequisites

Ensure you have the following installed on your machine:

- Go (version 1.16 or higher)

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
eolctl get product --name python --version 3.12

```

```
{
  "releaseDate": "2023-10-02",
  "eol": "2028-10-31",
  "latest": "3.12.7",
  "latestReleaseDate": "2024-10-01",
  "lts": false,
  "support": "2025-04-02"
}

```

```bash
bin/eolctl get prodeuct --name python --min 3.10 --max 3.12
```

```
[
    {
        "cycle": "3.12",
        "latest": "3.12.7",
        "latestReleaseDate": "2024-10-01",
        "releaseDate": "2023-10-02"
    },
    {
        "cycle": "3.11",
        "latest": "3.11.10",
        "latestReleaseDate": "2024-09-07",
        "releaseDate": "2022-10-24"
    },
    {
        "cycle": "3.10",
        "latest": "3.10.15",
        "latestReleaseDate": "2024-09-07",
        "releaseDate": "2021-10-04"
    }
]
```