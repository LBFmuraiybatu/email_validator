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

// WithStrictMode enables strict validation (RFC 5322 compliant)
func (v *EmailValidator) WithStrictMode() *EmailValidator {
	v.strictMode = true
	return v
}

// Validate performs comprehensive email validation
func (v *EmailValidator) Validate(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}

	// Basic format validation
	if err := v.validateFormat(email); err != nil {
		return err
	}

	// Local part validation
	if err := v.validateLocalPart(email); err != nil {
		return err
	}

	// Domain part validation
	if err := v.validateDomainPart(email); err != nil {
		return err
	}

	return nil
}

// validateFormat checks the basic email format
func (v *EmailValidator) validateFormat(email string) error {
	// Check for @ symbol
	if !strings.Contains(email, "@") {
		return errors.New("email must contain @ symbol")
	}

	// Split into local and domain parts
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("email must contain exactly one @ symbol")
	}

	localPart := parts[0]
	domainPart := parts[1]

	// Check parts are not empty
	if localPart == "" {
		return errors.New("local part cannot be empty")
	}
	if domainPart == "" {
		return errors.New("domain part cannot be empty")
	}

	// Check total length (RFC 5321)
	if len(email) > 254 {
		return errors.New("email address cannot exceed 254 characters")
	}

	// Check local part length (RFC 5321)
	if len(localPart) > 64 {
		return errors.New("local part cannot exceed 64 characters")
	}

	return nil
}

// validateLocalPart validates the local part of the email
func (v *EmailValidator) validateLocalPart(email string) error {
	localPart := strings.Split(email, "@")[0]

	// Check if local part starts or ends with dot
	if strings.HasPrefix(localPart, ".") || strings.HasSuffix(localPart, ".") {
		return errors.New("local part cannot start or end with a dot")
	}

	// Check for consecutive dots
	if strings.Contains(localPart, "..") {
		return errors.New("local part cannot contain consecutive dots")
	}

	// Validate characters based on mode
	for i, char := range localPart {
		if v.strictMode {
			if !isValidStrictLocalChar(char, i, localPart) {
				return errors.New("local part contains invalid characters")
			}
		} else {
			if !isValidRelaxedLocalChar(char) {
				return errors.New("local part contains invalid characters")
			}
		}
	}

	return nil
}

// validateDomainPart validates the domain part of the email
func (v *EmailValidator) validateDomainPart(email string) error {
	domainPart := strings.Split(email, "@")[1]

	// Check if domain part starts or ends with hyphen or dot
	if strings.HasPrefix(domainPart, "-") || strings.HasSuffix(domainPart, "-") {
		return errors.New("domain cannot start or end with hyphen")
	}
	if strings.HasPrefix(domainPart, ".") || strings.HasSuffix(domainPart, ".") {
		return errors.New("domain cannot start or end with dot")
	}

	// Split domain into labels
	labels := strings.Split(domainPart, ".")
	if len(labels) < 2 {
		return errors.New("domain must have at least two parts")
	}

	// Validate each label
	for _, label := range labels {
		if len(label) == 0 {
			return errors.New("domain label cannot be empty")
		}
		if len(label) > 63 {
			return errors.New("domain label cannot exceed 63 characters")
		}

		// Check label format
		for i, char := range label {
			if i == 0 || i == len(label)-1 {
				// First and last character cannot be hyphen
				if char == '-' {
					return errors.New("domain label cannot start or end with hyphen")
				}
			}

			if !isValidDomainChar(char) {
				return errors.New("domain contains invalid characters")
			}
		}
	}

	// Check TLD (last label)
	tld := labels[len(labels)-1]
	if len(tld) < 2 {
		return errors.New("TLD must be at least 2 characters")
	}

	// TLD should not contain numbers (with some exceptions, but we'll be strict)
	for _, char := range tld {
		if unicode.IsDigit(char) {
			return errors.New("TLD should not contain numbers")
		}
	}

	return nil
}

// Helper functions for character validation
func isValidStrictLocalChar(char rune, position int, localPart string) bool {
	// RFC 5322 compliant local part characters
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == '!' || char == '#' || char == '$' || char == '%' ||
		char == '&' || char == '\'' || char == '*' || char == '+' ||
		char == '-' || char == '/' || char == '=' || char == '?' ||
		char == '^' || char == '_' || char == '`' || char == '{' ||
		char == '|' || char == '}' || char == '~' || char == '.'
}

func isValidRelaxedLocalChar(char rune) bool {
	// More relaxed local part validation (common practice)
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == '.' || char == '-' || char == '_' || char == '+'
}

func isValidDomainChar(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == '-'
}

// IsValidSyntax is a quick syntax check using regex (less accurate but faster)
func (v *EmailValidator) IsValidSyntax(email string) bool {
	// Simple regex for basic email format
	pattern := `^[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}