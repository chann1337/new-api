package claude

import (
	"testing"

	"github.com/QuantumNous/new-api/setting/operation_setting"

	"github.com/stretchr/testify/assert"
)

const claudeCodeSystemPrompt = "You are Claude Code, Anthropic's official CLI for Claude.\n\nYou are an interactive CLI tool that helps users with software engineering tasks."

func TestCheckClaudeCodeAccess(t *testing.T) {
	cases := []struct {
		name          string
		restrict      bool
		minVersion    string
		maxVersion    string
		userAgent     string
		systemTexts   []string
		isCountTokens bool
		allowed       bool
	}{
		{
			name:      "restriction disabled allows any client",
			userAgent: "curl/8.5.0",
			allowed:   true,
		},
		{
			name:        "official claude code passes",
			restrict:    true,
			userAgent:   "claude-cli/2.1.207 (external, cli)",
			systemTexts: []string{claudeCodeSystemPrompt},
			allowed:     true,
		},
		{
			name:        "claude code user agent without matching system is rejected",
			restrict:    true,
			userAgent:   "claude-cli/2.1.207 (external, cli)",
			systemTexts: []string{"You are a helpful assistant."},
			allowed:     false,
		},
		{
			name:        "billing header marker passes without the prompt template",
			restrict:    true,
			userAgent:   "claude-cli/2.1.207 (external, cli)",
			systemTexts: []string{"x-anthropic-billing-header: cc_entrypoint=cli"},
			allowed:     true,
		},
		{
			name:          "count_tokens passes on user agent alone",
			restrict:      true,
			userAgent:     "claude-cli/2.1.207 (external, cli)",
			isCountTokens: true,
			allowed:       true,
		},
		{
			name:        "other clients are rejected",
			restrict:    true,
			userAgent:   "python-httpx/0.27.0",
			systemTexts: []string{claudeCodeSystemPrompt},
			allowed:     false,
		},
		{
			name:        "version below minimum is rejected",
			restrict:    true,
			minVersion:  "2.1.207",
			userAgent:   "claude-cli/2.1.100 (external, cli)",
			systemTexts: []string{claudeCodeSystemPrompt},
			allowed:     false,
		},
		{
			name:        "version equal to minimum passes",
			restrict:    true,
			minVersion:  "2.1.207",
			userAgent:   "claude-cli/2.1.207 (external, cli)",
			systemTexts: []string{claudeCodeSystemPrompt},
			allowed:     true,
		},
		{
			name:        "version above maximum is rejected",
			restrict:    true,
			maxVersion:  "2.1.207",
			userAgent:   "claude-cli/2.2.0 (external, cli)",
			systemTexts: []string{claudeCodeSystemPrompt},
			allowed:     false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			operation_setting.ClaudeCodeOnly = tc.restrict
			operation_setting.ClaudeCodeMinVersion = tc.minVersion
			operation_setting.ClaudeCodeMaxVersion = tc.maxVersion
			t.Cleanup(func() {
				operation_setting.ClaudeCodeOnly = false
				operation_setting.ClaudeCodeMinVersion = ""
				operation_setting.ClaudeCodeMaxVersion = ""
			})

			ok, message := CheckClaudeCodeAccess(tc.userAgent, tc.systemTexts, tc.isCountTokens)
			assert.Equal(t, tc.allowed, ok)
			if tc.allowed {
				assert.Empty(t, message)
			} else {
				assert.NotEmpty(t, message)
			}
		})
	}
}
