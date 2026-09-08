package operation_setting

import "strings"

var DemoSiteEnabled = false
var SelfUseModeEnabled = false

// Claude Code restriction: only requests coming from the official Claude Code
// CLI are relayed, optionally limited to a semver range. Configured at runtime
// through the option map.
var ClaudeCodeOnly = false    // relay only official Claude Code requests
var ClaudeCodeMinVersion = "" // minimum Claude Code version, "" means no lower bound
var ClaudeCodeMaxVersion = "" // maximum Claude Code version, "" means no upper bound

var AutomaticDisableKeywords = []string{
	"Your credit balance is too low",
	"This organization has been disabled.",
	"You exceeded your current quota",
	"Permission denied",
	"The security token included in the request is invalid",
	"Operation not allowed",
	"Your account is not authorized",
}

func AutomaticDisableKeywordsToString() string {
	return strings.Join(AutomaticDisableKeywords, "\n")
}

func AutomaticDisableKeywordsFromString(s string) {
	AutomaticDisableKeywords = []string{}
	ak := strings.Split(s, "\n")
	for _, k := range ak {
		k = strings.TrimSpace(k)
		k = strings.ToLower(k)
		if k != "" {
			AutomaticDisableKeywords = append(AutomaticDisableKeywords, k)
		}
	}
}
