package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/scalebit-com/reciept-invoice-ai-tool/pkg/ai"
	"github.com/scalebit-com/reciept-invoice-ai-tool/pkg/interfaces"
)

const maxFileSize = 200 * 1024 // 200KB in bytes

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract structured information from receipt and invoice files",
	Long: `Extract structured information from receipt and invoice data in text or markdown files.
The tool will parse the file and output structured JSON with receipt/invoice details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile, _ := cmd.Flags().GetString("input")
		outputFile, _ := cmd.Flags().GetString("output")
		return runExtract(inputFile, outputFile, logger)
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
	extractCmd.Flags().StringP("input", "i", "", "Path to the input file (required)")
	extractCmd.Flags().StringP("output", "o", "", "Path to the output JSON file (required)")
	extractCmd.MarkFlagRequired("input")
	extractCmd.MarkFlagRequired("output")
}

// runExtract handles the extract command logic
func runExtract(inputFile string, outputFile string, log interfaces.Logger) error {
	log.Info("Starting receipt/invoice extraction for file: %s", inputFile)

	// Check if output file already exists
	if _, err := os.Stat(outputFile); err == nil {
		log.Warn("Output file already exists: %s", outputFile)
		return nil
	}

	// Check if file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		log.Error("File does not exist: %s", inputFile)
		return fmt.Errorf("file does not exist: %s", inputFile)
	}

	// Check file size
	fileInfo, err := os.Stat(inputFile)
	if err != nil {
		log.Error("Failed to get file info for %s: %v", inputFile, err)
		return fmt.Errorf("failed to get file info: %w", err)
	}

	if fileInfo.Size() > maxFileSize {
		log.Error("File size (%d bytes) exceeds maximum allowed size (%d bytes): %s", 
			fileInfo.Size(), maxFileSize, inputFile)
		return fmt.Errorf("file size exceeds 200KB limit")
	}

	// Check if file is binary
	isBinary, err := isBinaryFile(inputFile)
	if err != nil {
		log.Error("Failed to check if file is binary %s: %v", inputFile, err)
		return fmt.Errorf("failed to check file type: %w", err)
	}

	if isBinary {
		log.Error("File appears to be binary, only text files are supported: %s", inputFile)
		return fmt.Errorf("binary files are not supported")
	}

	// Validate file extension
	ext := filepath.Ext(inputFile)
	if ext != ".txt" && ext != ".md" {
		log.Warn("File extension '%s' is not .txt or .md, proceeding anyway", ext)
	}

	log.Info("File validation successful")
	log.Info("File: %s, Size: %d bytes, Type: text", inputFile, fileInfo.Size())
	log.Info("Output will be written to: %s", outputFile)

	// Initialize OpenAI provider (it handles its own config)
	aiProvider, err := ai.NewOpenAIAIProvider(log)
	if err != nil {
		log.Error("Failed to initialize AI provider: %v", err)
		return fmt.Errorf("failed to initialize AI provider: %w", err)
	}

	// Read file content
	content, err := os.ReadFile(inputFile)
	if err != nil {
		log.Error("Failed to read file %s: %v", inputFile, err)
		return fmt.Errorf("failed to read file: %w", err)
	}

	log.Info("Processing document with AI provider...")

	// Extract information using AI
	result, err := aiProvider.GetReceiptInvoiceInfo(string(content))
	if err != nil {
		log.Error("Failed to extract information: %v", err)
		return fmt.Errorf("failed to extract information: %w", err)
	}

	log.Info("Successfully extracted information from document")

	// Generate suggested filename and populate the field
	result.SuggestedFileName = generateSuggestedFileName(result)
	log.Info("Generated suggested filename: %s", result.SuggestedFileName)

	// Convert result to JSON
	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Error("Failed to marshal result to JSON: %v", err)
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	// Output the JSON result to console
	fmt.Println(string(jsonOutput))

	// Write JSON to output file
	err = os.WriteFile(outputFile, jsonOutput, 0644)
	if err != nil {
		log.Error("Failed to write output file %s: %v", outputFile, err)
		return fmt.Errorf("failed to write output file: %w", err)
	}

	log.Info("Successfully wrote JSON output to %s", outputFile)

	return nil
}

// generateSuggestedFileName creates a suggested filename from extracted data
// Format: <date>-<company>-<description>-<amount>SEK (lowercase, non-alphanumeric chars become _)
func generateSuggestedFileName(info *interfaces.ReceiptInvoiceInfo) string {
	// Helper function to clean strings: lowercase and replace non-alphanumeric with _
	cleanString := func(s string) string {
		// Convert to lowercase
		s = strings.ToLower(s)
		// Replace non-alphanumeric characters with underscore
		reg := regexp.MustCompile(`[^a-z0-9]`)
		return reg.ReplaceAllString(s, "_")
	}
	
	// Extract date (use "unknown" if not available)
	date := "unknown"
	if info.DateIssued != nil && *info.DateIssued != "" && *info.DateIssued != "." {
		date = cleanString(*info.DateIssued)
	}
	
	// Extract company (use "unknown" if not available)
	company := "unknown"
	if info.Company != nil && *info.Company != "" && *info.Company != "." {
		company = cleanString(*info.Company)
	}
	
	// Extract description (always available as it's mandatory)
	description := cleanString(info.Description)
	
	// Extract amount in SEK (convert from se_cent_amount to SEK, rounded to nearest krona)
	amountSEK := "unknown"
	if info.SECentAmount != nil && *info.SECentAmount > 0 {
		// Convert öre to SEK and round to nearest krona
		sekAmount := float64(*info.SECentAmount) / 100.0
		roundedSEK := int(math.Round(sekAmount))
		amountSEK = fmt.Sprintf("%dsek", roundedSEK)
	}
	
	// Combine all parts
	return fmt.Sprintf("%s-%s-%s-%s", date, company, description, amountSEK)
}

// isBinaryFile checks if a file appears to be binary by examining the first 512 bytes
func isBinaryFile(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read first 512 bytes to check for binary content
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && n == 0 {
		return false, err
	}

	// Check for null bytes which typically indicate binary content
	// But skip isolated null bytes that might be encoding issues
	nullCount := 0
	for i := 0; i < n; i++ {
		if buffer[i] == 0 {
			nullCount++
			// Allow a few null bytes (could be encoding issues)
			// but many null bytes indicate binary
			if nullCount > 3 {
				return true, nil
			}
		}
	}

	// Reset file to beginning for additional check
	file.Seek(0, 0)
	
	// Additional check: scan for non-printable characters
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		text := scanner.Text()
		nonPrintableCount := 0
		for _, r := range text {
			// Check for control characters (except tab, newline, carriage return)
			if r < 32 && r != 9 && r != 10 && r != 13 {
				nonPrintableCount++
				// Allow some non-printable chars (could be special chars)
				if nonPrintableCount > 5 {
					return true, nil
				}
			}
		}
	}

	return false, nil
}