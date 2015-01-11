package travisci

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Log struct {
	ID    int    `json:"id"`
	JobID int    `json:"job_id"`
	Body  string `json:"body"`
}

func (c *client) LogByID(id int) (*Log, error) {
	req, err := newReq("GET", fmt.Sprintf("/logs/%d", id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "token "+c.travisToken)
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	var data Log
	err = json.NewDecoder(resp.Body).Decode(&data)
	return &data, err
}

func (c *client) ArchivedLogByJob(id int) ([]byte, error) {
	req, err := newReq("GET", fmt.Sprintf("/jobs/%d/log", id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "token "+c.travisToken)
	req.Header.Add("Accept", "text/plain")
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}
