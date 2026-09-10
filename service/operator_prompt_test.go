package service

import (
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildOperatorPrompt(t *testing.T) {
	tests := []struct {
		name          string
		modelPrompts  map[string]string
		channelPrompt string
		modelName     string
		wantEmpty     bool
		wantContains  []string
		wantOrder     []string
	}{
		{
			name:         "no prompts configured",
			modelPrompts: map[string]string{},
			modelName:    "gpt-4o",
			wantEmpty:    true,
		},
		{
			name:         "model exact match",
			modelPrompts: map[string]string{"gpt-4o": "be concise"},
			modelName:    "gpt-4o",
			wantContains: []string{"be concise", "HIGHEST PRIORITY"},
		},
		{
			name:         "wildcard fallback applies to any model",
			modelPrompts: map[string]string{"*": "stay polite"},
			modelName:    "claude-3-5-sonnet",
			wantContains: []string{"stay polite"},
		},
		{
			name:          "channel prompt applies without model prompt",
			modelPrompts:  map[string]string{},
			channelPrompt: "channel rules",
			modelName:     "gpt-4o",
			wantContains:  []string{"channel rules"},
		},
		{
			name:          "model prompt precedes channel prompt",
			modelPrompts:  map[string]string{"gpt-4o": "model persona"},
			channelPrompt: "channel rules",
			modelName:     "gpt-4o",
			wantOrder:     []string{"model persona", "channel rules"},
		},
		{
			name: "exact model match wins over wildcard",
			modelPrompts: map[string]string{
				"*":     "generic rules",
				"gpt-4o": "exact rules",
			},
			modelName:    "gpt-4o",
			wantContains: []string{"exact rules"},
		},
		{
			name:         "blank exact match falls back to wildcard",
			modelPrompts: map[string]string{"*": "generic rules", "gpt-4o": "   "},
			modelName:    "gpt-4o",
			wantContains: []string{"generic rules"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation_setting.ModelSystemPrompts = tt.modelPrompts
			t.Cleanup(func() { operation_setting.ModelSystemPrompts = map[string]string{} })

			prompt := BuildOperatorPrompt(tt.channelPrompt, tt.modelName)

			if tt.wantEmpty {
				require.Empty(t, prompt)
				return
			}
			require.NotEmpty(t, prompt)
			for _, want := range tt.wantContains {
				assert.Contains(t, prompt, want)
			}
			if len(tt.wantOrder) > 0 {
				first := strings.Index(prompt, tt.wantOrder[0])
				second := strings.Index(prompt, tt.wantOrder[1])
				require.NotEqual(t, -1, first)
				require.NotEqual(t, -1, second)
				assert.Less(t, first, second)
			}
		})
	}
}
