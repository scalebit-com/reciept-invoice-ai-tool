package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/scalebit-com/reciept-invoice-ai-tool/pkg/interfaces"
	pkglogger "github.com/scalebit-com/reciept-invoice-ai-tool/pkg/logger"
)

var logger interfaces.Logger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "reciept-invoice-ai-tool",
	Short: "A CLI tool for extracting structured information from receipts and invoices",
	Long: `A command-line tool that extracts structured information from receipt and invoice 
data in text or markdown files, outputting results as JSON.

The tool can process text (.txt) and markdown (.md) files containing receipt data
and parse information such as store name, date, items, prices, and totals.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Initialize logger
	logger = pkglogger.NewColorLogger()
	
	err := rootCmd.Execute()
	if err != nil {
		logger.Error("Command execution failed: %v", err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.reciept-invoice-ai-tool.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}