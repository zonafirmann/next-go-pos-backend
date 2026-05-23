package utils

import (
	"testing"
)

// TestCalculateFinalPrice validates the mathematical accuracy of our POS calculator.
func TestCalculateFinalPrice(t *testing.T) {
	// Define a struct for our test cases (Table-Driven Testing approach)
	tests := []struct {
		name            string
		subtotal        float64
		discountPercent float64
		taxPercent      float64
		expected        float64
	}{
		{
			name:            "Normal transaction with 0% discount and 10% tax",
			subtotal:        100000,
			discountPercent: 0,
			taxPercent:      10,
			expected:        110000, // 100k + 10k tax
		},
		{
			name:            "Transaction with 20% discount and 10% tax",
			subtotal:        100000,
			discountPercent: 20,
			taxPercent:      10,
			expected:        88000, // 100k - 20k discount = 80k. 80k + 8k tax = 88k
		},
		{
			name:            "Free transaction (100% discount)",
			subtotal:        500000,
			discountPercent: 100,
			taxPercent:      11,
			expected:        0, // Completely free, no tax should apply
		},
	}

	// Loop through each scenario and run the test
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CalculateFinalPrice(tc.subtotal, tc.discountPercent, tc.taxPercent)

			// If the machine's calculation doesn't match our human expectation, trigger an error
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}
