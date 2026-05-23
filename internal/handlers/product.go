package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

// ProductHandler injects the database connection into product-related HTTP handlers.
type ProductHandler struct {
	DB *pgx.Conn
}

// GetProducts retrieves all available items from the database catalog.
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	// 1. Query all products, ordered by ID
	rows, err := h.DB.Query(context.Background(), "SELECT id, sku, name, price, stock FROM products ORDER BY id ASC")
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch product catalog"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// 2. Map the SQL rows into a Go slice (array) of maps
	var products []map[string]interface{}
	for rows.Next() {
		var id, stock int
		var sku, name string
		var price float64

		// Scan row data into variables
		if err := rows.Scan(&id, &sku, &name, &price, &stock); err != nil {
			continue
		}

		// Append to the products slice
		products = append(products, map[string]interface{}{
			"id": id, "sku": sku, "name": name, "price": price, "stock": stock,
		})
	}

	// 3. Send the JSON response to the Frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Checkout processes a transaction and deducts the inventory stock securely.
func (h *ProductHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	// 1. Define the expected JSON payload struct from the Frontend
	var payload struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "Invalid checkout payload"}`, http.StatusBadRequest)
		return
	}

	// 2. Execute stock deduction ONLY IF the current stock is sufficient (Atomic operation)
	// This prevents the stock from becoming negative!
	cmdTag, err := h.DB.Exec(context.Background(),
		"UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1",
		payload.Quantity, payload.ProductID)

	// 3. If an error occurs or no rows were affected (meaning stock was insufficient)
	if err != nil || cmdTag.RowsAffected() == 0 {
		http.Error(w, `{"error": "Checkout failed: Insufficient stock or invalid product ID"}`, http.StatusConflict)
		return
	}

	// 4. Send success confirmation
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Transaction processed successfully"}`))
}

// GenerateDailyReport simulates processing thousands of records concurrently using Goroutines.
func (h *ProductHandler) GenerateDailyReport(w http.ResponseWriter, r *http.Request) {
	// 1. Record the start time to measure execution speed
	startTime := time.Now()

	// 2. Simulate 10 heavy background tasks (e.g., generating PDFs, sending emails)
	// In a synchronous language, 10 tasks * 1 second = 10 seconds total.
	totalTasks := 10

	// 3. Initialize a WaitGroup to act as the "supervisor" for our Goroutines
	var wg sync.WaitGroup

	// 4. Create a buffered channel to safely collect results from multiple concurrent workers
	results := make(chan string, totalTasks)

	// 5. Dispatch the worker Goroutines
	for i := 1; i <= totalTasks; i++ {
		wg.Add(1) // Notify the WaitGroup that a new worker is starting

		// The 'go' keyword spins up a new concurrent thread instantly
		go func(taskID int) {
			defer wg.Done() // Notify the WaitGroup when this worker finishes

			// Simulate a heavy operation taking exactly 1 second
			time.Sleep(1 * time.Second)

			// Send the result safely into the channel
			results <- fmt.Sprintf("✅ Transaction report #%d generated successfully", taskID)
		}(i)
	}

	// 6. Block the main thread until ALL workers report that they are done
	wg.Wait()
	close(results) // Close the channel since no more data will be sent

	// 7. Collect all processed results from the channel
	var reportDetails []string
	for res := range results {
		reportDetails = append(reportDetails, res)
	}

	// 8. Calculate the total execution time
	duration := time.Since(startTime)

	// 9. Send the final compiled report back to the Frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":         "Enterprise Daily Report Completed!",
		"total_processed": len(reportDetails),
		"time_taken_ms":   duration.Milliseconds(),
		"details":         reportDetails,
	})
}
