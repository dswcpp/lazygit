package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefaultConfigAIProfileUsesOpenAIGPT(t *testing.T) {
	cfg := GetDefaultConfig()

	if !assert.NotNil(t, cfg) {
		return
	}
	assert.False(t, cfg.AI.Enabled)
	assert.Equal(t, "gpt-4o-mini", cfg.AI.ActiveProfile)

	if !assert.Len(t, cfg.AI.Profiles, 1) {
		return
	}

	profile := cfg.AI.Profiles[0]
	assert.Equal(t, "gpt-4o-mini", profile.Name)
	assert.Equal(t, "openai", profile.Provider)
	assert.Equal(t, "gpt-4o-mini", profile.Model)
	assert.False(t, profile.EnableThinking)
	assert.Equal(t, 8000, profile.MaxTokens)
	assert.Equal(t, 60, profile.Timeout)
}
