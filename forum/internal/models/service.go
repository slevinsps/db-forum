package models

type Service struct {
	Forum  int `json:"forum,omitempty"`
	Post   int `json:"post,omitempty"`
	Thread int `json:"thread,omitempty"`
	User   int `json:"user,omitempty"`
}
