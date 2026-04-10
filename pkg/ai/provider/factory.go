package provider

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/dswcpp/lazygit/pkg/config"
)

const (
	deepSeekEndpoint = "https://api.deepseek.com/v1"
	openAIEndpoint   = "https://api.openai.com/v1"
	ollamaEndpoint   = "http://localhost:11434/v1"
)

// NewFromConfig creates a Provider from the user's active AI profile configuration.
// Returns nil, nil when AI is disabled or no active profile is found.
func NewFromConfig(cfg config.AIConfig) (Provider, error) {
	if !cfg.Enabled {
		return nil, nil
	}
	profile := cfg.GetActiveProfile()
	if profile == nil {
		return nil, nil
	}

	apiKey := resolveEnvVars(profile.APIKey)
	endpoint := resolveEndpoint(*profile)
	model := profile.Model
	if model == "" {
		return nil, fmt.Errorf("ai model must be set in the active profile config")
	}

	maxTokens := profile.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 500
	}
	timeoutSecs := profile.Timeout
	if timeoutSecs <= 0 {
		if profile.EnableThinking {
			timeoutSecs = 300
		} else {
			timeoutSecs = 60
		}
	}

	if strings.ToLower(profile.Provider) == "anthropic" {
		return NewAnthropicProvider(apiKey, model, profile.Endpoint, maxTokens, timeoutSecs, profile.CustomHeaders), nil
	}
	return NewOpenAIProvider(endpoint, apiKey, model, maxTokens, timeoutSecs, profile.EnableThinking, profile.CustomHeaders, profile.WireAPI), nil
}

func resolveEndpoint(profile config.AIProfileConfig) string {
	if profile.Endpoint != "" {
		return profile.Endpoint
	}
	switch strings.ToLower(profile.Provider) {
	case "deepseek":
		return deepSeekEndpoint
	case "ollama":
		return ollamaEndpoint
	case "openai":
		return openAIEndpoint
	default:
		return deepSeekEndpoint
	}
}

var envVarPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

func resolveEnvVars(s string) string {
	return envVarPattern.ReplaceAllStringFunc(s, func(match string) string {
		sub := envVarPattern.FindStringSubmatch(match)
		name := sub[1]
		if name == "" {
			name = sub[2]
		}
		if val := os.Getenv(name); val != "" {
			return val
		}
		return match
	})
}
