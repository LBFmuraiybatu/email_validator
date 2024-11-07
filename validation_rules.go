package emailvalidator

import "regexp"

// ValidationRule defines an interface for email validation rules
type ValidationRule interface {
	Validate(email string) error
	Name() string
}

// FormatRule validates email format
type FormatRule struct {
	pattern *regexp.Regexp
}

func NewFormatRule() *FormatRule {
	pattern := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
	return &FormatRule{
		pattern: regexp.MustCompile(pattern),
	}
}

func (r *FormatRule) Validate(email string) error {
	if !r.pattern.MatchString(email) {
		return ValidationError{Rule: r.Name(), Message: "Invalid email format"}
	}
	return nil
}

func (r *FormatRule) Name() string {
	return "format_rule"
}

// LengthRule validates email length constraints
type LengthRule struct{}

func NewLengthRule() *LengthRule {
	return &LengthRule{}
}

func (r *LengthRule) Validate(email string) error {
	if len(email) > 254 {
		return ValidationError{Rule: r.Name(), Message: "Email too long (max 254 characters)"}
	}
	
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ValidationError{Rule: r.Name(), Message: "Invalid email structure"}
	}
	
	if len(parts[0]) > 64 {
		return ValidationError{Rule: r.Name(), Message: "Local part too long (max 64 characters)"}
	}
	
	if len(parts[1]) > 253 {
		return ValidationError{Rule: r.Name(), Message: "Domain too long (max 253 characters)"}
	}
	
	return nil
}

func (r *LengthRule) Name() string {
	return "length_rule"
}

// DisposableDomainRule checks for disposable email domains
type DisposableDomainRule struct {
	disposableDomains map[string]bool
}

func NewDisposableDomainRule(domains map[string]bool) *DisposableDomainRule {
	return &DisposableDomainRule{
		disposableDomains: domains,
	}
}

func (r *DisposableDomainRule) Validate(email string) error {
	parts := strings.Split(strings.ToLower(email), "@")
	if len(parts) != 2 {
		return ValidationError{Rule: r.Name(), Message: "Invalid email structure"}
	}
	
	if r.disposableDomains[parts[1]] {
		return ValidationError{Rule: r.Name(), Message: "Disposable email addresses are not allowed"}
	}
	
	return nil
}

func (r *DisposableDomainRule) Name() string {
	return "disposable_domain_rule"
}

// ValidationError represents a validation error
type ValidationError struct {
	Rule    string
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}