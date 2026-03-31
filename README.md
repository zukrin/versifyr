# versifyr

[![CI](https://github.com/zukrin/versifyr/actions/workflows/pipeline.yml/badge.svg)](https://github.com/zukrin/versifyr/actions/workflows/pipeline.yml)
[![Latest Release](https://img.shields.io/github/v/release/zukrin/versifyr)](https://github.com/zukrin/versifyr/releases)
![Go Test Coverage](https://img.shields.io/badge/coverage-80.4%25-brightgreen)

`versifyr` is a specialized CLI tool designed to synchronize project versions across multiple files. It supports various formats (Go, YAML, XML, Java, JSON, etc.) by using powerful Go templates and Sprig functions to ensure consistency across your entire codebase.

Whether you need to update a Helm chart, a Maven `pom.xml`, or a Go constant, `versifyr` automates the process through a single command.

## Core Features

- **Format Agnostic**: Works with any text-based file format.
- **Template Driven**: Uses standard Go `text/template` syntax.
- **Extensible**: Includes [Sprig functions](http://masterminds.github.io/sprig/) for complex string manipulations.
- **Flexible Configuration**: Define templates directly in source comments or in a central YAML configuration.
- **Context Aware**: Automatically provides values like `latesttag`, `actualdate`, and `actualtimestamp`.

## Installation

```sh
go install github.com/zukrin/versifyr/cmd/versifyr@latest
```

## Getting Started

### 1. Initialize Configuration
Run the following command to create a default `.versifyr/configuration.yaml` in your project root:

```sh
versifyr init
```

### 2. Define Your Templates
You can define what needs to be updated in two ways:

#### Option A: Embedded in Source (Recommended)
Add a comment directly above the line you want to manage. The tool will replace the **immediately following line** with the result of your template.

**Go Example:**
```go
// $versifyr:template=const Version = "{{ .version }}"$
const Version = "v0.0.0"
```

**Maven Example (XML):**
```xml
<!--$versifyr:template=<version>{{ .version }}</version>$-->
<version>1.0.0-SNAPSHOT</version>
```

#### Option B: Centralized in Configuration
If you cannot add comments to a file, define the target row and template in `.versifyr/configuration.yaml`:

```yaml
files:
  - name: version.go
    type: go
    path: internal/versifyr/version.go
    templates:
      - row: 3
        template: 'const Version = "{{ .version }}"'
```

### 3. Apply Changes
Update your files by passing key-value pairs:

```sh
versifyr set version="v1.2.3"
```

## Built-in Variables

The following variables are always available in your templates:

- `version`: The primary value usually passed via CLI.
- `latesttag`: The most recent Git tag (e.g., `v0.1.0`).
- `actualdate`: Current date (YYYY-MM-DD).
- `actualtime`: Current time (HH:MM:SS).
- `actualtimestamp`: Current date and time.

## Advanced Usage

### Sprig Functions
You can use any Sprig function for transformations. For example, to generate a snake_case version for a constant:

```go
// $versifyr:template=const BUILD_ID = "{{ .version | replace "." "_" }}"$
const BUILD_ID = "v0_1_0"
```

### Escaping Quotes (JSON)
For formats like JSON where quotes must be escaped, use the `unescape` flag in your configuration:

```yaml
files:
  - name: package.json
    path: package.json
    type: json
    unescape: true
```

## Commands

| Command | Alias | Description |
| :--- | :--- | :--- |
| `init` | `i` | Creates the initial `.versifyr/configuration.yaml`. |
| `show` | `s` | Displays the current configuration and managed file content. |
| `set` | | Executes template replacements based on provided arguments. |

**Global Options:**
- `--debug`, `-d`: Enable verbose logging.
- `--nochange`, `-n`: Run in simulation mode (dry-run).

## Development

`versifyr` uses [Taskfile](https://taskfile.dev/) for a streamlined development experience.

### Local Tasks
- **Lint**: `task lint` - Runs `golangci-lint` exactly as it runs in CI.
- **Test**: `task test` - Runs the full test suite with coverage reporting.
- **Check Coverage**: `task test-coverage` - Verifies coverage against defined thresholds.
- **Advance Version**: `task advance-version [VERSION=vX.Y.Z]` - Increments the patch version from the latest Git tag.

### CI/CD Pipeline
Every push and pull request to `main` triggers a comprehensive pipeline:
1. **Linting**: Static analysis with `golangci-lint`.
2. **Testing**: Unit and integration tests with coverage enforcement.
3. **Release**: Automated multi-arch builds and GitHub Release creation (triggered on version tags).

## License
(c) 2023-2026 Stefano Zuccaro. All rights reserved.
