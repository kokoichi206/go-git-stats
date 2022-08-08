package mock

import (
	"github.com/kokoichi206/go-git-stats/api"
	"github.com/kokoichi206/go-git-stats/util"
)

type MockApi struct {
	config              util.Config
	ListRepos           []api.Repository
	ListCodeFreq        []api.CodeFrequency
	Error               error
	PublicCalled        bool
	AuthenticatedCalled bool
	WeeklyCodeCalled    bool
	PassedFullName      string
}

func (a *MockApi) InitMock() {
	a.Error = nil
	a.PublicCalled = false
	a.AuthenticatedCalled = false
	a.WeeklyCodeCalled = false
}

func (a *MockApi) ListPublicRepositories(userName string) ([]api.Repository, error) {
	a.PublicCalled = true
	return a.ListRepos, a.Error
}

func (a *MockApi) ListRepositoriesForAuthenticatedUser() ([]api.Repository, error) {
	a.AuthenticatedCalled = true
	return a.ListRepos, a.Error
}

func (a *MockApi) WeeklyCommitActivity(fullName string) ([]api.CodeFrequency, error) {
	a.WeeklyCodeCalled = true
	a.PassedFullName = fullName
	return a.ListCodeFreq, a.Error
}

func New(config util.Config) *MockApi {
	return &MockApi{
		config:              config,
		Error:               nil,
		PublicCalled:        false,
		AuthenticatedCalled: false,
	}
}
