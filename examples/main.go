package main

import (
	"fmt"
	"log"

	"github.com/yourusername/emailvalidator" // Adjust import path as needed
)

func main() {
	// Create validator instances
	validator := emailvalidator.New()
	dnsChecker := emailvalidator.NewDNSChecker()
	patternChecker := emailvalidator.NewCommonPatterns()

	// Test email addresses
	testEmails := []string{
		"user@example.com",
		"invalid-email",
		"user@nonexistentdomain.xyz",
		"test@mailinator.com",
		"admin@company.com",
		"user.name+tag@sub.domain.co.uk",
		"", // empty
		"@domain.com", // missing local part
		"user@", // missing domain
		"a@b.c", // too short TLD
		"user@192.168.1.1", // IP address (not supported in basic validation)
	}

	fmt.Println("Email Validation Results:")
	fmt.Println("========================")

	for _, email := range testEmails {
		fmt.Printf("\nTesting: %s\n", email)
		
		// Basic validation
		if err := validator.Validate(email); err != nil {
			fmt.Printf("  ❌ Validation failed: %s\n", err)
			continue
		}
		fmt.Printf("  ✅ Format is valid\n")
		
		// DNS check
		if valid, err := dnsChecker.ValidateEmailDomain(email); err != nil {
			fmt.Printf("  ⚠️  DNS check error: %s\n", err)
		} else if valid {
			fmt.Printf("  ✅ Domain exists\n")
		} else {
			fmt.Printf("  ❌ Domain does not exist\n")
		}
		
		// Pattern analysis
		if patternChecker.IsDisposable(email) {
			fmt.Printf("  ⚠️  Disposable email detected\n")
		}
		
		if patternChecker.IsRoleAccount(email) {
			fmt.Printf("  ⚠️  Role-based account detected\n")
		}
		
		pattern := patternChecker.HasCommonPattern(email)
		fmt.Printf("  📝 Pattern type: %s\n", pattern)
	}

	// Example with strict mode
	fmt.Println("\n\nStrict Mode Example:")
	fmt.Println("===================")
	strictValidator := emailvalidator.New().WithStrictMode()
	
	strictTestEmails := []string{
		"user!name@example.com", // Valid in strict mode
		"user#name@example.com", // Valid in strict mode
		"user name@example.com", // Invalid in both modes
	}
	
	for _, email := range strictTestEmails {
		if err := strictValidator.Validate(email); err != nil {
			fmt.Printf("❌ %s: %s\n", email, err)
		} else {
			fmt.Printf("✅ %s: Valid\n", email)
		}
	}
}

// Quick validation function for simple use cases
func quickValidate(email string) bool {
	validator := emailvalidator.New()
	return validator.IsValidSyntax(email)
}