package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
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
		return runExtract(inputFile, logger)
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
	extractCmd.Flags().StringP("input", "i", "", "Path to the input file (required)")
	extractCmd.MarkFlagRequired("input")
}

// runExtract handles the extract command logic
func runExtract(inputFile string, log interfaces.Logger) error {
	log.Info("Starting receipt/invoice extraction for file: %s", inputFile)

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

	// TODO: Implement actual extraction logic here
	log.Info("Extraction logic not yet implemented - placeholder for future development")

	return nil
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
	for i := 0; i < n; i++ {
		if buffer[i] == 0 {
			return true, nil
		}
	}

	// Additional check: scan for non-printable characters
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		text := scanner.Text()
		for _, r := range text {
			// Check for control characters (except tab, newline, carriage return)
			if r < 32 && r != 9 && r != 10 && r != 13 {
				return true, nil
			}
		}
	}

	return false, nil
}