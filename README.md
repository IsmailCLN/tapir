# Tapir – Test API Runner

*Automated API tests in a single YAML file, executed from a sleek terminal UI.*

---

## Table of Contents

- [Tapir – Test API Runner](#tapir--testapirunner)
  - [Table of Contents](#tableofcontents)
  - [Overview](#overview)
  - [Features](#features)
  - [Installation](#installation)
  - [Quick Start](#quickstart)
  - [Command Reference](#commandreference)
  - [YAML Suite Format](#yamlsuiteformat)
  - [Interactive TUI Shortcuts](#interactivetuishortcuts)
  - [Building from Source](#buildingfromsource)
  - [Contributing](#contributing)
  - [License](#license)

---

## Overview

**Tapir** is a Go‑based CLI that reads a YAML test‑suite file, fires HTTP requests in sequence,
checks the responses (status, headers, body, timing …) and displays the results in a
[Bubble Tea](https://github.com/charmbracelet/bubbletea) terminal UI.  Results can be exported to
Markdown or refreshed on demand – perfect for CI pipelines *and* local development.

![Tapir preview](docs/preview.gif)

---

## Features

* **YAML‑driven tests** – no code required; edit & commit your specs.
* **Assertions engine** – validate status code, JSON fields, headers and response time.
* **Interactive TUI** – coloured table of results with keyboard shortcuts.
* **One‑key export** – press **`p`** to save a styled Markdown report.
* **Hot reload** – press **`r`** to rerun the whole suite and update the table (1‑second cool‑down).
* **Schema validation** – `tapir validate <file>` ensures your YAML matches the expected format.
* **Sample generator** – `tapir generate example.yaml` creates a starter suite.
* **Configurable HTTP client** – shared transport, global timeout flag.

---

## Installation

```bash
# Requires Go 1.22+

# Install the latest released binary
go install github.com/IsmailCLN/tapir/cmd/tapir@latest

# Or clone and build from source
 git clone https://github.com/IsmailCLN/tapir.git
 cd tapir && go build -o tapir ./cmd
```

The resulting `tapir` binary can be copied anywhere in your `$PATH`.

---

## Quick Start

```bash
# Run a suite and open the TUI
 tapir run test-data/test.yaml

# Validate a file without executing requests
 tapir validate my-suite.yaml

# Generate a starter file
 tapir generate example.yaml
```

While the TUI is open you can press **`p`** to export a Markdown report, **`r`** to reload, **`c`** to copy and
**`q`** to quit.

---

## Command Reference

| Command                 | Description                                                       |
| ----------------------- | ----------------------------------------------------------------- |
| `tapir run <file>`      | Execute the test suite in *file* and show the interactive report. |
| `tapir validate <file>` | Check *file* against Tapir schema – returns non‑zero on error.    |
| `tapir generate <file>` | Write a minimal example suite to *file*.                          |

Global flags:

```text
--timeout   HTTP timeout per request (default 10s)
--verbose   Print request/response details to stdout while running
```

---

## YAML Suite Format

```yaml
tests:
  - name: Get Single User
    method: GET
    url: https://jsonplaceholder.typicode.com/users/1
    headers:
        Content-Type: application/json
    expect:
        status: 200
        body: | 
            {
                "id": 1,
                "name": "Leanne Graham",
                "username": "Bret",
                "email": "Sincere@april.biz",
                "address": {
                    "street": "Kulas Light",
                    "suite": "Apt. 556",
                    "city": "Gwenborough",
                    "zipcode": "92998-3874",
                    "geo": {
                        "lat": "-37.3159",
                        "lng": "81.1496"
                    }
                },
                "phone": "1-770-736-8031 x56442",
                "website": "hildegard.org",
                "company": {
                    "name": "Romaguera-Crona",
                    "catchPhrase": "Multi-layered client-server neural-net",
                    "bs": "harness real-time e-markets"
                }
            }
```

See **`test-data/test.yaml`** for a complete example.

---

## Interactive TUI Shortcuts

| Key     | Action                                                                                                  |
| ------- | ------------------------------------------------------------------------------------------------------- |
| **`q`** | Quit Tapir                                                                                              |
| **`p`** | Print report to `tapir-report-YYYYMMDD.md`                                                              |
| **`c`** | Copy report to clipboard (if OS supported)                                                              |
| **`r`** | Reload the entire suite (blocked if pressed again within 1 second → *"Refresh requests too frequent."*) |

---

## Building from Source

```bash
# Clone repository
 git clone https://github.com/IsmailCLN/tapir.git
 cd tapir

# Run tests
 go test ./...

# Build for your platform
 go build -o tapir ./cmd
```

To cross‑compile:

```bash
GOOS=linux  GOARCH=amd64 go build -o tapir-linux  ./cmd
GOOS=darwin GOARCH=arm64 go build -o tapir-mac    ./cmd
```

---

## Contributing

Pull requests are welcome – please open an issue first to discuss what you would like to change.
Make sure `go vet ./...` and `go test ./...` pass before submitting.

---

## License

Tapir is released under the MIT License.  See [`LICENSE`](LICENSE) for details.
