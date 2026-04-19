// Package models defines the data types shared across the backend.
// These structs are serialised to JSON and consumed by the React frontend.
package models

// Transaction represents a single credit card transaction extracted from a PDF statement.
type Transaction struct {
	Date        string `json:"date"`
	Description string `json:"description"`
	Amount      float64 `json:"amount"`
	Category    string `json:"category"`
	SubCategory string `json:"subCategory"`
	Hash        string `json:"hash"`
}

// ParseResponse is the per-file result returned by the /api/parse endpoint.
type ParseResponse struct {
	Transactions []Transaction `json:"transactions"`
	BillPeriod   BillPeriod    `json:"billPeriod"`
	CardName     string        `json:"cardName"`
	TotalAmount  float64       `json:"totalAmount"`
	FileName     string        `json:"fileName"`
}

// BillPeriod holds the statement's billing date range (e.g. "22/07/2024" to "22/08/2024").
type BillPeriod struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Category groups subcategories under a top-level spend category for the /api/categories endpoint.
type Category struct {
	Name          string   `json:"name"`
	SubCategories []string `json:"subCategories"`
}
