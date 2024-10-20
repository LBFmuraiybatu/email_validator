package emailvalidator

// Option defines functional options for EmailValidator
type Option func(*EmailValidator)

// WithAllowedTLDs sets allowed top-level domains
func WithAllowedTLDs(tlds []string) Option {
	return func(ev *EmailValidator) {
		ev.allowTLDs = tlds
	}
}

// WithBlockedDomains sets blocked domains
func WithBlockedDomains(domains []string) Option {
	return func(ev *EmailValidator) {
		ev.blockedDomains = make(map[string]bool)
		for _, domain := range domains {
			ev.blockedDomains[strings.ToLower(domain)] = true
		}
	}
}

// WithIPAddresses allows IP address domains
func WithIPAddresses(allow bool) Option {
	return func(ev *EmailValidator) {
		ev.allowIPAddresses = allow
	}
}