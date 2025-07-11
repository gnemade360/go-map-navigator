package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/conditions"
	models "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-nav-models"
	actionModels "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-action-repository/models"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-conditional-modifier"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/templates"
)

func main() {
	// Example 1: Basic conditional processing
	fmt.Println("=== Basic Conditional Processing ===")
	basicConditionalExample()

	// Example 2: Complex conditional logic
	fmt.Println("\n=== Complex Conditional Logic ===")
	complexConditionalExample()

	// Example 3: Conditional data transformation
	fmt.Println("\n=== Conditional Data Transformation ===")
	conditionalTransformationExample()

	// Example 4: Multi-condition processing
	fmt.Println("\n=== Multi-condition Processing ===")
	multiConditionExample()
}

func basicConditionalExample() {
	// Sample user data
	data := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{
				"id":       1,
				"name":     "Alice",
				"age":      25,
				"status":   "active",
				"premium":  true,
				"balance":  1000.50,
			},
			map[string]interface{}{
				"id":       2,
				"name":     "Bob",
				"age":      17,
				"status":   "inactive",
				"premium":  false,
				"balance":  50.00,
			},
			map[string]interface{}{
				"id":       3,
				"name":     "Charlie",
				"age":      30,
				"status":   "active",
				"premium":  true,
				"balance":  2500.00,
			},
		},
	}

	// Create conditional modifier: if age >= 18, mark as adult
	config := &actionModels.MapNodeConditionalModifierConfig{
		Condition: "{{.NodeValue}} >= 18",
		TemplateConfigHolder: &actionModels.TemplateConfigHolder{
			TemplateConfig: templates.TemplateConfig{
				LeftDelim:  "{{",
				RightDelim: "}}",
			},
		},
	}

	conditionalModifier := &map_node_conditional_modifier.MapNodeConditionalModifier{
		Config:             config,
		ConditionsExecutor: conditions.NewSaConditionsExecutor(),
		IfTrue: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			return "adult"
		}),
		IfFalse: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			return "minor"
		}),
	}

	navigator := &mapnavigator.MapNavigator{
		NodeModifier: conditionalModifier,
		ReadOnly:     false,
	}

	// Apply conditional logic to all ages
	if result, err := navigator.VisitNode(data, "users", "*", "age"); err == nil {
		fmt.Printf("Age classifications: %v\n", result)
	}

	// Print the modified data
	if jsonData, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Printf("Modified data:\n%s\n", jsonData)
	}
}

func complexConditionalExample() {
	// E-commerce order data
	data := map[string]interface{}{
		"orders": []interface{}{
			map[string]interface{}{
				"id":       "order-1",
				"amount":   150.00,
				"status":   "pending",
				"customer": "premium",
				"items":    3,
			},
			map[string]interface{}{
				"id":       "order-2",
				"amount":   75.00,
				"status":   "completed",
				"customer": "regular",
				"items":    1,
			},
			map[string]interface{}{
				"id":       "order-3",
				"amount":   300.00,
				"status":   "pending",
				"customer": "premium",
				"items":    5,
			},
		},
	}

	// Create a complex conditional modifier for order priority
	// Priority logic: premium customers with amount > 100 get "high" priority
	priorityModifier := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		if orderMap, ok := o.(map[string]interface{}); ok {
			amount, hasAmount := orderMap["amount"].(float64)
			customer, hasCustomer := orderMap["customer"].(string)
			status, hasStatus := orderMap["status"].(string)

			if hasAmount && hasCustomer && hasStatus {
				if customer == "premium" && amount > 100 && status == "pending" {
					orderMap["priority"] = "high"
				} else if customer == "premium" && amount > 50 {
					orderMap["priority"] = "medium"
				} else {
					orderMap["priority"] = "low"
				}
			}
		}
		return o
	})

	navigator := &mapnavigator.MapNavigator{
		NodeModifier: priorityModifier,
		ReadOnly:     false,
	}

	// Apply priority logic to all orders
	if result, err := navigator.VisitNode(data, "orders", "*", "-"); err == nil {
		fmt.Printf("Processed %d orders\n", len(result.([]interface{})))
	}

	// Extract all priorities
	readOnlyNav := &mapnavigator.MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			return o
		}),
		ReadOnly: true,
	}

	if result, err := readOnlyNav.VisitNode(data, "orders", "*", "priority"); err == nil {
		fmt.Printf("Order priorities: %v\n", result)
	}

	// Print the modified data
	if jsonData, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Printf("Orders with priorities:\n%s\n", jsonData)
	}
}

func conditionalTransformationExample() {
	// Product inventory data
	data := map[string]interface{}{
		"inventory": []interface{}{
			map[string]interface{}{
				"sku":       "LAPTOP-001",
				"name":      "Gaming Laptop",
				"price":     1299.99,
				"stock":     5,
				"category":  "electronics",
				"discount":  0.0,
			},
			map[string]interface{}{
				"sku":       "BOOK-001",
				"name":      "Programming Guide",
				"price":     49.99,
				"stock":     0,
				"category":  "books",
				"discount":  0.0,
			},
			map[string]interface{}{
				"sku":       "PHONE-001",
				"name":      "Smartphone",
				"price":     699.99,
				"stock":     15,
				"category":  "electronics",
				"discount":  0.0,
			},
		},
	}

	// Conditional discount application
	discountModifier := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		if product, ok := o.(map[string]interface{}); ok {
			price, hasPrice := product["price"].(float64)
			stock, hasStock := product["stock"].(float64)
			category, hasCategory := product["category"].(string)

			if hasPrice && hasStock && hasCategory {
				// Apply different discounts based on conditions
				if stock == 0 {
					product["discount"] = 0.0 // No discount for out of stock
					product["status"] = "out_of_stock"
				} else if category == "electronics" && price > 500 {
					product["discount"] = 0.15 // 15% discount for expensive electronics
					product["status"] = "premium_discount"
				} else if category == "books" && stock < 10 {
					product["discount"] = 0.10 // 10% discount for low stock books
					product["status"] = "clearance"
				} else {
					product["discount"] = 0.05 // 5% standard discount
					product["status"] = "regular"
				}

				// Calculate final price
				finalPrice := price * (1 - product["discount"].(float64))
				product["final_price"] = finalPrice
			}
		}
		return o
	})

	navigator := &mapnavigator.MapNavigator{
		NodeModifier: discountModifier,
		ReadOnly:     false,
	}

	// Apply discount logic to all products
	if result, err := navigator.VisitNode(data, "inventory", "*", "-"); err == nil {
		fmt.Printf("Processed %d products\n", len(result.([]interface{})))
	}

	// Extract product statuses
	readOnlyNav := &mapnavigator.MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			return o
		}),
		ReadOnly: true,
	}

	if result, err := readOnlyNav.VisitNode(data, "inventory", "*", "status"); err == nil {
		fmt.Printf("Product statuses: %v\n", result)
	}

	if result, err := readOnlyNav.VisitNode(data, "inventory", "*", "final_price"); err == nil {
		fmt.Printf("Final prices: %v\n", result)
	}

	// Print the modified data
	if jsonData, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Printf("Inventory with discounts:\n%s\n", jsonData)
	}
}

func multiConditionExample() {
	// Employee performance data
	data := map[string]interface{}{
		"employees": []interface{}{
			map[string]interface{}{
				"id":          1,
				"name":        "Alice",
				"department":  "engineering",
				"years":       5,
				"performance": 4.5,
				"salary":      75000,
			},
			map[string]interface{}{
				"id":          2,
				"name":        "Bob",
				"department":  "sales",
				"years":       2,
				"performance": 3.2,
				"salary":      50000,
			},
			map[string]interface{}{
				"id":          3,
				"name":        "Charlie",
				"department":  "engineering",
				"years":       8,
				"performance": 4.8,
				"salary":      90000,
			},
		},
	}

	// Multi-condition promotion eligibility
	promotionModifier := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		if employee, ok := o.(map[string]interface{}); ok {
			years, hasYears := employee["years"].(float64)
			performance, hasPerf := employee["performance"].(float64)
			department, hasDept := employee["department"].(string)
			salary, hasSalary := employee["salary"].(float64)

			if hasYears && hasPerf && hasDept && hasSalary {
				eligible := false
				bonusPercent := 0.0

				// Complex promotion logic
				if years >= 5 && performance >= 4.0 {
					eligible = true
					bonusPercent = 0.15
				} else if years >= 3 && performance >= 4.5 {
					eligible = true
					bonusPercent = 0.10
				} else if department == "sales" && performance >= 3.5 {
					eligible = true
					bonusPercent = 0.08
				}

				employee["promotion_eligible"] = eligible
				employee["bonus_percent"] = bonusPercent

				if eligible {
					employee["new_salary"] = salary * (1 + bonusPercent)
					employee["promotion_reason"] = fmt.Sprintf("Years: %.0f, Performance: %.1f, Dept: %s", years, performance, department)
				}
			}
		}
		return o
	})

	navigator := &mapnavigator.MapNavigator{
		NodeModifier: promotionModifier,
		ReadOnly:     false,
	}

	// Apply promotion logic to all employees
	if result, err := navigator.VisitNode(data, "employees", "*", "-"); err == nil {
		fmt.Printf("Evaluated %d employees\n", len(result.([]interface{})))
	}

	// Extract promotion results
	readOnlyNav := &mapnavigator.MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			return o
		}),
		ReadOnly: true,
	}

	if result, err := readOnlyNav.VisitNode(data, "employees", "*", "promotion_eligible"); err == nil {
		fmt.Printf("Promotion eligibility: %v\n", result)
	}

	if result, err := readOnlyNav.VisitNode(data, "employees", "*", "bonus_percent"); err == nil {
		fmt.Printf("Bonus percentages: %v\n", result)
	}

	// Print the modified data
	if jsonData, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Printf("Employee evaluation results:\n%s\n", jsonData)
	}
}