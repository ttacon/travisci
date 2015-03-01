package travisci

import (
	"encoding/json"
	"fmt"
)

type Branch struct {
	ID           int                    `json:"id"`
	RepositoryID int                    `json:"repository_id"`
	CommitID     int                    `json:"commit_id"`
	Number       string                 `json:"number"`
	Config       map[string]interface{} `json:"config"`
	State        string                 `json:"state"`
	StartedAt    string                 `json:"started_at"`
	FinishedAt   string                 `json:"finished_at"`
	Duration     int                    `json:"duration"`
	JobIDs       []int                  `json:"job_ids"`
	PullRequest  bool                   `json:"pull_request"`
}

type branchList struct {
	Branches []Branch `json:"branches"`
}

func (c client) ListBranches(identifier string) ([]Branch, error) {
	// NOTE(ttacon): at the top level there are two fields: "branches" and "commits"
	// should we parse out commits?
	req, err := newReq("GET",
		fmt.Sprintf("/repos/%s/branches", identifier),
		nil,
	)
	if err != nil {
		return nil, err
	}

	if len(c.travisToken) > 0 {
		req.Header.Add("Authorization", "token "+c.travisToken)
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	var data branchList
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Branches, nil
}

type getBranchResponse struct {
	Branch Branch `json:"branch"`
}

func (c client) GetBranch(repoIdentifier, branch string) (Branch, error) {
	var empty Branch
	req, err := newReq("GET",
		fmt.Sprintf("/repos/%s/branches/%s", repoIdentifier, branch),
		nil,
	)
	if err != nil {
		return empty, err
	}

	if len(c.travisToken) > 0 {
		req.Header.Add("Authorization", "token "+c.travisToken)
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return empty, err
	}

	var data getBranchResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return empty, err
	}

	return data.Branch, nil
}
