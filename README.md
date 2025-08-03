# Receipt Invoice AI Tool

A CLI tool to extract structured information from text files or markdown files containing receipt and invoice data, powered by OpenAI, and output it as structured JSON.

## Features

- 📄 Extract data from text (.txt) and markdown (.md) files
- 🤖 AI-powered parsing using OpenAI with structured outputs
- 📊 Output structured JSON format with document classification
- 🌐 Professional HTML report generation with embedded templates
- ⚡ Fast CLI interface with Cobra framework
- 🎨 Colored logging with timestamps and detailed AI interaction logs
- 🔧 Easy to build and deploy
- ✅ Comprehensive file validation (existence, binary detection, size limits)
- 💰 Currency conversion to Swedish cents (öre)
- 🏢 Company extraction from financial documents
- 💵 Original amount and currency preservation
- 🧾 VAT amount extraction in original currency
- 🔢 ID field extraction (invoice numbers, receipt numbers, etc.)
- 📝 Mandatory output file specification
- 🖨️ Print-optimized HTML reports with professional styling
- 📄 Auto-generated filesystem-safe filename suggestions
- 🐳 Docker support with multi-stage builds (42.8MB image)
- 📦 Container registry integration

## Installation

### Prerequisites

- Go 1.24.2 or higher (for building from source)
- [Task](https://taskfile.dev/) (optional, for development)
- Docker (for containerized usage)
- OpenAI API key

### Build from Source

1. Clone the repository:
```bash
git clone https://github.com/scalebit-com/reciept-invoice-ai-tool.git
cd reciept-invoice-ai-tool
```

2. Set up environment variables:
```bash
# Create .env file with your OpenAI API key
echo "OPENAI_KEY=sk-proj-your-actual-api-key-here" > .env
echo "OPENAI_MODEL=gpt-4o-2024-08-06" >> .env
```

3. Build the application:
```bash
# Using Task (recommended)
task build

# Process all sample MD files and generate JSON and HTML outputs (depends on build)
task run

# Or using Go directly
go build -o target/reciept-invoice-ai-tool main.go
```

4. The binary will be available in the `target/` directory.

### Using Docker

You can use the pre-built Docker image without installing Go:

```bash
# Pull the latest image
docker pull perarneng/reciept-invoice-ai-tool:latest

# Or use a specific version
docker pull perarneng/reciept-invoice-ai-tool:2.1.0
```

## Usage

The tool uses a Cobra-based CLI with structured commands:

```bash
# Show help
./target/reciept-invoice-ai-tool --help

# Extract from receipt/invoice file (both input and output files are required)
./target/reciept-invoice-ai-tool extract -i <input-file> -o <output-file>

# Generate HTML overview from JSON file (both input and output files are required)
./target/reciept-invoice-ai-tool htmloverview -i <json-file> -o <html-file>
```

### Basic Examples

```bash
# Process a receipt text file
./target/reciept-invoice-ai-tool extract -i receipt.txt -o receipt.json

# Process a markdown file with receipt data
./target/reciept-invoice-ai-tool extract -i invoice.md -o invoice.json

# Generate HTML overview from JSON
./target/reciept-invoice-ai-tool htmloverview -i receipt.json -o receipt.html

# Show command help
./target/reciept-invoice-ai-tool extract --help
./target/reciept-invoice-ai-tool htmloverview --help
```

### Docker Usage

The same commands work with Docker. You need to mount your working directory and pass environment variables:

```bash
# Show help
docker run --rm perarneng/reciept-invoice-ai-tool:latest --help

# Extract from receipt/invoice file using Docker
docker run --rm -v $(pwd):/app/data \
  -e OPENAI_KEY="your-api-key" \
  perarneng/reciept-invoice-ai-tool:latest extract \
  -i /app/data/receipt.md -o /app/data/output.json

# Generate HTML overview using Docker
docker run --rm -v $(pwd):/app/data \
  perarneng/reciept-invoice-ai-tool:latest htmloverview \
  -i /app/data/output.json -o /app/data/overview.html
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

### Required Environment Variables

- `OPENAI_KEY`: Your OpenAI API key (required)
- `OPENAI_MODEL`: Model to use (optional, defaults to `gpt-4o-2024-08-06`)

### Configuration Methods

1. **Using .env file** (recommended):
```bash
# Create .env file in the project directory
echo "OPENAI_KEY=sk-proj-your-actual-api-key-here" > .env
echo "OPENAI_MODEL=gpt-4o-2024-08-06" >> .env
```

2. **Using environment variables**:
```bash
export OPENAI_KEY="sk-proj-your-actual-api-key-here"
export OPENAI_MODEL="gpt-4o-2024-08-06"
```

## Input Format

The tool accepts text and markdown files containing receipt or invoice information. Examples:

### Text File Format
```
Receipt

Invoice number D9F68A38-0009

Date paid
August 2, 2025

Anthropic, PBC
548 Market Street
San Francisco, California 94104

€95.37 paid on August 2, 2025

Description
Max plan - 5x
Aug 2 – Sep 2, 2025
```

### Markdown Format
```markdown
# Receipt - Anthropic, PBC

**Date:** August 2, 2025  
**Invoice:** D9F68A38-0009  
**Company:** Anthropic, PBC  

## Service Description
Max plan - 5x subscription for Aug 2 – Sep 2, 2025

**Total:** €95.37
```

## Output Format

The tool outputs structured JSON containing extracted information:

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

- **`document_type`**: Always present - `"None"` (not financial), `"Invoice"`, or `"Receipt"`
- **`description`**: **Mandatory** - Accountant-friendly categorization (max 50 chars):
  - For "None" documents: describes what the document is about
  - For "Invoice/Receipt": generic service category (e.g., "AI Services", "Cloud Services")
  - Generated by analyzing entire document: headers, company name, service details, context
- **`company`**: Optional - The company offering the service and requesting payment
- **`date_issued`**: Optional - Date in YYYY-MM-DD format
- **`service_description`**: Optional - Description of services or items
- **`se_cent_amount`**: Optional - Amount in Swedish cents (öre), where last 2 digits are cents
- **`original_amount`**: Optional - Total amount in original currency as it appears in document
- **`original_currency`**: Optional - ISO 3-letter currency code (e.g., "EUR", "USD", "SEK")
- **`original_vat_amount`**: Optional - VAT/tax amount in original currency
- **`id_fields`**: Optional list of identification fields found in document:
  - Each entry has `name` (identifier type) and `value` (actual identifier)
  - Examples: Invoice Number, Receipt Number, Customer ID, Order Number
- **`suggested_filename`**: **Auto-generated** - Filesystem-safe filename suggestion based on extracted data:
  - Format: `<date>-<company>-<description>-<amount>SEK`
  - All lowercase with non-alphanumeric characters replaced with `_`
  - Amount converted from öre to SEK rounded to nearest krona
  - Missing fields default to "unknown"
  - Example: `2025_08_02-anthropic__pbc-ai_services-1097sek`

### Currency Handling

The tool converts all currencies to Swedish cents (öre):
- **SEK amounts**: multiply by 100 (95.37 SEK = 9537 öre)
- **EUR amounts**: convert using approximate rate (1 EUR ≈ 11.5 SEK), then to öre
- **Other currencies**: convert to SEK first, then to öre

## HTML Overview Generation

The `htmloverview` command generates professional HTML reports from JSON output created by the `extract` command.

### Features

- **Professional Design**: Clean, formal layout with Arial/Helvetica fonts and minimal styling
- **Print Optimization**: A4-optimized CSS for perfect printing with proper page breaks
- **Comprehensive Data Display**: Shows all extracted information in organized sections:
  - Document information (type, description, company, date)
  - Financial information with currency conversion display
  - Identification fields in a professional table format
- **Process Timestamp**: Includes generation date and time in the footer
- **Embedded Template**: HTML template is embedded in the binary for single-file deployment
- **Responsive Design**: Mobile-friendly layout that adapts to different screen sizes

### Usage Examples

```bash
# Generate HTML from existing JSON
./target/reciept-invoice-ai-tool htmloverview -i receipt.json -o receipt.html

# The task run command automatically generates both JSON and HTML files
task run
```

### HTML Template

The HTML template (`cmd/overview-template.html`) includes:
- Clean CSS with embedded styles for offline use
- Formal black and white design optimized for business use
- A4-specific media queries for optimal printing with Puppeteer
- Responsive design for mobile devices
- Arial/Helvetica typography with minimal borders

The template is embedded in the binary using Go's `//go:embed` directive, ensuring the tool remains a single, deployable binary without external dependencies.

## Logging

The tool provides comprehensive logging with colored, timestamped output:

```
[2025-08-03 21:03:43] [INFO] Starting receipt/invoice extraction for file: receipt.md
[2025-08-03 21:03:43] [INFO] File validation successful
[2025-08-03 21:03:43] [INFO] Output will be written to: output.json
[2025-08-03 21:03:43] [INFO] Using OpenAI model: gpt-4o-2024-08-06
[2025-08-03 21:03:43] [INFO] OpenAI API call successful (took 4.88s)
[2025-08-03 21:03:43] [INFO] Token usage - Prompt: 765, Completion: 48, Total: 813
[2025-08-03 21:03:43] [INFO] Extracted document type: Receipt
[2025-08-03 21:03:43] [INFO] Extracted description: AI Services
[2025-08-03 21:03:43] [INFO] Extracted company: Anthropic, PBC
[2025-08-03 21:03:43] [INFO] Successfully wrote JSON output to output.json
```

### Log Levels
- **INFO**: Green - progress updates, successful operations
- **ERROR**: Red - operation failures with context
- **WARN**: Yellow - non-critical issues (e.g., missing OPENAI_MODEL)
- **DEBUG**: Cyan - detailed debugging information

## Architecture

### AI Provider Pattern

The system uses a provider pattern for AI services, making it easy to add new providers:

```go
type AIProvider interface {
    GetReceiptInvoiceInfo(content string) (*ReceiptInvoiceInfo, error)
}
```

Current implementation:
- **OpenAI Provider**: Uses structured outputs with JSON schema validation
- **Future providers**: Could include Anthropic Claude, local models, etc.

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

## Development

### Available Tasks

```bash
# Build the application
task build

# Process all sample MD files and generate JSON and HTML outputs
task run

# Docker tasks
task docker-build    # Build Docker image with version tags
task docker-push     # Push to registry (both latest and version tags)
task docker-run      # Process all sample files using Docker (equivalent to task run)

# Run tests (when implemented)
task test

# Clean build artifacts
task clean
```

### Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/fatih/color` - Terminal color output
- `github.com/openai/openai-go` - OpenAI API client
- `github.com/invopop/jsonschema` - JSON schema generation for structured outputs
- `github.com/joho/godotenv` - Environment variable loading from .env files

### Building

The project uses [Task](https://taskfile.dev/) for build automation. The main tasks are defined in `Taskfile.yaml`:

- `task build`: Compiles the Go application to `target/reciept-invoice-ai-tool`

## Implementation Status

- ✅ **CLI Framework** - Complete Cobra-based command structure
- ✅ **Logging System** - Color-coded, timestamped logging with AI interaction details
- ✅ **File Validation** - Comprehensive input file validation
- ✅ **Error Handling** - Proper error handling and user feedback
- ✅ **OpenAI Integration** - Structured outputs with JSON schema validation
- ✅ **JSON Output** - Output to both console and specified file
- ✅ **Environment Configuration** - .env file support and environment variables
- ✅ **Document Classification** - Automatic classification of document types
- ✅ **Company Extraction** - Extract company information from financial documents
- ✅ **Currency Conversion** - Convert all currencies to Swedish cents (öre)
- ✅ **Original Amount Preservation** - Extract and preserve original amounts and currencies
- ✅ **VAT Amount Extraction** - Extract VAT/tax amounts in original currency
- ✅ **ID Field Extraction** - Extract identification fields (invoice numbers, receipt numbers, etc.)
- ✅ **Provider Pattern** - Extensible architecture for multiple AI providers
- ✅ **HTML Report Generation** - Professional, print-optimized HTML reports
- ✅ **Embedded Templates** - Single binary deployment with embedded HTML templates
- ✅ **Automated Pipeline** - Build process generates both JSON and HTML outputs
- ✅ **File Existence Protection** - Graceful warnings when output files already exist
- ✅ **Git Ignore Patterns** - Generated files are properly excluded from version control
- ✅ **SuggestedFileName Field** - Auto-generated filesystem-safe filename suggestions
- ✅ **Docker Support** - Multi-stage builds with minimal 42.8MB image size
- ✅ **Container Registry** - Published to `perarneng/reciept-invoice-ai-tool`
- ✅ **Docker Build Automation** - Task-based Docker build and push workflows

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the terms found in the [LICENSE](LICENSE) file.

## Support

For issues and feature requests, please use the [GitHub Issues](https://github.com/scalebit-com/reciept-invoice-ai-tool/issues) page.