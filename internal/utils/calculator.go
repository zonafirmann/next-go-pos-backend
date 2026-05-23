package utils

// CalculateFinalPrice computes the grand total after applying a percentage discount and tax.
// It returns the final amount that the customer needs to pay.
func CalculateFinalPrice(subtotal float64, discountPercent float64, taxPercent float64) float64 {
	// 1. Calculate discount deduction
	discountAmount := subtotal * (discountPercent / 100)
	priceAfterDiscount := subtotal - discountAmount
	taxAmount := priceAfterDiscount * (taxPercent / 100)
	// 3. Return final grand total
	return priceAfterDiscount + taxAmount
}
