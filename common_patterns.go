package emailvalidator

import (
	"regexp"
	"strings"
)

// CommonPatterns provides detection for common email patterns
type CommonPatterns struct{}

// NewCommonPatterns creates a new CommonPatterns instance
func NewCommonPatterns() *CommonPatterns {
	return &CommonPatterns{}
}

// IsDisposable checks if the email is from a known disposable email provider
func (c *CommonPatterns) IsDisposable(email string) bool {
	domain := strings.ToLower(strings.Split(email, "@")[1])
	
	// Common disposable email domains (partial list)
	disposableDomains := []string{
		"tempmail.com", "guerrillamail.com", "mailinator.com",
		"10minutemail.com", "throwawaymail.com", "yopmail.com",
		"fakeinbox.com", "trashmail.com", "getairmail.com",
		"dispostable.com", "maildrop.cc", "tmpmail.org",
	}
	
	for _, disposable := range disposableDomains {
		if strings.Contains(domain, disposable) {
			return true
		}
	}
	
	return false
}

// IsRoleAccount checks if the email is a role-based account
func (c *CommonPatterns) IsRoleAccount(email string) bool {
	localPart := strings.ToLower(strings.Split(email, "@")[0])
	
	roleAccounts := []string{
		"admin", "administrator", "webmaster", "postmaster",
		"hostmaster", "abuse", "noc", "security", "info",
		"sales", "support", "contact", "help", "mail",
		"hello", "noreply", "no-reply", "newsletter",
	}
	
	for _, role := range roleAccounts {
		if localPart == role {
			return true
		}
	}
	
	return false
}

// HasCommonPattern checks for common email patterns
func (c *CommonPatterns) HasCommonPattern(email string) string {
	localPart := strings.ToLower(strings.Split(email, "@")[0])
	
	patterns := map[string]string{
		`^\d+$`:                              "numeric_only",
		`^[a-z]+\.[a-z]+$`:                   "first.last",
		`^[a-z]+\d+$`:                        "name_with_numbers",
		`^[a-z]+\.[a-z]+\d+$`:                "first.last_with_numbers",
		`^[a-z]{1,2}\d+$`:                    "initials_with_numbers",
		`^test`:                              "test_account",
		`^demo`:                              "demo_account",
	}
	
	for pattern, patternType := range patterns {
		matched, _ := regexp.MatchString(pattern, localPart)
		if matched {
			return patternType
		}
	}
	
	return "custom"
}