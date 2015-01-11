package travisci

import (
	"encoding/json"
	"fmt"
	"time"
)

type Repo struct {
	ID                  int        `json:"id"`
	Slug                string     `json:"slug"`
	Description         string     `json:"description"`
	LastBuildID         int64      `json:"last_build_id"`
	LastBuildNumber     string     `json:"last_build_number"`
	LastBuildState      string     `json:"last_build_state"`
	LastBuildDuration   int        `json:"last_build_duration"`
	LastBuildStartedAt  *time.Time `json:"last_build_started_at"`
	LastBuildFinishedAt *time.Time `json:"last_build_finished_at"`
	GitHubLanguage      string     `json:"github_language"`
}

func (c *client) GetRepo(id string) (*Repo, error) {
	req, err := newReq("GET", fmt.Sprintf("/repos/%s", id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "token "+c.travisToken)
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	// weird, there's an undocumented field "last_build_language"
	// that comes back as null
	var data = make(map[string]*Repo)
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data["repo"], err
}
