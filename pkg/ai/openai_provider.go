package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/invopop/jsonschema"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/scalebit-com/reciept-invoice-ai-tool/pkg/interfaces"
)

// OpenAIAIProvider implements the AIProvider interface using OpenAI's API
type OpenAIAIProvider struct {
	client *openai.Client
	logger interfaces.Logger
	model  string
}

// NewOpenAIAIProvider creates a new OpenAI AI provider
func NewOpenAIAIProvider(logger interfaces.Logger) (*OpenAIAIProvider, error) {
	// Try to load .env file from current directory
	// It's okay if the file doesn't exist
	err := godotenv.Load()
	if err != nil {
		// Only log if it's not a "file not found" error
		if !os.IsNotExist(err) {
			logger.Warn("Could not load .env file: %v", err)
		}
	}

	// Get OpenAI API key from environment
	apiKey := os.Getenv("OPENAI_KEY")
	if apiKey == "" {
		logger.Error("OPENAI_KEY environment variable is not set")
		logger.Error("Please set OPENAI_KEY in your environment or create a .env file with OPENAI_KEY=your-api-key")
		return nil, fmt.Errorf("OPENAI_KEY environment variable is required")
	}

	logger.Debug("Initializing OpenAI provider with API key: sk-...%s", apiKey[len(apiKey)-4:])
	
	// Get model from environment or use default
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = string(openai.ChatModelGPT4o2024_08_06)
		logger.Warn("OPENAI_MODEL environment variable is not set, defaulting to %s", model)
	} else {
		logger.Info("Using OpenAI model: %s", model)
	}
	
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	logger.Info("OpenAI provider initialized successfully")

	return &OpenAIAIProvider{
		client: &client,
		logger: logger,
		model:  model,
	}, nil
}

// GenerateSchema generates a JSON schema for structured outputs
func GenerateSchema[T any]() interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

// Generate the JSON schema at initialization time
var ReceiptInvoiceInfoSchema = GenerateSchema[interfaces.ReceiptInvoiceInfo]()

// GetReceiptInvoiceInfo extracts structured information from receipt/invoice text
func (p *OpenAIAIProvider) GetReceiptInvoiceInfo(content string) (*interfaces.ReceiptInvoiceInfo, error) {
	startTime := time.Now()
	ctx := context.Background()

	p.logger.Info("Starting OpenAI API request for receipt/invoice extraction")
	p.logger.Debug("Document content length: %d characters", len(content))
	
	// Log truncated content for debugging (first 200 chars)
	contentPreview := content
	if len(contentPreview) > 200 {
		contentPreview = contentPreview[:200] + "..."
	}
	p.logger.Debug("Document preview: %s", strings.ReplaceAll(contentPreview, "\n", " "))

	// System prompt for the AI to act as an accountant
	systemPrompt := `You are an experienced accountant reviewing financial documents. Your task is to:
1. Classify the document as either "None" (not a financial document), "Invoice", or "Receipt"
2. Extract the company name that is offering the service and requesting payment
3. Extract the date the document was issued (in YYYY-MM-DD format)
4. Extract a concise description of the service or items paid for
5. Extract the total amount in Swedish currency (SEK) and convert it to Swedish cents (öre)
   - For amounts in SEK: multiply by 100 (e.g., 95.37 SEK = 9537)
   - For amounts in EUR or other currencies: convert to SEK first using approximate rates (1 EUR ≈ 11.5 SEK), then to cents
   - Return null if no amount is found or if conversion is not possible

Be precise and extract only information that is clearly present in the document.`

	userPrompt := fmt.Sprintf("Please analyze the following document and extract the required information:\n\n%s", content)

	p.logger.Debug("System prompt length: %d characters", len(systemPrompt))
	p.logger.Debug("User prompt length: %d characters", len(userPrompt))

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "receipt_invoice_info",
		Description: openai.String("Structured information extracted from a receipt or invoice"),
		Schema:      ReceiptInvoiceInfoSchema,
		Strict:      openai.Bool(true),
	}

	p.logger.Info("Calling OpenAI Chat Completions API")
	p.logger.Info("Model: %s", p.model)
	p.logger.Info("Using structured output with schema: %s (strict mode: true)", schemaParam.Name)

	// Query the Chat Completions API with structured output
	chat, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
		// Use configured model
		Model: openai.ChatModel(p.model),
	})

	duration := time.Since(startTime)

	if err != nil {
		p.logger.Error("OpenAI API call failed after %v: %v", duration, err)
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	p.logger.Info("OpenAI API call successful (took %v)", duration)
	
	// Log response details
	if len(chat.Choices) > 0 {
		choice := chat.Choices[0]
		p.logger.Debug("Finish reason: %s", choice.FinishReason)
		p.logger.Debug("Response content length: %d characters", len(choice.Message.Content))
		
		// Log token usage if available
		if chat.Usage.TotalTokens > 0 {
			p.logger.Info("Token usage - Prompt: %d, Completion: %d, Total: %d",
				chat.Usage.PromptTokens,
				chat.Usage.CompletionTokens,
				chat.Usage.TotalTokens)
		}
	}

	// Parse the JSON response into our struct
	p.logger.Debug("Parsing JSON response from OpenAI")
	var result interfaces.ReceiptInvoiceInfo
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &result)
	if err != nil {
		p.logger.Error("Failed to parse OpenAI JSON response: %v", err)
		p.logger.Debug("Raw response that failed to parse: %s", chat.Choices[0].Message.Content)
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	p.logger.Info("Successfully parsed OpenAI response")
	p.logger.Info("Extracted document type: %s", result.DocumentType)
	
	if result.DateIssued != nil {
		p.logger.Info("Extracted date: %s", *result.DateIssued)
	} else {
		p.logger.Debug("No date found in document")
	}
	
	if result.Company != nil {
		p.logger.Info("Extracted company: %s", *result.Company)
	} else {
		p.logger.Debug("No company found in document")
	}
	
	if result.ServiceDescription != nil {
		p.logger.Info("Extracted service: %s", *result.ServiceDescription)
	} else {
		p.logger.Debug("No service description found in document")
	}
	
	if result.SECentAmount != nil {
		p.logger.Info("Extracted amount: %d öre (%.2f SEK)", *result.SECentAmount, float64(*result.SECentAmount)/100)
	} else {
		p.logger.Debug("No amount found in document")
	}

	p.logger.Info("Total processing time: %v", time.Since(startTime))

	return &result, nil
}