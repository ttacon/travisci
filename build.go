package travisci

import (
	"encoding/json"
	"errors"
	"time"
)

type BuildResponse struct {
	Builds  []*Build  `json:"builds"`
	Commits []*Commit `json:"commits"`
}

type Build struct {
	ID                int        `json:"id"`
	RepositoryID      int        `json:"repository_id"`
	CommitID          int        `json:"commit_id"`
	Number            string     `json:"number"`
	PullRequest       bool       `json:"pull_request"`
	PullRequestTitle  *string    `json:"pull_request_title"`
	PullRequestNumber *int       `json:"pull_request_number"`
	Config            *Config    `json:"config"`
	State             string     `json:"state"`
	StartedAt         *time.Time `json:"started_at"`
	FinishedAt        *time.Time `json:"finished_at"`
	Duration          int        `json:"duration"`
	JobIDs            []int      `json:"job_ids"`
}

type Commit struct {
	ID                int        `json:"id"`
	SHA               string     `json:"sha"`
	Branch            string     `json:"branch"`
	Message           string     `json:"message"`
	CommittedAt       *time.Time `json:"committed_at"`
	AuthorName        string     `json:"author_name"`
	AuthorEmail       string     `json:"author_email"`
	CommitterName     string     `json:"committer_name"`
	CommitterEmail    string     `json:"committer_email"`
	CompareUrl        string     `json:"compare_url"`
	PullRequestNumber *int       `json:"pull_request_number"`
}

type Config struct {
	Language string `json:"language"`
	// TODO(ttacon): figure out how to find language and versions used
	// without putting them all in here, ignore for now, but put go on
	// for now since it's the best language
	Go           interface{} `json:"go"` // gross but deal w/ it, ty Travis :( eg: [1.3, "tip"]
	Services     []string    `json:"services"`
	BeforeScript []string    `json:"before_script"`
	Result       string      `json:".result"` // TODO(ttacon): will this work?
	OS           string      `json:"os"`
}

func (c *client) GetBuildByID(id string) (*Build, error) {
	req, err := newReq("GET", "/builds", map[string][]string{
		"ids": []string{id},
	})
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "token "+c.travisToken)
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	var data BuildResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if len(data.Builds) == 0 {
		return nil, errors.New("no builds found")
	}
	return data.Builds[0], nil
}
