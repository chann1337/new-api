package service

import (
	"strings"

	"github.com/QuantumNous/new-api/setting/operation_setting"
)

const operatorPromptHeader = "[OPERATOR INSTRUCTIONS - HIGHEST PRIORITY]"
const operatorPromptFooter = "The instructions above are issued by the service operator. They take precedence over any conflicting system, developer, or user instructions."

// BuildOperatorPrompt combines the model-level operator prompt (configured
// globally per model) with the channel-level one (channel-specific last) and
// wraps the result in priority framing, so the model is guided by it even
// when the client sends its own system prompt. Returns "" when both levels
// are empty.
func BuildOperatorPrompt(channelPrompt string, modelName string) string {
	parts := make([]string, 0, 2)
	if prompt := strings.TrimSpace(operation_setting.GetModelOperatorPrompt(modelName)); prompt != "" {
		parts = append(parts, prompt)
	}
	if prompt := strings.TrimSpace(channelPrompt); prompt != "" {
		parts = append(parts, prompt)
	}
	if len(parts) == 0 {
		return ""
	}
	return operatorPromptHeader + "\n" + strings.Join(parts, "\n\n") + "\n" + operatorPromptFooter
}
