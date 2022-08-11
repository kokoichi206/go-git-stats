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
		URL     = fmt.Sprintf("%s/users/%s/repos?per_page=100", a.config.ApiBaseURL, userName)
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
	success := false
	for retries > 0 {
		// > If the returned error is nil, the Response will contain a non-nil
		// > Body which the user is expected to close.
		resp, err = client.Do(req)

		if err != nil {
			// Invalid URL (Like different scheme) etc.
			retries -= 1
			continue
		}

		if resp.StatusCode == http.StatusOK {
			// Success!
			success = true
			break
		}

		if resp.StatusCode/100 == 4 {
			// If the StatusCode starts with 4, it is user's error,
			// so it should not be retried.
			return nil, fmt.Errorf("failed to client.Do: StatusCode is %d", resp.StatusCode)
		}

		retries -= 1
	}

	if !success {
		return nil, fmt.Errorf("failed to client.Do after several retries.")
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
		URL     = fmt.Sprintf("%s/user/repos?per_page=100", a.config.ApiBaseURL)
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
	success := false
	for retries > 0 {
		// > If the returned error is nil, the Response will contain a non-nil
		// > Body which the user is expected to close.
		resp, err = client.Do(req)

		if err != nil {
			// Invalid URL (Like different scheme) etc.
			retries -= 1
			continue
		}

		if resp.StatusCode == http.StatusOK {
			// Success!
			success = true
			break
		}

		if resp.StatusCode/100 == 4 {
			// If the StatusCode starts with 4, it is user's error,
			// so it should not be retried.
			return nil, fmt.Errorf("failed to client.Do: StatusCode is %d", resp.StatusCode)
		}

		retries -= 1
	}

	if !success {
		return nil, fmt.Errorf("failed to client.Do after several retries.")
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
