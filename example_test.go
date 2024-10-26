package emailvalidator

import (
	"fmt"
	"testing"
)

func ExampleEmailValidator() {
	// Create validator with default settings
	validator := New()
	
	// Validate an email
	result := validator.Validate("user@example.com")
	fmt.Printf("Valid: %t\n", result.IsValid)
	fmt.Printf("Normalized: %s\n", result.Normalized)
	
	// Output:
	// Valid: true
	// Normalized: user@example.com
}

func ExampleEmailValidator_withOptions() {
	// Create validator with custom options
	validator := New(
		WithAllowedTLDs([]string{"com", "org", "net"}),
		WithBlockedDomains([]string{"spam.com", "fake.org"}),
		WithIPAddresses(true),
	)
	
	// Test various emails
	emails := []string{
		"user@example.com",
		"user@spam.com",      // blocked domain
		"user@example.io",    // invalid TLD
		"user@[192.168.1.1]", // IP address domain
	}
	
	for _, email := range emails {
		result := validator.Validate(email)
		fmt.Printf("%s: %t (Errors: %v)\n", email, result.IsValid, result.Errors)
	}
}

func TestEmailValidation(t *testing.T) {
	validator := New()
	
	testCases := []struct {
		email    string
		expected bool
	}{
		{"user@example.com", true},
		{"user.name@example.com", true},
		{"user+tag@example.com", true},
		{"invalid-email", false},
		{"user@", false},
		{"@example.com", false},
		{"user@.com", false},
		{"user@com", false},
		{"", false},
	}
	
	for _, tc := range testCases {
		result := validator.Validate(tc.email)
		if result.IsValid != tc.expected {
			t.Errorf("Email %s: expected %t, got %t", tc.email, tc.expected, result.IsValid)
		}
	}
}