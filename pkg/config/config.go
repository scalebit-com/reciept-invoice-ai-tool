package config

// Config holds the application configuration
// This is now a generic config that can be extended for future needs
// AI providers handle their own configuration internally
type Config struct {
	// Add any global application configuration here in the future
}

// LoadConfig loads application-wide configuration
// Currently returns an empty config as providers handle their own config
func LoadConfig() (*Config, error) {
	return &Config{}, nil
}