package models

import (
	"fmt"
	"strings"
)

// ExampleMapNodeModifierFunc demonstrates using a function as a modifier.
func ExampleMapNodeModifierFunc() {
	// Create a modifier function that converts strings to uppercase
	uppercaseModifier := MapNodeModifierFunc(func(node interface{}) interface{} {
		if str, ok := node.(string); ok {
			return strings.ToUpper(str)
		}
		return node
	})

	// Use the modifier
	result := uppercaseModifier.ModifyNode("hello world")
	fmt.Printf("Modified result: %s\n", result)

	// Test with non-string input
	result = uppercaseModifier.ModifyNode(42)
	fmt.Printf("Non-string input unchanged: %v\n", result)

	// Output:
	// Modified result: HELLO WORLD
	// Non-string input unchanged: 42
}

// Define a custom modifier that multiplies numbers by 2
type DoubleModifier struct{}

func (d DoubleModifier) ModifyNode(node interface{}) interface{} {
	switch v := node.(type) {
	case int:
		return v * 2
	case float64:
		return v * 2
	default:
		return node
	}
}

// ExampleMapNodeModifier demonstrates implementing a custom modifier.
func ExampleMapNodeModifier() {
	// Create and use the modifier
	modifier := DoubleModifier{}

	result := modifier.ModifyNode(5)
	fmt.Printf("Integer doubled: %v\n", result)

	result = modifier.ModifyNode(3.14)
	fmt.Printf("Float doubled: %v\n", result)

	result = modifier.ModifyNode("text")
	fmt.Printf("String unchanged: %v\n", result)

	// Output:
	// Integer doubled: 10
	// Float doubled: 6.28
	// String unchanged: text
}

// ExampleMapNodeModifier_chaining demonstrates chaining modifiers.
func ExampleMapNodeModifier_chaining() {
	// Create a modifier that adds 10 to numbers
	add10 := MapNodeModifierFunc(func(node interface{}) interface{} {
		if num, ok := node.(int); ok {
			return num + 10
		}
		return node
	})

	// Create a modifier that multiplies by 2
	multiply2 := MapNodeModifierFunc(func(node interface{}) interface{} {
		if num, ok := node.(int); ok {
			return num * 2
		}
		return node
	})

	// Chain the modifiers
	input := 5
	step1 := add10.ModifyNode(input)
	result := multiply2.ModifyNode(step1)

	fmt.Printf("Input: %v\n", input)
	fmt.Printf("After adding 10: %v\n", step1)
	fmt.Printf("After multiplying by 2: %v\n", result)

	// Output:
	// Input: 5
	// After adding 10: 15
	// After multiplying by 2: 30
}
