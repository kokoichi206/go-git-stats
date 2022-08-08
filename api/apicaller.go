package api

// Github REST Api caller
// For detailed information, see the official documentations:
// https://docs.github.com/ja/rest
type ApiCaller interface {
	ListPublicRepositories(userName string) ([]Repository, error)
	ListRepositoriesForAuthenticatedUser() ([]Repository, error)
}
