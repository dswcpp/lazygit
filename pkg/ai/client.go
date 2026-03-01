package ai

import (
	"context"
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

// Client is the entry point for all AI features.
type Client struct {
	provider Provider
}

// NewClient creates a Client from the user configuration.
// Returns nil, nil when AI is disabled or no active profile is configured (callers must check for nil).
func NewClient(cfg config.AIConfig) (*Client, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	profile := cfg.GetActiveProfile()
	if profile == nil {
		return nil, nil
	}

	apiKey := resolveEnvVars(profile.APIKey)
	endpoint := resolveEndpointForProfile(*profile)
	model := profile.Model
	maxTokens := profile.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 500
	}
	timeout := profile.Timeout
	if timeout <= 0 {
		if profile.EnableThinking {
			timeout = 300 // thinking/reasoning models need more time
		} else {
			timeout = 60
		}
	}

	if model == "" {
		return nil, fmt.Errorf("ai model must be set in the active profile config")
	}

	var provider Provider
	if strings.ToLower(profile.Provider) == "anthropic" {
		provider = newAnthropicProvider(apiKey, model, profile.Endpoint, maxTokens, timeout, profile.CustomHeaders)
	} else {
		provider = newOpenAIProvider(endpoint, apiKey, model, maxTokens, timeout, profile.EnableThinking, profile.CustomHeaders)
	}
	return &Client{provider: provider}, nil
}

// Complete sends a prompt and returns the full AI result (content + reasoning chain).
func (c *Client) Complete(ctx context.Context, prompt string) (Result, error) {
	return c.provider.Complete(ctx, prompt)
}

// CompleteStream sends a prompt and streams the response, calling onChunk for each fragment.
func (c *Client) CompleteStream(ctx context.Context, prompt string, onChunk func(string)) error {
	return c.provider.CompleteStream(ctx, prompt, onChunk)
}

// resolveEndpointForProfile returns the effective API endpoint for a profile.
func resolveEndpointForProfile(profile config.AIProfileConfig) string {
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

// resolveEnvVars replaces ${VAR} or $VAR patterns with environment variable values.
var envVarPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

func resolveEnvVars(s string) string {
	return envVarPattern.ReplaceAllStringFunc(s, func(match string) string {
		// Extract variable name from ${VAR} or $VAR
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
