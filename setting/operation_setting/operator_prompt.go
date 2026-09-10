package operation_setting

import (
	"strings"

	"github.com/QuantumNous/new-api/common"
)

// ModelSystemPrompts maps a model name to an operator-grade system prompt.
// The special key "*" provides a fallback prompt for every model; an exact
// model name match always wins over "*".
var ModelSystemPrompts = map[string]string{}

// OperatorPromptWildcard is the fallback key in ModelSystemPrompts.
const OperatorPromptWildcard = "*"

func UpdateModelSystemPromptsByJSONString(s string) error {
	prompts := map[string]string{}
	if strings.TrimSpace(s) != "" {
		if err := common.UnmarshalJsonStr(s, &prompts); err != nil {
			return err
		}
	}
	ModelSystemPrompts = prompts
	return nil
}

func ModelSystemPrompts2JSONString() string {
	if ModelSystemPrompts == nil {
		return "{}"
	}
	jsonBytes, err := common.Marshal(ModelSystemPrompts)
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}

// GetModelOperatorPrompt returns the operator system prompt configured for
// the given model. An exact model name match wins; "*" applies to all models.
func GetModelOperatorPrompt(model string) string {
	if model != "" {
		if prompt, ok := ModelSystemPrompts[model]; ok && strings.TrimSpace(prompt) != "" {
			return prompt
		}
	}
	if prompt, ok := ModelSystemPrompts[OperatorPromptWildcard]; ok {
		return prompt
	}
	return ""
}
