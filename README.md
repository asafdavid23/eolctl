# End of Life Control (eolctl)

`eolctl` is a command-line tool that helps developers and system administrators manage and monitor the end-of-life (EOL) status of various software products. This tool queries the [endoflife.date](https://endoflife.date/) API to provide real-time information about software versions and their support lifecycle.

## Features

- Check the EOL status of various programming languages and frameworks.
- Export to a JSON file.
- Scan your code project (GO and Python only supported for now.)
- Get a custom range for versions.
- Compare between existing and future version

## Installation

### Prerequisites

Ensure you have the following installed on your machine:

- Go (version 1.16 or higher)

### Clone the Repository

```bash
git clone https://github.com/asafdavid23/eolctl.git
cd eolctl/
go build bin/eolctl
```


## Examples

```bash
bin/eolctl get product --name python --version 3.12 | jq

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
bin/eolctl get prodeuct --name python --min --max 3.12
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