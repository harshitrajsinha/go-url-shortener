package common

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Reference string      `json:"reference,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}