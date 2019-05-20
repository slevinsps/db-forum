package models

type Post struct {
	Author   string `json:"author,omitempty"`
	Created  string `json:"created,omitempty"`
	Forum    string `json:"forum,omitempty"`
	Id       int    `json:"id,omitempty"`
	IsEdited bool   `json:"isEdited,omitempty"`
	Message  string `json:"message,omitempty"`
	Parent   int    `json:"parent,omitempty"`
	Thread   int    `json:"thread,omitempty"`
	Path     string `json:"path,omitempty"`
	Branch   string `json:"branch,omitempty"`
}

type PostArray struct {
    posts []Post `json:"posts,omitempty"`
}
