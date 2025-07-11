package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator"
	models "github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-nav-models"
)

func main() {
	// Example 1: Basic data validation
	fmt.Println("=== Basic Data Validation ===")
	basicValidationExample()

	// Example 2: Complex data cleaning
	fmt.Println("\n=== Complex Data Cleaning ===")
	dataCleaningExample()

	// Example 3: Format validation and normalization
	fmt.Println("\n=== Format Validation and Normalization ===")
	formatValidationExample()

	// Example 4: Batch validation with error reporting
	fmt.Println("\n=== Batch Validation with Error Reporting ===")
	batchValidationExample()
}

func basicValidationExample() {
	// Sample user registration data with validation issues
	data := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{
				"id":       1,
				"username": "alice123",
				"email":    "alice@example.com",
				"age":      25,
				"phone":    "123-456-7890",
			},
			map[string]interface{}{
				"id":       2,
				"username": "", // Invalid: empty username
				"email":    "invalid-email", // Invalid: bad email format
				"age":      -5, // Invalid: negative age
				"phone":    "123", // Invalid: too short
			},
			map[string]interface{}{
				"id":       3,
				"username": "charlie",
				"email":    "charlie@test.com",
				"age":      30,
				"phone":    "555-123-4567",
			},
		},
	}

	// Email validation modifier
	emailValidator := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		if email, ok := o.(string); ok {
			emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
			if !emailRegex.MatchString(email) {
				return "INVALID_EMAIL"
			}
		}
		return o
	})

	navigator := &mapnavigator.MapNavigator{
		NodeModifier: emailValidator,
		ReadOnly:     false,
	}

	// Validate all emails
	if result, err := navigator.VisitNode(data, "users", "*", "email"); err == nil {
		fmt.Printf("Email validation results: %v\n", result)
	}

	// Age validation modifier
	ageValidator := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		if age, ok := o.(float64); ok {
			if age < 0 || age > 150 {
				return "INVALID_AGE"
			}
		}
		return o
	})

	ageNav := &mapnavigator.MapNavigator{
		NodeModifier: ageValidator,
		ReadOnly:     false,
	}

	// Validate all ages
	if result, err := ageNav.VisitNode(data, "users", "*", "age"); err == nil {
		fmt.Printf("Age validation results: %v\n", result)
	}

	// Print the validated data
	if jsonData, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Printf("Validated data:\n%s\n", jsonData)
	}
}

func dataCleaningExample() {
	// Sample product data with inconsistent formatting
	data := map[string]interface{}{
		"products": []interface{}{
			map[string]interface{}{
				"name":        "  Gaming Laptop  ", // Extra whitespace
				"description": "High-performance laptop for gaming",
				"price":       "1299.99", // String instead of number
				"category":    "ELECTRONICS", // Inconsistent case
				"tags":        []interface{}{"gaming", "LAPTOP", "High-Performance"},
			},
			map[string]interface{}{
				"name":        "smartphone",
				"description": "Latest smartphone with great features",
				"price":       599.99,
				"category":    "electronics",
				"tags":        []interface{}{"mobile", "PHONE", "technology"},
			},
			map[string]interface{}{
				"name":        "Programming Book",
				"description": "Learn programming with this comprehensive guide",
				"price":       "49.99",
				"category":    "Books",
				"tags":        []interface{}{"programming", "education", "BOOK"},
			},
		},
	}

	// String cleaning modifier
	stringCleaner := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		if str, ok := o.(string); ok {
			// Trim whitespace and normalize case
			cleaned := strings.TrimSpace(str)
			if cleaned != "" {
				// Title case for names, lowercase for categories
				return cleaned
			}
		}
		return o
	})

	// Clean product names
	nameNav := &mapnavigator.MapNavigator{
		NodeModifier: stringCleaner,
		ReadOnly:     false,
	}

	if result, err := nameNav.VisitNode(data, "products", "*", "name"); err == nil {
		fmt.Printf("Cleaned names: %v\n", result)
	}

	// Category normalizer
	categoryNormalizer := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		if category, ok := o.(string); ok {
			return strings.ToLower(strings.TrimSpace(category))
		}
		return o
	})

	categoryNav := &mapnavigator.MapNavigator{
		NodeModifier: categoryNormalizer,
		ReadOnly:     false,
	}

	if result, err := categoryNav.VisitNode(data, "products", "*", "category"); err == nil {
		fmt.Printf("Normalized categories: %v\n", result)
	}

	// Tags normalizer
	tagsNormalizer := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		if tags, ok := o.([]interface{}); ok {
			normalizedTags := make([]interface{}, len(tags))
			for i, tag := range tags {
				if tagStr, ok := tag.(string); ok {
					normalizedTags[i] = strings.ToLower(strings.TrimSpace(tagStr))
				} else {
					normalizedTags[i] = tag
				}
			}
			return normalizedTags
		}
		return o
	})

	tagsNav := &mapnavigator.MapNavigator{
		NodeModifier: tagsNormalizer,
		ReadOnly:     false,
	}

	if result, err := tagsNav.VisitNode(data, "products", "*", "tags"); err == nil {
		fmt.Printf("Normalized tags: %v\n", result)
	}

	// Print the cleaned data
	if jsonData, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Printf("Cleaned product data:\n%s\n", jsonData)
	}
}

func formatValidationExample() {
	// Sample contact data with various formats
	data := map[string]interface{}{
		"contacts": []interface{}{
			map[string]interface{}{
				"name":   "John Doe",
				"phone":  "123-456-7890",
				"email":  "john@example.com",
				"zip":    "12345",
				"ssn":    "123-45-6789",
			},
			map[string]interface{}{
				"name":   "Jane Smith",
				"phone":  "1234567890", // No formatting
				"email":  "jane.smith@company.co.uk",
				"zip":    "12345-6789",
				"ssn":    "987654321", // No formatting
			},
			map[string]interface{}{
				"name":   "Bob Johnson",
				"phone":  "(555) 123-4567", // Different format
				"email":  "bob@test.org",
				"zip":    "98765",
				"ssn":    "invalid-ssn", // Invalid format
			},
		},
	}

	// Phone number formatter
	phoneFormatter := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		if phone, ok := o.(string); ok {
			// Remove all non-digit characters
			digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
			
			// Format as XXX-XXX-XXXX if 10 digits
			if len(digits) == 10 {
				return fmt.Sprintf("%s-%s-%s", digits[:3], digits[3:6], digits[6:])
			}
			// Format as +1-XXX-XXX-XXXX if 11 digits and starts with 1
			if len(digits) == 11 && digits[0] == '1' {
				return fmt.Sprintf("+1-%s-%s-%s", digits[1:4], digits[4:7], digits[7:])
			}
			return "INVALID_PHONE"
		}
		return o
	})

	phoneNav := &mapnavigator.MapNavigator{
		NodeModifier: phoneFormatter,
		ReadOnly:     false,
	}

	if result, err := phoneNav.VisitNode(data, "contacts", "*", "phone"); err == nil {
		fmt.Printf("Formatted phone numbers: %v\n", result)
	}

	// SSN formatter
	ssnFormatter := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		if ssn, ok := o.(string); ok {
			// Remove all non-digit characters
			digits := regexp.MustCompile(`\D`).ReplaceAllString(ssn, "")
			
			// Format as XXX-XX-XXXX if 9 digits
			if len(digits) == 9 {
				return fmt.Sprintf("%s-%s-%s", digits[:3], digits[3:5], digits[5:])
			}
			return "INVALID_SSN"
		}
		return o
	})

	ssnNav := &mapnavigator.MapNavigator{
		NodeModifier: ssnFormatter,
		ReadOnly:     false,
	}

	if result, err := ssnNav.VisitNode(data, "contacts", "*", "ssn"); err == nil {
		fmt.Printf("Formatted SSNs: %v\n", result)
	}

	// ZIP code validator
	zipValidator := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		if zip, ok := o.(string); ok {
			// Check for 5-digit or 9-digit ZIP codes
			if regexp.MustCompile(`^\d{5}$`).MatchString(zip) {
				return zip // Valid 5-digit ZIP
			}
			if regexp.MustCompile(`^\d{5}-\d{4}$`).MatchString(zip) {
				return zip // Valid 9-digit ZIP
			}
			return "INVALID_ZIP"
		}
		return o
	})

	zipNav := &mapnavigator.MapNavigator{
		NodeModifier: zipValidator,
		ReadOnly:     false,
	}

	if result, err := zipNav.VisitNode(data, "contacts", "*", "zip"); err == nil {
		fmt.Printf("Validated ZIP codes: %v\n", result)
	}

	// Print the formatted data
	if jsonData, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Printf("Formatted contact data:\n%s\n", jsonData)
	}
}

func batchValidationExample() {
	// Sample data for batch validation
	data := map[string]interface{}{
		"records": []interface{}{
			map[string]interface{}{
				"id":       1,
				"name":     "Alice Johnson",
				"email":    "alice@example.com",
				"age":      28,
				"salary":   75000.00,
				"status":   "active",
			},
			map[string]interface{}{
				"id":       2,
				"name":     "", // Invalid: empty name
				"email":    "invalid-email", // Invalid: bad email
				"age":      -5, // Invalid: negative age
				"salary":   "not-a-number", // Invalid: non-numeric salary
				"status":   "unknown", // Invalid: unknown status
			},
			map[string]interface{}{
				"id":       3,
				"name":     "Charlie Brown",
				"email":    "charlie@test.com",
				"age":      45,
				"salary":   95000.00,
				"status":   "active",
			},
		},
	}

	// Comprehensive validator that collects all validation errors
	validationErrors := make(map[string][]string)

	comprehensiveValidator := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		if record, ok := o.(map[string]interface{}); ok {
			recordID := fmt.Sprintf("record_%v", record["id"])
			errors := []string{}

			// Validate name
			if name, exists := record["name"]; exists {
				if nameStr, ok := name.(string); !ok || strings.TrimSpace(nameStr) == "" {
					errors = append(errors, "name is empty or invalid")
				}
			}

			// Validate email
			if email, exists := record["email"]; exists {
				if emailStr, ok := email.(string); ok {
					emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
					if !emailRegex.MatchString(emailStr) {
						errors = append(errors, "email format is invalid")
					}
				}
			}

			// Validate age
			if age, exists := record["age"]; exists {
				if ageNum, ok := age.(float64); !ok || ageNum < 0 || ageNum > 150 {
					errors = append(errors, "age must be between 0 and 150")
				}
			}

			// Validate salary
			if salary, exists := record["salary"]; exists {
				if _, ok := salary.(float64); !ok {
					errors = append(errors, "salary must be a number")
				}
			}

			// Validate status
			if status, exists := record["status"]; exists {
				if statusStr, ok := status.(string); ok {
					validStatuses := []string{"active", "inactive", "pending"}
					valid := false
					for _, validStatus := range validStatuses {
						if statusStr == validStatus {
							valid = true
							break
						}
					}
					if !valid {
						errors = append(errors, "status must be one of: active, inactive, pending")
					}
				}
			}

			// Store validation errors
			if len(errors) > 0 {
				validationErrors[recordID] = errors
				record["validation_status"] = "INVALID"
				record["validation_errors"] = errors
			} else {
				record["validation_status"] = "VALID"
			}
		}
		return o
	})

	navigator := &mapnavigator.MapNavigator{
		NodeModifier: comprehensiveValidator,
		ReadOnly:     false,
	}

	// Validate all records
	if result, err := navigator.VisitNode(data, "records", "*", "-"); err == nil {
		fmt.Printf("Validated %d records\n", len(result.([]interface{})))
	}

	// Print validation summary
	fmt.Printf("\nValidation Summary:\n")
	for recordID, errors := range validationErrors {
		fmt.Printf("- %s: %d errors\n", recordID, len(errors))
		for _, error := range errors {
			fmt.Printf("  * %s\n", error)
		}
	}

	// Extract validation statuses
	readOnlyNav := &mapnavigator.MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			return o
		}),
		ReadOnly: true,
	}

	if result, err := readOnlyNav.VisitNode(data, "records", "*", "validation_status"); err == nil {
		fmt.Printf("\nValidation statuses: %v\n", result)
	}

	// Print the validated data
	if jsonData, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Printf("\nValidated data:\n%s\n", jsonData)
	}
}