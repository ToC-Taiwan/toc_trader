// Package pyresponse package pyapiserver
package pyresponse

// PyServerResponse PyServerResponse
type PyServerResponse struct {
	Status  string `json:"status"`
	OrderID string `json:"order_id"`
}

// SuccessStatus SuccessStatus
const SuccessStatus string = "success"
