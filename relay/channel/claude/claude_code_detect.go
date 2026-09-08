package claude

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/QuantumNous/new-api/setting/operation_setting"
)

// Official Claude Code client detection plus version gating. Used by the
// "Claude Code only" restriction: custom clients are rejected, while the real
// Claude Code CLI is allowed through.

var (
	// User-Agent: claude-cli/x.x.x (official CLI only, case-insensitive)
	claudeCodeUAPattern        = regexp.MustCompile(`(?i)^claude-cli/\d+\.\d+\.\d+`)
	claudeCodeUAVersionPattern = regexp.MustCompile(`(?i)^claude-cli/(\d+\.\d+\.\d+)`)
)

const (
	// Similarity threshold between the request system prompt and the Claude Code template.
	ccSystemPromptThreshold = 0.5
	// Claude Code billing header marker carried inside the system text.
	ccBillingHeaderPrefix = "x-anthropic-billing-header"
	ccEntrypointMarker    = "cc_entrypoint="
)

// Claude Code system prompt templates (the leading identity block).
var ccSystemPromptTemplates = []string{
	"You are Claude Code, Anthropic's official CLI for Claude.",
	"You are Claude Code, Anthropic's official CLI for Claude, running within the Claude Agent SDK.",
}

// ccExtractVersion returns the version from a claude-cli/X.Y.Z user agent ("" when not Claude Code).
func ccExtractVersion(ua string) string {
	m := claudeCodeUAVersionPattern.FindStringSubmatch(ua)
	if len(m) >= 2 {
		return m[1]
	}
	return ""
}

// ccParseSemver returns [major, minor, patch].
func ccParseSemver(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.Split(v, ".")
	out := [3]int{0, 0, 0}
	for i := 0; i < len(parts) && i < 3; i++ {
		if p, err := strconv.Atoi(strings.TrimSpace(parts[i])); err == nil {
			out[i] = p
		}
	}
	return out
}

// ccCompareVersions returns -1 (a<b), 0 (a==b) or 1 (a>b).
func ccCompareVersions(a, b string) int {
	ap, bp := ccParseSemver(a), ccParseSemver(b)
	for i := 0; i < 3; i++ {
		if ap[i] < bp[i] {
			return -1
		}
		if ap[i] > bp[i] {
			return 1
		}
	}
	return 0
}

// diceCoefficient measures bigram similarity between two strings (0..1).
func diceCoefficient(a, b string) float64 {
	a, b = strings.ToLower(a), strings.ToLower(b)
	if a == b {
		return 1.0
	}
	if len(a) < 2 || len(b) < 2 {
		return 0.0
	}
	bigrams := func(s string) map[string]int {
		m := make(map[string]int)
		r := []rune(s)
		for i := 0; i < len(r)-1; i++ {
			m[string(r[i:i+2])]++
		}
		return m
	}
	ba, bb := bigrams(a), bigrams(b)
	overlap := 0
	for g, ca := range ba {
		if cb, ok := bb[g]; ok {
			if ca < cb {
				overlap += ca
			} else {
				overlap += cb
			}
		}
	}
	total := 0
	for _, c := range ba {
		total += c
	}
	for _, c := range bb {
		total += c
	}
	if total == 0 {
		return 0.0
	}
	return 2.0 * float64(overlap) / float64(total)
}

// ccBestSimilarity returns the best similarity between text and the Claude Code
// templates. The template is the head of the system block, so only a prefix of
// the same length (plus slack) is compared.
func ccBestSimilarity(text string) float64 {
	best := 0.0
	for _, tpl := range ccSystemPromptTemplates {
		head := text
		if len(head) > len(tpl)+64 {
			head = head[:len(tpl)+64]
		}
		if s := diceCoefficient(head, tpl); s > best {
			best = s
		}
		// A literal template match is an immediate maximum.
		if strings.Contains(text, tpl) {
			return 1.0
		}
	}
	return best
}

// ccHasSystemPrompt reports whether any system text carries a Claude Code
// signature: the billing marker cc_entrypoint= or a similar enough prompt.
func ccHasSystemPrompt(systemTexts []string) bool {
	for _, text := range systemTexts {
		if text == "" {
			continue
		}
		if strings.Contains(text, ccBillingHeaderPrefix) && strings.Contains(text, ccEntrypointMarker) {
			return true
		}
		if ccBestSimilarity(text) >= ccSystemPromptThreshold {
			return true
		}
	}
	return false
}

// IsClaudeCodeRequest performs the strict official Claude Code detection:
// (1) a claude-cli/x.x.x user agent AND (2) a system prompt that matches the
// Claude Code template (or carries the billing marker). count_tokens is an
// internal Claude Code call without the full system, so the user agent is enough.
func IsClaudeCodeRequest(userAgent string, systemTexts []string, isCountTokens bool) bool {
	if !claudeCodeUAPattern.MatchString(strings.TrimSpace(userAgent)) {
		return false
	}
	if isCountTokens {
		return true
	}
	return ccHasSystemPrompt(systemTexts)
}

// ClaudeCodeRestrictEnabled reports whether the "Claude Code only" restriction is on.
func ClaudeCodeRestrictEnabled() bool {
	return operation_setting.ClaudeCodeOnly
}

// CheckClaudeCodeAccess is the entry point for the restricted mode. It returns
// (false, message) when the request must be rejected:
//   - restriction disabled: always allowed;
//   - client is not official Claude Code: rejected;
//   - Claude Code with a version outside the configured min/max range: rejected.
func CheckClaudeCodeAccess(userAgent string, systemTexts []string, isCountTokens bool) (bool, string) {
	if !ClaudeCodeRestrictEnabled() {
		return true, ""
	}
	if !IsClaudeCodeRequest(userAgent, systemTexts, isCountTokens) {
		return false, "This service is restricted to official Claude Code."
	}
	minV := strings.TrimSpace(operation_setting.ClaudeCodeMinVersion)
	maxV := strings.TrimSpace(operation_setting.ClaudeCodeMaxVersion)
	if minV == "" && maxV == "" {
		return true, ""
	}
	cv := ccExtractVersion(userAgent)
	if cv == "" {
		return false, "Unable to determine Claude Code version. Please update Claude Code: npm update -g @anthropic-ai/claude-code"
	}
	if minV != "" && ccCompareVersions(cv, minV) < 0 {
		return false, "Your Claude Code version (" + cv + ") is below the minimum required (" + minV + "). Please update: npm update -g @anthropic-ai/claude-code"
	}
	if maxV != "" && ccCompareVersions(cv, maxV) > 0 {
		return false, "Your Claude Code version (" + cv + ") exceeds the maximum allowed (" + maxV + ")."
	}
	return true, ""
}
