// Package global package global
package global

// PyServerResponse PyServerResponse
type PyServerResponse struct {
	Status  string `json:"status"`
	OrderID string `json:"order_id"`
}

// SuccessStatus SuccessStatus
const SuccessStatus string = "success"
