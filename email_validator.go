package emailvalidator

import (
	"errors"
	"net"
	"regexp"
	"strings"
	"unicode"
)

// EmailValidator provides methods to validate email addresses
type EmailValidator struct {
	strictMode bool
}

// New creates a new EmailValidator instance
func New() *EmailValidator {
	return &EmailValidator{
		strictMode: false,
	}
}

// NewStrict creates a new EmailValidator with strict validation
func NewStrict() *EmailValidator {
	return &EmailValidator{
		strictMode: true,
	}
}

// ValidationResult contains detailed validation results
type ValidationResult struct {
	IsValid      bool     `json:"is_valid"`
	Errors       []string `json:"errors,omitempty"`
	Warnings     []string `json:"warnings,omitempty"`
	Normalized   string   `json:"normalized,omitempty"`
	Domain       string   `json:"domain,omitempty"`
	Username     string   `json:"username,omitempty"`
}

// Validate performs comprehensive email validation
func (v *EmailValidator) Validate(email string) ValidationResult {
	result := ValidationResult{}
	
	// Basic format check
	if !v.isValidFormat(email) {
		result.Errors = append(result.Errors, "Invalid email format")
		return result
	}
	
	// Extract parts
	username, domain := v.splitEmail(email)
	result.Username = username
	result.Domain = domain
	
	// Validate username
	if err := v.validateUsername(username); err != nil {
		result.Errors = append(result.Errors, err.Error())
	}
	
	// Validate domain
	if err := v.validateDomain(domain); err != nil {
		result.Errors = append(result.Errors, err.Error())
	}
	
	// Check for common typos
	if warning := v.checkForTypos(email); warning != "" {
		result.Warnings = append(result.Warnings, warning)
	}
	
	// Normalize email (lowercase)
	result.Normalized = strings.ToLower(strings.TrimSpace(email))
	
	result.IsValid = len(result.Errors) == 0
	return result
}

// isValidFormat checks basic email format using regex
func (v *EmailValidator) isValidFormat(email string) bool {
	// RFC 5322 compliant regex (simplified version)
	pattern := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// splitEmail splits email into username and domain parts
func (v *EmailValidator) splitEmail(email string) (string, string) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// validateUsername checks username part constraints
func (v *EmailValidator) validateUsername(username string) error {
	if len(username) == 0 {
		return errors.New("username cannot be empty")
	}
	
	if len(username) > 64 {
		return errors.New("username too long (max 64 characters)")
	}
	
	// Check for consecutive dots
	if strings.Contains(username, "..") {
		return errors.New("username cannot contain consecutive dots")
	}
	
	// Check if starts or ends with dot
	if strings.HasPrefix(username, ".") || strings.HasSuffix(username, ".") {
		return errors.New("username cannot start or end with a dot")
	}
	
	// In strict mode, check for special characters
	if v.strictMode {
		for _, char := range username {
			if !v.isValidUsernameChar(char) {
				return errors.New("username contains invalid characters")
			}
		}
	}
	
	return nil
}

// validateDomain checks domain part constraints
func (v *EmailValidator) validateDomain(domain string) error {
	if len(domain) == 0 {
		return errors.New("domain cannot be empty")
	}
	
	if len(domain) > 253 {
		return errors.New("domain too long (max 253 characters)")
	}
	
	// Check for valid domain structure
	domainParts := strings.Split(domain, ".")
	if len(domainParts) < 2 {
		return errors.New("domain must have at least two parts")
	}
	
	// Check each domain part
	for _, part := range domainParts {
		if len(part) == 0 {
			return errors.New("domain part cannot be empty")
		}
		if len(part) > 63 {
			return errors.New("domain part too long (max 63 characters)")
		}
		if strings.HasPrefix(part, "-") || strings.HasSuffix(part, "-") {
			return errors.New("domain part cannot start or end with hyphen")
		}
		
		// Check for valid characters in domain part
		for _, char := range part {
			if !v.isValidDomainChar(char) {
				return errors.New("domain contains invalid characters")
			}
		}
	}
	
	return nil
}

// isValidUsernameChar checks if character is valid in email username
func (v *EmailValidator) isValidUsernameChar(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char) ||
		char == '!' || char == '#' || char == '$' || char == '%' || char == '&' ||
		char == '\'' || char == '*' || char == '+' || char == '-' || char == '/' ||
		char == '=' || char == '?' || char == '^' || char == '_' || char == '`' ||
		char == '{' || char == '|' || char == '}' || char == '~' || char == '.'
}

// isValidDomainChar checks if character is valid in domain part
func (v *EmailValidator) isValidDomainChar(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char) || char == '-'
}

// checkForTypos looks for common email typos
func (v *EmailValidator) checkForTypos(email string) string {
	lowerEmail := strings.ToLower(email)
	
	// Check for common domain typos
	commonTypos := map[string]string{
		"gmial.com":  "gmail.com",
		"gmal.com":   "gmail.com",
		"gmai.com":   "gmail.com",
		"yahooo.com": "yahoo.com",
		"yaho.com":   "yahoo.com",
		"hotmal.com": "hotmail.com",
		"hotmai.com": "hotmail.com",
	}
	
	for typo, correct := range commonTypos {
		if strings.Contains(lowerEmail, "@"+typo) {
			return "Possible typo detected: " + typo + " should be " + correct
		}
	}
	
	return ""
}

// IsDisposableDomain checks if the email domain is from a known disposable email service
func (v *EmailValidator) IsDisposableDomain(email string) bool {
	_, domain := v.splitEmail(email)
	
	// List of common disposable email domains (truncated for example)
	disposableDomains := map[string]bool{
		"tempmail.com":    true,
		"throwaway.com":   true,
		"guerrillamail.com": true,
		"mailinator.com":  true,
		"yopmail.com":     true,
		"10minutemail.com": true,
	}
	
	return disposableDomains[strings.ToLower(domain)]
}

// HasMXRecord checks if the domain has valid MX records
func (v *EmailValidator) HasMXRecord(email string) bool {
	_, domain := v.splitEmail(email)
	
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return false
	}
	
	return true
}