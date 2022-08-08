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
		URL         = fmt.Sprintf("https://api.github.com/repos/%s/stats/code_frequency", fullName)
		retries     = 3
		waitSeconds = 3 * time.Second
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	// Set request header.
	req.Header.Add("Accept", "application/vnd.github+json")
	if a.config.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("token %s", a.config.Token))
	}

	var resp *http.Response
	for retries > 0 {
		resp, err = client.Do(req)

		// If the StatusCode starts with 4, it is user's error,
		// so it should not be retried.
		if resp.StatusCode/100 == 4 {
			return nil, err
		}

		if resp.StatusCode == 202 {
			// Failed to find cache and GitHub started to create statistics.
			// Sleep some time and retry.
			time.Sleep(waitSeconds)
		} else if err != nil {
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

	var cf [][]int
	if err := json.Unmarshal(body, &cf); err != nil {
		return nil, err
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
