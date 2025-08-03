# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a CLI tool for extracting structured information from receipt and invoice data in text or markdown files, outputting results as JSON. The project now has a complete implementation with OpenAI integration and comprehensive logging.

## Build System

The project uses [Task](https://taskfile.dev/) for build automation via `Taskfile.yaml`:

```bash
# Build the application (creates binary in target/)
task build

# Process all sample MD files and generate JSON and HTML outputs (depends on build)
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
- **HTML Reporting**: Embedded Go templates for professional HTML output

### Project Structure

```
├── cmd/                    # Cobra CLI commands
│   ├── root.go            # Root command and CLI setup
│   ├── extract.go         # Extract command implementation
│   ├── htmloverview.go    # HTML overview generation command
│   └── overview-template.html # HTML template (embedded in binary)
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
└── version.txt          # Current version: 2.1.0
```

## CLI Usage

The tool implements a Cobra-based CLI with the `extract` and `htmloverview` commands:

```bash
# Extract from receipt/invoice file with output file (both required)
./target/reciept-invoice-ai-tool extract -i path/to/receipt.md -o output.json

# Generate HTML overview from JSON file (both required)
./target/reciept-invoice-ai-tool htmloverview -i output.json -o overview.html

# Show help
./target/reciept-invoice-ai-tool --help
./target/reciept-invoice-ai-tool extract --help
./target/reciept-invoice-ai-tool htmloverview --help
```

### Command Flags

**Extract Command:**
- `-i, --input` (required): Path to the input file
- `-o, --output` (required): Path to the output JSON file

**HTML Overview Command:**
- `-i, --input` (required): Path to the input JSON file
- `-o, --output` (required): Path to the output HTML file

### File Validation

Both commands perform comprehensive validation and file existence checks:

**Extract Command Validation:**
- ✅ **Output file existence** - warns and exits gracefully if output file already exists
- ✅ **Input file existence** - errors and exits if input file doesn't exist
- ✅ **Binary detection** - errors and exits if file is binary (with tolerance for occasional null bytes)
- ✅ **Size limits** - errors and exits if file > 200KB
- ⚠️ **Extension check** - warns for non-.txt/.md files but continues

**HTML Overview Command Validation:**
- ✅ **Output file existence** - warns and exits gracefully if output file already exists
- ✅ **Input file existence** - errors and exits if input JSON file doesn't exist
- ✅ **JSON validation** - errors and exits if input file is not valid JSON

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
  "se_cent_amount": 109677,
  "original_amount": 95.37,
  "original_currency": "EUR",
  "original_vat_amount": 19.07,
  "id_fields": [
    {
      "name": "Receipt Number",
      "value": "2844-5789-6006"
    },
    {
      "name": "Invoice Number", 
      "value": "D8F78A38-0007"
    }
  ],
  "suggested_filename": "2025_08_02-anthropic__pbc-ai_services-1097sek"
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
- `original_amount`: Optional - Total amount in original currency as it appears in document
- `original_currency`: Optional - ISO 3-letter currency code (e.g., "EUR", "USD", "SEK")
- `original_vat_amount`: Optional - VAT/tax amount in original currency
- `id_fields`: Optional list of identification fields found in document:
  - Each entry has `name` (identifier type) and `value` (actual identifier)
  - Examples: Invoice Number, Receipt Number, Customer ID, Order Number
- `suggested_filename`: **Auto-generated** - Filesystem-safe filename suggestion based on extracted data:
  - Format: `<date>-<company>-<description>-<amount>SEK`
  - All lowercase with non-alphanumeric characters replaced with `_`
  - Amount converted from öre to SEK rounded to nearest krona
  - Missing fields default to "unknown"
  - Example: `2025_08_02-anthropic__pbc-ai_services-1097sek`

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

## HTML Overview Generation

The `htmloverview` command generates professional HTML reports from JSON output:

### Features
- **Professional Design**: Clean, modern layout with embedded CSS
- **Print Optimized**: Styles optimized for printing with proper page breaks
- **Responsive**: Mobile-friendly design that adapts to different screen sizes
- **Comprehensive Data Display**: Shows all extracted information in organized sections
- **Process Timestamp**: Includes generation date and time in the footer
- **Template Embedding**: HTML template is embedded in the binary for easy deployment

### Template Structure
- **Document Information**: Type, description, company, date
- **Financial Information**: Amounts in original currency and Swedish kronor
- **Identification Fields**: Table of all found ID fields
- **Process Information**: Generation timestamp and tool information

### Usage
```bash
# Generate HTML from existing JSON
./target/reciept-invoice-ai-tool htmloverview -i receipt.json -o receipt.html

# The task run command automatically generates both JSON and HTML files
task run
```

The HTML template (`cmd/overview-template.html`) is embedded in the binary using Go's `//go:embed` directive, ensuring the tool remains a single, deployable binary.

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
- ✅ Original amount and currency extraction
- ✅ VAT amount extraction in original currency
- ✅ ID field extraction (invoice numbers, receipt numbers, etc.)
- ✅ HTML overview generation with embedded templates
- ✅ Professional print-optimized HTML reports
- ✅ Automated HTML generation in build pipeline
- ✅ File existence protection with graceful warnings
- ✅ Git ignore patterns for generated files
- ✅ SuggestedFileName field with auto-generated filesystem-safe filenames