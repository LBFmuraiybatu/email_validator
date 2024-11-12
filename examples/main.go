package main

import (
	"encoding/json"
	"fmt"
	"log"

	"emailvalidator" // Replace with your actual module path
)

func main() {
	// Create a new email validator
	validator := emailvalidator.NewEmailValidator()

	// Test email addresses
	testEmails := []string{
		"user@example.com",
		"invalid-email",
		"user@tempmail.com",
		"verylonglocalpart123456789012345678901234567890123456789012345678901234567890@example.com",
		"",
		"user@nonexistentdomain12345.com",
		"test.email+tag@sub.domain.co.uk",
	}

	fmt.Println("Email Validation Results:")
	fmt.Println("=========================")

	for _, email := range testEmails {
		fmt.Printf("\nTesting: %s\n", email)
		
		result := validator.Validate(email)
		
		// Pretty print the result
		jsonResult, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			log.Printf("Error marshaling result: %v", err)
			continue
		}
		
		fmt.Println(string(jsonResult))
	}

	// Example of adding custom disposable domain
	fmt.Println("\n--- Adding custom disposable domain ---")
	validator.AddDisposableDomain("mycustomdisposable.com")
	
	testEmail := "user@mycustomdisposable.com"
	result := validator.Validate(testEmail)
	fmt.Printf("Testing %s - IsDisposable: %t\n", testEmail, result.IsDisposable)
}

// Example function showing batch validation
func batchValidateEmails(emails []string) {
	validator := emailvalidator.NewEmailValidator()
	
	validEmails := []string{}
	invalidEmails := []string{}
	
	for _, email := range emails {
		result := validator.Validate(email)
		if result.IsValid {
			validEmails = append(validEmails, email)
		} else {
			invalidEmails = append(invalidEmails, email)
		}
	}
	
	fmt.Printf("Valid emails: %d\n", len(validEmails))
	fmt.Printf("Invalid emails: %d\n", len(invalidEmails))
}