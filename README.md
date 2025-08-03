# Receipt Invoice AI Tool

A CLI tool to extract structured information from text files or markdown files containing receipt data and output it as structured JSON.

## Features

- 📄 Extract data from text (.txt) and markdown (.md) files
- 🤖 AI-powered parsing of receipt and invoice information
- 📊 Output structured JSON format
- ⚡ Fast CLI interface
- 🔧 Easy to build and deploy

## Installation

### Prerequisites

- Go 1.19 or higher
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

```bash
./target/reciept-invoice-ai-tool [options] <input-file>
```

### Basic Example

```bash
# Process a receipt text file
./target/reciept-invoice-ai-tool receipt.txt

# Process a markdown file with receipt data
./target/reciept-invoice-ai-tool invoice.md
```

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

The tool outputs structured JSON containing extracted information:

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
.
├── README.md
├── LICENSE
├── go.mod
├── main.go
├── Taskfile.yaml
├── .gitignore
└── target/          # Build output (git-ignored)
```

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