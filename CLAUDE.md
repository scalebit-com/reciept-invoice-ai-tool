# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a CLI tool for extracting structured information from receipt and invoice data in text or markdown files, outputting results as JSON. The project now has a complete implementation with OpenAI integration and comprehensive logging.

## Build System

The project uses [Task](https://taskfile.dev/) for build automation via `Taskfile.yaml`:

```bash
# Build the application (creates binary in target/)
task build

# Process all sample MD files and generate JSON outputs (depends on build)
task run

# Direct Go build alternative
go build -o target/reciept-invoice-ai-tool main.go
```

Build output goes to `target/` directory which is git-ignored.

## Current Architecture

- **Language**: Go 1.24.2
- **Module**: `github.com/scalebit-com/reciept-invoice-ai-tool`
- **Entry Point**: `main.go` → `cmd.Execute()`
- **CLI Framework**: Cobra with structured commands
- **AI Provider**: OpenAI with structured outputs
- **Logging**: Color-coded logging system inspired by getgmail
- **Build Target**: Single binary CLI tool

### Project Structure

```
├── cmd/                    # Cobra CLI commands
│   ├── root.go            # Root command and CLI setup
│   └── extract.go         # Extract command implementation
├── pkg/
│   ├── interfaces/        # Interface definitions
│   │   ├── logger.go      # Logger interface
│   │   └── ai_provider.go # AI provider interface and data structures
│   ├── logger/           # Logging implementation
│   │   └── logger.go     # ColorLogger with timestamped output
│   ├── ai/               # AI provider implementations
│   │   └── openai_provider.go # OpenAI provider with structured outputs
│   └── config/           # Configuration management
│       └── config.go     # Generic configuration (provider-agnostic)
├── sampledata/           # Sample receipt/invoice files and extracted JSON
├── target/              # Build output (git-ignored)
├── main.go              # Application entry point
├── Taskfile.yaml        # Build automation
├── .env                 # Environment variables (git-ignored)
└── version.txt          # Current version: 1.3.0
```

## CLI Usage

The tool implements a Cobra-based CLI with the `extract` command:

```bash
# Extract from receipt/invoice file with output file (both required)
./target/reciept-invoice-ai-tool extract -i path/to/receipt.md -o output.json

# Show help
./target/reciept-invoice-ai-tool --help
./target/reciept-invoice-ai-tool extract --help
```

### Command Flags

- `-i, --input` (required): Path to the input file
- `-o, --output` (required): Path to the output JSON file

### File Validation

The `extract` command performs comprehensive validation:
- ✅ **File existence** - errors and exits if file doesn't exist
- ✅ **Binary detection** - errors and exits if file is binary (with tolerance for occasional null bytes)
- ✅ **Size limits** - errors and exits if file > 200KB
- ⚠️ **Extension check** - warns for non-.txt/.md files but continues

## Environment Configuration

### OpenAI Provider Configuration

The OpenAI provider handles its own environment variables internally:

- `OPENAI_KEY` (required): OpenAI API key
- `OPENAI_MODEL` (optional): Model to use, defaults to `gpt-4o-2024-08-06` with warning if not set

Configuration is loaded from `.env` file in the current directory or from environment variables.

### Example .env file:
```
OPENAI_KEY=sk-proj-your-actual-api-key-here
OPENAI_MODEL=gpt-4o-2024-08-06
```

## Logging System

Implements colored, timestamped logging inspired by getgmail:
- **INFO**: Green - progress updates, successful operations
- **ERROR**: Red - operation failures with context
- **WARN**: Yellow - non-critical issues (e.g., missing OPENAI_MODEL)
- **DEBUG**: Cyan - detailed debugging information

Example output:
```
[2025-08-03 21:03:43] [INFO] Starting receipt/invoice extraction for file: receipt.md
[2025-08-03 21:03:43] [INFO] File validation successful
[2025-08-03 21:03:43] [INFO] Output will be written to: output.json
[2025-08-03 21:03:43] [INFO] Using OpenAI model: gpt-4o-2024-08-06
[2025-08-03 21:03:43] [INFO] OpenAI API call successful (took 4.88s)
[2025-08-03 21:03:43] [INFO] Token usage - Prompt: 765, Completion: 48, Total: 813
[2025-08-03 21:03:43] [INFO] Extracted document type: Receipt
[2025-08-03 21:03:43] [INFO] Extracted company: Anthropic, PBC
[2025-08-03 21:03:43] [INFO] Successfully wrote JSON output to output.json
```

## Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/fatih/color` - Terminal color output
- `github.com/openai/openai-go` - OpenAI API client
- `github.com/invopop/jsonschema` - JSON schema generation for structured outputs
- `github.com/joho/godotenv` - Environment variable loading from .env files

## Data Extraction

The tool extracts structured information from receipts and invoices:

### Output JSON Structure
```json
{
  "document_type": "Receipt",
  "description": "AI Services",
  "company": "Anthropic, PBC",
  "date_issued": "2025-08-02",
  "service_description": "Max plan - 5x subscription",
  "se_cent_amount": 109677
}
```

### Field Descriptions
- `document_type`: Always present - "None" (not financial), "Invoice", or "Receipt"
- `description`: **Mandatory** - Accountant-friendly categorization (max 50 chars):
  - For "None" documents: describes what the document is about
  - For "Invoice/Receipt": generic service category (e.g., "AI Services", "Cloud Services")
  - Generated by analyzing entire document: headers, company name, service details, context
- `company`: Optional - The company offering the service and requesting payment
- `date_issued`: Optional - Date in YYYY-MM-DD format
- `service_description`: Optional - Description of services or items
- `se_cent_amount`: Optional - Amount in Swedish cents (öre), where last 2 digits are cents

### Currency Handling
- SEK amounts: multiply by 100 (95.37 SEK = 9537 öre)
- EUR amounts: convert using approximate rate (1 EUR ≈ 11.5 SEK), then to öre
- Other currencies: convert to SEK first, then to öre

## AI Provider Architecture

The system uses a provider pattern for AI services:

### AIProvider Interface
```go
type AIProvider interface {
    GetReceiptInvoiceInfo(content string) (*ReceiptInvoiceInfo, error)
}
```

### OpenAI Implementation
- Uses structured outputs with JSON schema validation
- Handles environment variable configuration internally
- Provides detailed logging of API interactions
- Supports configurable models via `OPENAI_MODEL`

This architecture allows easy addition of other AI providers (Anthropic Claude, local models, etc.) without changing the core application logic.

## Development Status

- ✅ CLI framework with Cobra
- ✅ Colored logging system with comprehensive AI interaction logging
- ✅ File validation and error handling
- ✅ Command structure and help documentation
- ✅ OpenAI integration with structured outputs
- ✅ JSON output to files and console
- ✅ Environment variable configuration
- ✅ Provider pattern for AI services
- ✅ Company extraction functionality
- ✅ Currency conversion to Swedish cents
- ✅ Document type classification (None/Invoice/Receipt)