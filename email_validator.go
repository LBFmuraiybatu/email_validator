package emailvalidator

import (
	"errors"
	"net"
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailValidator provides email validation functionality
type EmailValidator struct {
	allowIPAddresses bool
	allowTLDs        []string
	blockedDomains   map[string]bool
}

// ValidationResult contains detailed validation results
type ValidationResult struct {
	IsValid    bool     `json:"is_valid"`
	Errors     []string `json:"errors,omitempty"`
	Normalized string   `json:"normalized,omitempty"`
}

// New creates a new EmailValidator instance
func New(options ...Option) *EmailValidator {
	validator := &EmailValidator{
		allowIPAddresses: false,
		blockedDomains:   make(map[string]bool),
	}
	
	for _, option := range options {
		option(validator)
	}
	
	return validator
}

// Validate performs comprehensive email validation
func (v *EmailValidator) Validate(email string) ValidationResult {
	result := ValidationResult{}
	
	if !v.isValidFormat(email) {
		result.Errors = append(result.Errors, "Invalid email format")
		return result
	}
	
	localPart, domain := v.splitEmail(email)
	
	// Validate local part
	if err := v.validateLocalPart(localPart); err != nil {
		result.Errors = append(result.Errors, err.Error())
	}
	
	// Validate domain part
	if err := v.validateDomain(domain); err != nil {
		result.Errors = append(result.Errors, err.Error())
	}
	
	// Check blocked domains
	if v.isDomainBlocked(domain) {
		result.Errors = append(result.Errors, "Domain is blocked")
	}
	
	result.IsValid = len(result.Errors) == 0
	if result.IsValid {
		result.Normalized = v.normalizeEmail(email)
	}
	
	return result
}

// isValidFormat checks basic email format using regex
func (v *EmailValidator) isValidFormat(email string) bool {
	if utf8.RuneCountInString(email) > 254 {
		return false
	}
	
	// RFC 5322 compliant regex (simplified)
	pattern := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// splitEmail splits email into local part and domain
func (v *EmailValidator) splitEmail(email string) (string, string) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// validateLocalPart validates the local part of email
func (v *EmailValidator) validateLocalPart(localPart string) error {
	if localPart == "" {
		return errors.New("local part cannot be empty")
	}
	
	if utf8.RuneCountInString(localPart) > 64 {
		return errors.New("local part exceeds maximum length (64 characters)")
	}
	
	// Check for consecutive dots
	if strings.Contains(localPart, "..") {
		return errors.New("local part cannot contain consecutive dots")
	}
	
	// Check if starts or ends with dot
	if strings.HasPrefix(localPart, ".") || strings.HasSuffix(localPart, ".") {
		return errors.New("local part cannot start or end with a dot")
	}
	
	return nil
}

// validateDomain validates the domain part of email
func (v *EmailValidator) validateDomain(domain string) error {
	if domain == "" {
		return errors.New("domain cannot be empty")
	}
	
	// Check for IP address domain
	if strings.HasPrefix(domain, "[") && strings.HasSuffix(domain, "]") {
		if !v.allowIPAddresses {
			return errors.New("IP address domains are not allowed")
		}
		return v.validateIPDomain(domain[1 : len(domain)-1])
	}
	
	// Check domain parts
	domainParts := strings.Split(domain, ".")
	if len(domainParts) < 2 {
		return errors.New("domain must have at least two parts")
	}
	
	// Validate TLD
	tld := domainParts[len(domainParts)-1]
	if !v.isValidTLD(tld) {
		return errors.New("invalid top-level domain")
	}
	
	// Check domain length restrictions
	for _, part := range domainParts {
		if len(part) > 63 {
			return errors.New("domain part exceeds maximum length (63 characters)")
		}
		if strings.HasPrefix(part, "-") || strings.HasSuffix(part, "-") {
			return errors.New("domain part cannot start or end with hyphen")
		}
	}
	
	return nil
}

// validateIPDomain validates IP address domains
func (v *EmailValidator) validateIPDomain(ip string) error {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return errors.New("invalid IP address in domain")
	}
	return nil
}

// isValidTLD checks if TLD is valid
func (v *EmailValidator) isValidTLD(tld string) bool {
	if len(v.allowTLDs) > 0 {
		for _, allowedTLD := range v.allowTLDs {
			if strings.EqualFold(tld, allowedTLD) {
				return true
			}
		}
		return false
	}
	
	// Basic TLD validation - in production, you might want a more comprehensive list
	validTLDs := map[string]bool{
		"com": true, "org": true, "net": true, "edu": true, "gov": true,
		"io": true, "co": true, "info": true, "biz": true, "me": true,
	}
	
	return validTLDs[strings.ToLower(tld)]
}

// isDomainBlocked checks if domain is in blocked list
func (v *EmailValidator) isDomainBlocked(domain string) bool {
	return v.blockedDomains[strings.ToLower(domain)]
}

// normalizeEmail normalizes email address
func (v *EmailValidator) normalizeEmail(email string) string {
	// Convert to lowercase and trim spaces
	return strings.TrimSpace(strings.ToLower(email))
}