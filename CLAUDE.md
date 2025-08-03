# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a CLI tool for extracting structured information from receipt and invoice data in text or markdown files, outputting results as JSON. The project now has a complete CLI framework with logging and file validation implemented.

## Build System

The project uses [Task](https://taskfile.dev/) for build automation via `Taskfile.yaml`:

```bash
# Build the application (creates binary in target/)
task build

# Direct Go build alternative
go build -o target/reciept-invoice-ai-tool main.go
```

Build output goes to `target/` directory which is git-ignored.

## Current Architecture

- **Language**: Go 1.24.2
- **Module**: `github.com/scalebit-com/reciept-invoice-ai-tool`
- **Entry Point**: `main.go` → `cmd.Execute()`
- **CLI Framework**: Cobra with structured commands
- **Logging**: Color-coded logging system inspired by getgmail
- **Build Target**: Single binary CLI tool

### Project Structure

```
├── cmd/                    # Cobra CLI commands
│   ├── root.go            # Root command and CLI setup
│   └── extract.go         # Extract command implementation
├── pkg/
│   ├── interfaces/        # Interface definitions
│   │   └── logger.go      # Logger interface
│   └── logger/           # Logging implementation
│       └── logger.go     # ColorLogger with timestamped output
├── sampledata/           # Sample receipt/invoice files
├── target/              # Build output (git-ignored)
├── main.go              # Application entry point
└── Taskfile.yaml        # Build automation
```

## CLI Usage

The tool implements a Cobra-based CLI with the `extract` command:

```bash
# Extract from receipt/invoice file
./target/reciept-invoice-ai-tool extract -i path/to/receipt.md

# Show help
./target/reciept-invoice-ai-tool --help
./target/reciept-invoice-ai-tool extract --help
```

### File Validation

The `extract` command performs comprehensive validation:
- ✅ **File existence** - errors and exits if file doesn't exist
- ✅ **Binary detection** - errors and exits if file is binary
- ✅ **Size limits** - errors and exits if file > 200KB
- ⚠️ **Extension check** - warns for non-.txt/.md files but continues

## Logging System

Implements colored, timestamped logging inspired by getgmail:
- **INFO**: Green - progress updates, successful operations
- **ERROR**: Red - operation failures with context
- **WARN**: Yellow - non-critical issues
- **DEBUG**: Cyan - detailed debugging information

Example output:
```
[2025-08-03 11:09:32] [INFO] Starting receipt/invoice extraction for file: receipt.md
[2025-08-03 11:09:32] [INFO] File validation successful
```

## Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/fatih/color` - Terminal color output

## Expected Data Flow

The tool framework is ready for:

1. Accept text (.txt) or markdown (.md) files containing receipt data
2. Parse receipt information (store name, date, items, prices, totals)
3. Output structured JSON with fields like:
   - `store_name`, `date`, `location`
   - `items[]` with `name`, `quantity`, `price`
   - `subtotal`, `tax`, `total`, `currency`

## Development Status

- ✅ CLI framework with Cobra
- ✅ Colored logging system
- ✅ File validation and error handling
- ✅ Command structure and help documentation
- ⏳ **TODO**: Actual extraction/parsing logic implementation
- ⏳ **TODO**: JSON output formatting
- ⏳ **TODO**: AI/LLM integration for content parsing