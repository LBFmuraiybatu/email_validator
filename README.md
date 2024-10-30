# Go Email Validator

A comprehensive email validation package for Go that checks format compliance according to RFC 5322 standards.

## Features

- ✅ RFC 5322 compliant email format validation
- ✅ Local part and domain validation
- ✅ Customizable TLD restrictions
- ✅ Domain blocking functionality
- ✅ IP address domain support (optional)
- ✅ Disposable email detection
- ✅ Email normalization
- ✅ Extensible with functional options

## Installation

```bash
go get github.com/yourusername/emailvalidator
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/yourusername/emailvalidator"
)

func main() {
    validator := emailvalidator.New()
    
    result := validator.Validate("user@example.com")
    if result.IsValid {
        fmt.Printf("Valid email: %s\n", result.Normalized)
    } else {
        fmt.Printf("Invalid email. Errors: %v\n", result.Errors)
    }
}
```

### Advanced Configuration

```go
package main

import (
    "fmt"
    "github.com/yourusername/emailvalidator"
)

func main() {
    // Create validator with custom options
    validator := emailvalidator.New(
        emailvalidator.WithAllowedTLDs([]string{"com", "org", "net"}),
        emailvalidator.WithBlockedDomains([]string{"spam.com", "fake.org"}),
        emailvalidator.WithIPAddresses(false),
    )
    
    emails := []string{
        "user@example.com",
        "user@spam.com",
        "user@example.io",
        "user@[192.168.1.1]",
    }
    
    for _, email := range emails {
        result := validator.Validate(email)
        fmt.Printf("%s: %t\n", email, result.IsValid)
        if len(result.Errors) > 0 {
            fmt.Printf("  Errors: %v\n", result.Errors)
        }
    }
}
```

### Utility Functions

```go
validator := emailvalidator.New()

// Extract domain
domain := validator.ExtractDomain("user@example.com") // "example.com"

// Extract username
username := validator.ExtractUsername("user@example.com") // "user"

// Check for disposable email
isDisposable := validator.IsDisposableEmail("user@tempmail.com") // true
```

## Validation Rules

The validator checks:

- **Format**: Basic email structure compliance
- **Length**: Local part ≤ 64 chars, entire email ≤ 254 chars
- **Local Part**: No consecutive dots, no leading/trailing dots
- **Domain**: Valid TLD, proper domain structure
- **Special Characters**: Proper handling of special characters in local part

## Options

- `WithAllowedTLDs([]string)`: Restrict to specific TLDs
- `WithBlockedDomains([]string)`: Block specific domains
- `WithIPAddresses(bool)`: Allow IP address domains

## Testing

Run the included tests:

```bash
go test -v
```

## License

MIT License - feel free to use in your projects.

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for bugs and feature requests.