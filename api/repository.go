package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Lists public repositories for a user.
// See documentation:
// https://docs.github.com/ja/rest/repos/repos#list-public-repositories
func (a *Api) ListPublicRepositories(userName string) ([]Repository, error) {

	var (
		URL     = fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100", userName)
		retries = 3
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to http.NewRequest: %w", err)
	}

	// Set request header.
	req.Header.Add("Accept", "application/vnd.github+json")

	var resp *http.Response
	for retries > 0 {
		resp, err = client.Do(req)

		// If the StatusCode starts with 4, it is user's error,
		// so it should not be retried.
		if resp.StatusCode/100 == 4 {
			return nil, fmt.Errorf("failed to client.Do: StatusCode is %d", resp.StatusCode)
		}

		if err != nil {
			retries -= 1
		} else {
			// Success!
			break
		}
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to io.ReadAll: %w", err)
	}

	var repositories []Repository
	if err := json.Unmarshal(body, &repositories); err != nil {
		return nil, fmt.Errorf("failed to json.Unmarshal: %w", err)
	}

	return repositories, nil
}

// Lists repositories for the authenticated user.
// Config must have the github access token.
// See documentation:
// https://docs.github.com/ja/rest/repos/repos#list-repositories-for-the-authenticated-user
func (a *Api) ListRepositoriesForAuthenticatedUser() ([]Repository, error) {

	var (
		URL     = "https://api.github.com/user/repos?per_page=100"
		retries = 3
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to http.NewRequest: %w", err)
	}

	// Set request header.
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("token %s", a.config.Token))

	var resp *http.Response
	for retries > 0 {
		resp, err = client.Do(req)

		// If the StatusCode starts with 4, it is user's error,
		// so it should not be retried.
		if resp.StatusCode/100 == 4 {
			return nil, fmt.Errorf("failed to client.Do: StatusCode is %d", resp.StatusCode)
		}

		if err != nil {
			retries -= 1
		} else {
			// Success!
			break
		}
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to io.ReadAll: %w", err)
	}

	var repositories []Repository
	if err := json.Unmarshal(body, &repositories); err != nil {
		return nil, fmt.Errorf("failed to json.Unmarshal: %w", err)
	}

	return repositories, nil
}
