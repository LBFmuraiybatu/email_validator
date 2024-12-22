package emailvalidator

import "regexp"

// Common email patterns for additional validation
type EmailPatterns struct {
	CommonProviders map[string]*regexp.Regexp
}

// NewEmailPatterns creates a new EmailPatterns instance
func NewEmailPatterns() *EmailPatterns {
	return &EmailPatterns{
		CommonProviders: map[string]*regexp.Regexp{
			"gmail":     regexp.MustCompile(`^[a-z0-9](\.?[a-z0-9]){5,}@gmail\.com$`),
			"outlook":   regexp.MustCompile(`^[a-z0-9](\.?[a-z0-9_-]){1,}@(outlook|hotmail)\.com$`),
			"yahoo":     regexp.MustCompile(`^[a-z0-9](\.?[a-z0-9_-]){1,}@yahoo\.com$`),
			"icloud":    regexp.MustCompile(`^[a-z0-9](\.?[a-z0-9]){1,}@icloud\.com$`),
			"proton":    regexp.MustCompile(`^[a-z0-9](\.?[a-z0-9_-]){1,}@proton(mail)?\.(com|ch)$`),
		},
	}
}

// MatchesProviderPattern checks if email matches known provider patterns
func (p *EmailPatterns) MatchesProviderPattern(email string) string {
	lowerEmail := strings.ToLower(email)
	
	for provider, pattern := range p.CommonProviders {
		if pattern.MatchString(lowerEmail) {
			return provider
		}
	}
	
	return ""
}