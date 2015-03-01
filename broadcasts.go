package travisci

import "encoding/json"

// http://docs.travis-ci.com/api/#broadcasts

type Broadcast struct {
	ID      int
	Message string
}

type broadcastList struct {
	Broadcasts []Broadcast `json:"broadcasts"`
}

func (c client) ListBroadcasts() ([]Broadcast, error) {
	req, err := newReq("GET", "/broadcasts", nil)
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

	var data broadcastList
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Broadcasts, nil
}
