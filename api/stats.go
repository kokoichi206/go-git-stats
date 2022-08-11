package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Get the weekly commit activity of a specific repository.
// See documentation:
// https://docs.github.com/ja/rest/metrics/statistics#get-the-weekly-commit-activity
func (a *Api) WeeklyCommitActivity(fullName string) ([]CodeFrequency, error) {

	var (
		URL         = fmt.Sprintf("%s/repos/%s/stats/code_frequency", a.config.ApiBaseURL, fullName)
		retries     = 3
		waitSeconds = 3 * time.Second
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to http.NewRequest: %w", err)
	}

	// Set request header.
	req.Header.Add("Accept", "application/vnd.github+json")
	if a.config.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("token %s", a.config.Token))
	}

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

		if resp.StatusCode == http.StatusCreated {
			// Failed to find cache and GitHub started to create statistics.
			// Sleep some time and retry.
			//
			// See GitHub documentation: https://docs.github.com/en/rest/metrics/statistics#a-word-about-caching
			time.Sleep(waitSeconds)
			continue
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

	var cf [][]int
	if err := json.Unmarshal(body, &cf); err != nil {
		return nil, fmt.Errorf("failed to json.Unmarshal: %w", err)
	}

	// MAYBE: there's better way.
	// Convert array to CodeFrequency
	var codeFreqs []CodeFrequency
	for i := len(cf) - 1; i >= 0; i-- {
		c := cf[i]
		if len(c) == 3 {
			codeFreqs = append(codeFreqs, CodeFrequency{
				Time:      c[0],
				Additions: c[1],
				Deletions: c[2],
			})
		}
	}

	return codeFreqs, nil
}
