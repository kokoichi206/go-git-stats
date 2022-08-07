package api

// Github REST Api caller
type ApiCaller interface {
	ListRepositoriesForAuthenticatedUser() ([]Repository, error)
}
