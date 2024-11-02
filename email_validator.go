package emailvalidator

import (
	"errors"
	"net"
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailValidator provides methods for email validation
type EmailValidator struct {
	// Common email validation regex pattern
	emailRegex *regexp.Regexp
	// List of common disposable email domains
	disposableDomains map[string]bool
}

// NewEmailValidator creates a new EmailValidator instance
func NewEmailValidator() *EmailValidator {
	// RFC 5322 compliant email regex (simplified version)
	pattern := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
	
	ev := &EmailValidator{
		emailRegex: regexp.MustCompile(pattern),
		disposableDomains: map[string]bool{
			"tempmail.com": true,
			"throwaway.com": true,
			"guerrillamail.com": true,
			"mailinator.com": true,
			"yopmail.com": true,
			"10minutemail.com": true,
		},
	}
	
	return ev
}

// ValidationResult holds the result of email validation
type ValidationResult struct {
	IsValid          bool     `json:"is_valid"`
	Errors           []string `json:"errors,omitempty"`
	NormalizedEmail  string   `json:"normalized_email,omitempty"`
	IsDisposable     bool     `json:"is_disposable"`
	HasMXRecord      bool     `json:"has_mx_record,omitempty"`
}

// Validate performs comprehensive email validation
func (ev *EmailValidator) Validate(email string) ValidationResult {
	result := ValidationResult{}
	var errors []string

	// Trim spaces and convert to lowercase
	normalized := strings.TrimSpace(strings.ToLower(email))
	result.NormalizedEmail = normalized

	// Check if email is empty
	if normalized == "" {
		errors = append(errors, "Email cannot be empty")
		result.IsValid = false
		result.Errors = errors
		return result
	}

	// Check email length
	if utf8.RuneCountInString(normalized) > 254 {
		errors = append(errors, "Email address too long (max 254 characters)")
	}

	// Validate format using regex
	if !ev.emailRegex.MatchString(normalized) {
		errors = append(errors, "Invalid email format")
	}

	// Split email into local part and domain
	parts := strings.Split(normalized, "@")
	if len(parts) != 2 {
		errors = append(errors, "Invalid email structure")
		result.IsValid = false
		result.Errors = errors
		return result
	}

	localPart := parts[0]
	domain := parts[1]

	// Validate local part
	if len(localPart) > 64 {
		errors = append(errors, "Local part too long (max 64 characters)")
	}

	if localPart == "" {
		errors = append(errors, "Local part cannot be empty")
	}

	// Validate domain
	if len(domain) > 253 {
		errors = append(errors, "Domain too long (max 253 characters)")
	}

	if domain == "" {
		errors = append(errors, "Domain cannot be empty")
	}

	// Check for disposable email
	if ev.isDisposableDomain(domain) {
		result.IsDisposable = true
		errors = append(errors, "Disposable email addresses are not allowed")
	}

	// Check for MX records (optional - can be slow)
	hasMX, err := ev.checkMXRecords(domain)
	if err == nil {
		result.HasMXRecord = hasMX
		if !hasMX {
			errors = append(errors, "Domain does not have valid MX records")
		}
	}

	result.IsValid = len(errors) == 0
	result.Errors = errors

	return result
}

// isDisposableDomain checks if the domain is in the disposable list
func (ev *EmailValidator) isDisposableDomain(domain string) bool {
	return ev.disposableDomains[domain]
}

// checkMXRecords verifies if the domain has MX records
func (ev *EmailValidator) checkMXRecords(domain string) (bool, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return false, err
	}
	return len(mxRecords) > 0, nil
}

// AddDisposableDomain adds a domain to the disposable domains list
func (ev *EmailValidator) AddDisposableDomain(domain string) {
	ev.disposableDomains[strings.ToLower(domain)] = true
}

// RemoveDisposableDomain removes a domain from the disposable domains list
func (ev *EmailValidator) RemoveDisposableDomain(domain string) {
	delete(ev.disposableDomains, strings.ToLower(domain))
}