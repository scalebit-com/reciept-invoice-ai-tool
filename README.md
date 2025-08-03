# Receipt Invoice AI Tool

A CLI tool to extract structured information from text files or markdown files containing receipt data and output it as structured JSON.

## Features

- 📄 Extract data from text (.txt) and markdown (.md) files
- 🤖 AI-powered parsing of receipt and invoice information (planned)
- 📊 Output structured JSON format (planned)
- ⚡ Fast CLI interface with Cobra framework
- 🎨 Colored logging with timestamps
- 🔧 Easy to build and deploy
- ✅ Comprehensive file validation (existence, binary detection, size limits)

## Installation

### Prerequisites

- Go 1.24.2 or higher
- [Task](https://taskfile.dev/) (optional, for development)

### Build from Source

1. Clone the repository:
```bash
git clone https://github.com/scalebit-com/reciept-invoice-ai-tool.git
cd reciept-invoice-ai-tool
```

2. Build the application:
```bash
# Using Task (recommended)
task build

# Or using Go directly
go build -o target/reciept-invoice-ai-tool main.go
```

3. The binary will be available in the `target/` directory.

## Usage

The tool uses a Cobra-based CLI with structured commands:

```bash
# Show help
./target/reciept-invoice-ai-tool --help

# Extract from receipt/invoice file
./target/reciept-invoice-ai-tool extract -i <input-file>
```

### Basic Examples

```bash
# Process a receipt text file
./target/reciept-invoice-ai-tool extract -i receipt.txt

# Process a markdown file with receipt data
./target/reciept-invoice-ai-tool extract -i invoice.md

# Show extract command help
./target/reciept-invoice-ai-tool extract --help
```

### File Validation

The tool performs comprehensive validation on input files:
- ✅ **File existence** - errors and exits if file doesn't exist
- ✅ **Binary detection** - errors and exits if file is binary
- ✅ **Size limits** - errors and exits if file > 200KB
- ⚠️ **Extension check** - warns for non-.txt/.md files but continues

## Input Format

The tool accepts text and markdown files containing receipt or invoice information. Examples:

### Text File Format
```
Store: Best Buy
Date: 2024-01-15
Total: $299.99

Items:
- iPhone Cable - $19.99
- Phone Case - $24.99
- Screen Protector - $14.99
- Tax: $24.00
```

### Markdown Format
```markdown
# Receipt - Electronics Store

**Date:** 2024-01-15  
**Store:** Best Buy  
**Location:** 123 Main St, City, State  

## Items
| Item | Quantity | Price |
|------|----------|-------|
| iPhone Cable | 1 | $19.99 |
| Phone Case | 1 | $24.99 |
| Screen Protector | 1 | $14.99 |

**Subtotal:** $59.97  
**Tax:** $4.80  
**Total:** $64.77
```

## Output Format

**Note: JSON output functionality is planned for future implementation.**

The tool will output structured JSON containing extracted information:

```json
{
  "store_name": "Best Buy",
  "date": "2024-01-15",
  "location": "123 Main St, City, State",
  "items": [
    {
      "name": "iPhone Cable",
      "quantity": 1,
      "price": 19.99
    },
    {
      "name": "Phone Case", 
      "quantity": 1,
      "price": 24.99
    },
    {
      "name": "Screen Protector",
      "quantity": 1,
      "price": 14.99
    }
  ],
  "subtotal": 59.97,
  "tax": 4.80,
  "total": 64.77,
  "currency": "USD"
}
```

## Current Implementation Status

- ✅ **CLI Framework** - Complete Cobra-based command structure
- ✅ **Logging System** - Color-coded, timestamped logging inspired by getgmail
- ✅ **File Validation** - Comprehensive input file validation
- ✅ **Error Handling** - Proper error handling and user feedback
- ⏳ **Extraction Logic** - Planned for future implementation
- ⏳ **JSON Output** - Planned for future implementation
- ⏳ **AI Integration** - Planned for future implementation

## Development

### Available Tasks

```bash
# Build the application
task build

# Run tests (when implemented)
task test

# Clean build artifacts
task clean
```

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
├── Taskfile.yaml        # Build automation
├── go.mod               # Go module dependencies
└── README.md
```

### Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/fatih/color` - Terminal color output

### Building

The project uses [Task](https://taskfile.dev/) for build automation. The main tasks are defined in `Taskfile.yaml`:

- `task build`: Compiles the Go application to `target/reciept-invoice-ai-tool`

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