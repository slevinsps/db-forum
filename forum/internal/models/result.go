package models

type Result struct {
	Place   string `json:"place,omitempty"`
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}
