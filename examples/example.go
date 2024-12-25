package main

import (
	"encoding/json"
	"fmt"
	"log"

	"yourmodule/emailvalidator"
)

func main() {
	// Create a validator instance
	validator := emailvalidator.New()
	
	// Test email addresses
	testEmails := []string{
		"user@example.com",
		"invalid-email",
		"user@company.co.uk",
		"user.name+tag@domain.com",
		"user@nonexistentdomain.xyz",
		"user@gmial.com", // Common typo
		"", // Empty string
		"a@b.c", // Too short domain
		"verylongusername1234567890123456789012345678901234567890123456789012345@example.com", // Long username
	}
	
	fmt.Println("Email Validation Results:")
	fmt.Println("=========================")
	
	for _, email := range testEmails {
		fmt.Printf("\nTesting: %s\n", email)
		
		result := validator.Validate(email)
		
		// Pretty print the result
		jsonResult, err := json.MarshalIndent(result, "  ", "  ")
		if err != nil {
			log.Printf("Error marshaling result: %v", err)
			continue
		}
		
		fmt.Printf("  Result: %s\n", jsonResult)
		
		// Additional checks
		if result.IsValid {
			fmt.Printf("  Disposable Domain: %t\n", validator.IsDisposableDomain(email))
			fmt.Printf("  Has MX Records: %t\n", validator.HasMXRecord(email))
		}
	}
	
	// Example with strict validation
	fmt.Println("\n\nStrict Validation Example:")
	fmt.Println("=========================")
	strictValidator := emailvalidator.NewStrict()
	strictResult := strictValidator.Validate("user!name@example.com")
	jsonStrict, _ := json.MarshalIndent(strictResult, "  ", "  ")
	fmt.Printf("Result: %s\n", jsonStrict)
}