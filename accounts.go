package travisci

import "encoding/json"

type Account struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Login      string `json:"login"`
	Type       string `json:"type"`
	ReposCount int    `json:"repos_count"`
	Subscribed bool   `json:"subscribed"`
}

type accountList struct {
	Accounts []Account `json:"accounts"`
}

func (c client) Accounts(includeNonAdmin bool) ([]Account, error) {
	req, err := newReq("GET", "/accounts", map[string]bool{
		"all": includeNonAdmin,
	})
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

	var data accountList
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data.Accounts, nil
}
