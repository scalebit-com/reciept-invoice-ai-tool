# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a CLI tool for extracting structured information from receipt and invoice data in text or markdown files, outputting results as JSON. The project is in early development with a basic Go structure established.

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
- **Entry Point**: `main.go` (currently prints "Hello World")
- **Build Target**: Single binary CLI tool

## Expected Data Flow

Based on README specifications, the tool should:

1. Accept text (.txt) or markdown (.md) files containing receipt data
2. Parse receipt information (store name, date, items, prices, totals)
3. Output structured JSON with fields like:
   - `store_name`, `date`, `location`
   - `items[]` with `name`, `quantity`, `price`
   - `subtotal`, `tax`, `total`, `currency`

## Input Format Examples

Text format:
```
Store: Best Buy
Date: 2024-01-15
Total: $299.99
Items:
- iPhone Cable - $19.99
```

Markdown format with tables and structured headers.

## Development Notes

- Current state: Basic Go project scaffold with Hello World
- No CLI argument parsing implemented yet
- No file I/O or JSON output implemented yet
- No AI/parsing logic implemented yet
- Task file mentions `task test` and `task clean` but these tasks don't exist yet