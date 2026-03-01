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
// Returns nil, nil when AI is disabled (callers must check for nil).
func NewClient(cfg config.AIConfig) (*Client, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	apiKey := resolveEnvVars(cfg.APIKey)
	endpoint := resolveEndpoint(cfg)
	model := cfg.Model
	maxTokens := cfg.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 500
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	if model == "" {
		return nil, fmt.Errorf("ai.model must be set in config")
	}

	provider := newOpenAIProvider(endpoint, apiKey, model, maxTokens, timeout, cfg.EnableThinking)
	return &Client{provider: provider}, nil
}

// Complete sends a prompt and returns the full AI result (content + reasoning chain).
func (c *Client) Complete(ctx context.Context, prompt string) (Result, error) {
	return c.provider.Complete(ctx, prompt)
}

// resolveEndpoint returns the effective API endpoint based on provider setting.
func resolveEndpoint(cfg config.AIConfig) string {
	if cfg.Endpoint != "" {
		return cfg.Endpoint
	}
	switch strings.ToLower(cfg.Provider) {
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
