package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
		return nil, err
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
			return nil, err
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
		return nil, err
	}

	var repositories []Repository
	if err := json.Unmarshal(body, &repositories); err != nil {
		return nil, err
	}

	return repositories, nil
}
