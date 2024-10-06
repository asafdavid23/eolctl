# End of Life Control (eolctl)

`eolctl` is a command-line tool that helps developers and system administrators manage and monitor the end-of-life (EOL) status of various software products. This tool queries the [endoflife.date](https://endoflife.date/) API to provide real-time information about software versions and their support lifecycle.

## Features

- Check the EOL status of various programming languages and frameworks.

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

