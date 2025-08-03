package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/scalebit-com/reciept-invoice-ai-tool/pkg/interfaces"
)

//go:embed overview-template.html
var htmlTemplateContent string

// htmloverviewCmd represents the htmloverview command
var htmloverviewCmd = &cobra.Command{
	Use:   "htmloverview",
	Short: "Generate HTML overview from extracted JSON data",
	Long: `Generate a professional HTML overview report from JSON data extracted by the extract command.
The HTML output is optimized for printing and includes all extracted information in a nicely formatted layout.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile, _ := cmd.Flags().GetString("input")
		outputFile, _ := cmd.Flags().GetString("output")
		return runHTMLOverview(inputFile, outputFile, logger)
	},
}

func init() {
	rootCmd.AddCommand(htmloverviewCmd)
	htmloverviewCmd.Flags().StringP("input", "i", "", "Path to the input JSON file (required)")
	htmloverviewCmd.Flags().StringP("output", "o", "", "Path to the output HTML file (required)")
	htmloverviewCmd.MarkFlagRequired("input")
	htmloverviewCmd.MarkFlagRequired("output")
}

// TemplateData represents the data structure passed to the HTML template
type TemplateData struct {
	Data        *interfaces.ReceiptInvoiceInfo `json:"data"`
	ProcessedAt string                         `json:"processed_at"`
}

// runHTMLOverview handles the htmloverview command logic
func runHTMLOverview(inputFile string, outputFile string, log interfaces.Logger) error {
	log.Info("Starting HTML overview generation for file: %s", inputFile)

	// Check if output file already exists
	if _, err := os.Stat(outputFile); err == nil {
		log.Warn("Output file already exists: %s", outputFile)
		return nil
	}

	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		log.Error("Input file does not exist: %s", inputFile)
		return fmt.Errorf("input file does not exist: %s", inputFile)
	}

	log.Info("Reading JSON data from: %s", inputFile)

	// Read JSON file content
	jsonContent, err := os.ReadFile(inputFile)
	if err != nil {
		log.Error("Failed to read input file %s: %v", inputFile, err)
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Parse JSON data
	var receiptData interfaces.ReceiptInvoiceInfo
	if err := json.Unmarshal(jsonContent, &receiptData); err != nil {
		log.Error("Failed to parse JSON from %s: %v", inputFile, err)
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	log.Info("Successfully parsed JSON data")
	log.Info("Document type: %s", receiptData.DocumentType)
	log.Info("Description: %s", receiptData.Description)
	if receiptData.Company != nil {
		log.Info("Company: %s", *receiptData.Company)
	}

	// Prepare template data
	templateData := TemplateData{
		Data:        &receiptData,
		ProcessedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	log.Info("Generating HTML output...")

	// Create template functions
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"default": func(defaultValue, value string) string {
			if value == "" {
				return defaultValue
			}
			return value
		},
		"div": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"divInt": func(a *int, b float64) float64 {
			if a == nil || b == 0 {
				return 0
			}
			return float64(*a) / b
		},
		"formatFloat": func(f *float64) string {
			if f == nil {
				return "0.00"
			}
			return fmt.Sprintf("%.2f", *f)
		},
		"formatFloatWithCurrency": func(f *float64, currency *string) string {
			if f == nil {
				return "0.00"
			}
			if currency == nil || *currency == "" {
				return fmt.Sprintf("%.2f", *f)
			}
			return fmt.Sprintf("%.2f %s", *f, *currency)
		},
	}

	// Parse template
	tmpl, err := template.New("overview").Funcs(funcMap).Parse(htmlTemplateContent)
	if err != nil {
		log.Error("Failed to parse HTML template: %v", err)
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	// Create output file
	outputFileHandle, err := os.Create(outputFile)
	if err != nil {
		log.Error("Failed to create output file %s: %v", outputFile, err)
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFileHandle.Close()

	// Execute template
	if err := tmpl.Execute(outputFileHandle, templateData); err != nil {
		log.Error("Failed to execute HTML template: %v", err)
		return fmt.Errorf("failed to execute HTML template: %w", err)
	}

	log.Info("Successfully generated HTML overview: %s", outputFile)

	// Also log some statistics
	if receiptData.SECentAmount != nil {
		log.Info("Amount (SEK): %.2f SEK (%d öre)", float64(*receiptData.SECentAmount)/100.0, *receiptData.SECentAmount)
	}
	if receiptData.OriginalAmount != nil {
		currency := "units"
		if receiptData.OriginalCurrency != nil {
			currency = *receiptData.OriginalCurrency
		}
		log.Info("Original amount: %.2f %s", *receiptData.OriginalAmount, currency)
	}
	if len(receiptData.IdFields) > 0 {
		log.Info("Found %d identification fields", len(receiptData.IdFields))
	}

	return nil
}