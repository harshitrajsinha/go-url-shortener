package common

type Response struct {
	Code      int                      `json:"code"`
	Message   string                   `json:"message"`
	Reference string                   `json:"reference,omitempty"`
	Data      []map[string]interface{} `json:"data,omitempty"`
}