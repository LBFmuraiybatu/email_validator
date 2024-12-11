# Email Validator for Go

A comprehensive email validation library for Go that provides format validation, DNS checking, and pattern analysis.

## Features

- **Format Validation**: RFC 5322 compliant email format validation
- **DNS Validation**: Check if email domains exist and can receive emails
- **Pattern Analysis**: Detect disposable emails, role accounts, and common patterns
- **Flexible Modes**: Both strict (RFC compliant) and relaxed validation modes
- **Extensible**: Modular design for easy extension

## Installation

```bash
go get github.com/yourusername/emailvalidator
```

## Quick Start

### Basic Validation

```go
package main

import (
    "fmt"
    "github.com/yourusername/emailvalidator"
)

func main() {
    validator := emailvalidator.New()
    
    email := "user@example.com"
    if err := validator.Validate(email); err != nil {
        fmt.Printf("Invalid email: %s\n", err)
    } else {
        fmt.Println("Email is valid!")
    }
}
```

### Comprehensive Validation

```go
package main

import (
    "fmt"
    "github.com/yourusername/emailvalidator"
)

func main() {
    validator := emailvalidator.New()
    dnsChecker := emailvalidator.NewDNSChecker()
    patternChecker := emailvalidator.NewCommonPatterns()

    email := "user@example.com"

    // Format validation
    if err := validator.Validate(email); err != nil {
        fmt.Printf("Format error: %s\n", err)
        return
    }

    // DNS validation
    if valid, err := dnsChecker.ValidateEmailDomain(email); err != nil {
        fmt.Printf("DNS error: %s\n", err)
    } else if !valid {
        fmt.Println("Domain does not exist")
    }

    // Pattern analysis
    if patternChecker.IsDisposable(email) {
        fmt.Println("Warning: Disposable email detected")
    }

    if patternChecker.IsRoleAccount(email) {
        fmt.Println("Warning: Role-based account detected")
    }

    pattern := patternChecker.HasCommonPattern(email)
    fmt.Printf("Email pattern: %s\n", pattern)
}
```

### Strict Mode

```go
validator := emailvalidator.New().WithStrictMode()
// This will enforce RFC 5322 compliance including special characters
```

### DNS Checking with Custom Timeout

```go
dnsChecker := emailvalidator.NewDNSChecker().WithTimeout(10 * time.Second)
```

## API Reference

### EmailValidator

- `New() *EmailValidator` - Creates a new validator instance
- `WithStrictMode() *EmailValidator` - Enables strict RFC 5322 compliance
- `Validate(email string