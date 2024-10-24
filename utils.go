package emailvalidator

import (
	"strings"
)

// IsDisposableEmail checks if email is from common disposable email providers
func (v *EmailValidator) IsDisposableEmail(email string) bool {
	_, domain := v.splitEmail(email)
	
	disposableDomains := map[string]bool{
		"tempmail.com":     true,
		"guerrillamail.com": true,
		"mailinator.com":   true,
		"10minutemail.com": true,
		"yopmail.com":      true,
		"throwawaymail.com": true,
	}
	
	return disposableDomains[strings.ToLower(domain)]
}

// ExtractDomain extracts domain from email address
func (v *EmailValidator) ExtractDomain(email string) string {
	_, domain := v.splitEmail(email)
	return domain
}

// ExtractUsername extracts username/local part from email address
func (v *EmailValidator) ExtractUsername(email string) string {
	localPart, _ := v.splitEmail(email)
	return localPart
}