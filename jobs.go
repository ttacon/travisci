package travisci

import (
	"encoding/json"
	"fmt"
	"time"
)

type Job struct {
	ID            int         `json:"id"`
	BuildID       int         `json:"build_id"`
	RepositoryID  int         `json:"repository_id"`
	CommitID      int         `json:"commit_id"`
	LogID         int         `json:"log_id"`
	Number        string      `json:"number"`
	Config        *Config     `json:"config"`
	State         string      `json:"state"`
	StartedAt     *time.Time  `json:"started_at"`
	FinishedAt    *time.Time  `json:"finished_at"`
	Duration      int         `json:"duration"`
	Queue         string      `json:"queue"`
	AllowFailure  bool        `json:"allow_failure"`
	Tags          interface{} `json:"tags"` // no idea what this field is
	AnnotationIDs []int       `json:"annotation_ids"`
}

type JobResponse struct {
	Job         *Job          `json:"job"`
	Commit      *Commit       `json:"commit"`
	Annotations []interface{} `json:"annotations"`
}

func (c *client) GetJobByID(id int) (*Job, error) {
	req, err := newReq("GET", fmt.Sprintf("/jobs/%d", id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "token "+c.travisToken)
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	var data JobResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Job, err
}
