package api

type Repository struct {
	ID       int    `json:"id"`
	Private  bool   `json:"private"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}
