package emailvalidator

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// DNSChecker provides DNS validation for email domains
type DNSChecker struct {
	timeout time.Duration
}

// NewDNSChecker creates a new DNSChecker instance
func NewDNSChecker() *DNSChecker {
	return &DNSChecker{
		timeout: 5 * time.Second,
	}
}

// WithTimeout sets the DNS lookup timeout
func (d *DNSChecker) WithTimeout(timeout time.Duration) *DNSChecker {
	d.timeout = timeout
	return d
}

// HasMXRecords checks if the domain has MX records
func (d *DNSChecker) HasMXRecords(domain string) (bool, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return false, fmt.Errorf("DNS lookup failed: %v", err)
	}
	return len(mxRecords) > 0, nil
}

// HasARecords checks if the domain has A records (fallback for domains without MX)
func (d *DNSChecker) HasARecords(domain string) (bool, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return false, fmt.Errorf("DNS lookup failed: %v", err)
	}
	return len(ips) > 0, nil
}

// IsDomainValid checks if the domain exists and can receive emails
func (d *DNSChecker) IsDomainValid(domain string) (bool, error) {
	// First check for MX records
	hasMX, err := d.HasMXRecords(domain)
	if err != nil {
		return false, err
	}

	// If no MX records, check for A records
	if !hasMX {
		hasA, err := d.HasARecords(domain)
		if err != nil {
			return false, err
		}
		return hasA, nil
	}

	return true, nil
}

// ValidateEmailDomain validates the domain part of an email address
func (d *DNSChecker) ValidateEmailDomain(email string) (bool, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid email format")
	}

	domain := parts[1]
	return d.IsDomainValid(domain)
}