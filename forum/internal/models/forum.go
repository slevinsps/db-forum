package models

type Forum struct {
	Posts   int    `json:"posts,omitempty"`
	Slug    string `json:"slug,omitempty"`
	Threads int    `json:"threads,omitempty"`
	Title   string `json:"title,omitempty"`
	User    string `json:"user,omitempty"`
}
