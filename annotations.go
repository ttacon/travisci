package travisci

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Annotation struct {
	ID          int    `json:"id"`
	JobID       int    `json:"job_id"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Status      string `json:"status"`
}

func (c client) CreateAnnotation(jobID int, annotation Annotation) error {
	// NOTE(ttacon): this currently fails - 401, unsure what username and secret
	// need to be (what does it mean by provider - Travis, GH?)
	req, err := newReq(
		"POST",
		fmt.Sprintf("/jobs/%d/annotations", jobID),
		map[string]interface{}{
			"url":         annotation.URL,
			"description": annotation.Description,
			"status":      annotation.Status,
		},
	)
	if err != nil {
		return err
	}

	if len(c.travisToken) > 0 {
		req.Header.Add("Authorization", "token "+c.travisToken)
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}

	fmt.Printf("%#v\n", resp)
	dbytes, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(dbytes))
	fmt.Println("err[in]: ", err)

	return nil
}

type annotationList struct {
	Annotations []Annotation `json:"annotations"`
}

func (c client) ListAnnotations(jobID int) ([]Annotation, error) {
	req, err := newReq("GET",
		fmt.Sprintf("/jobs/%d/annotations", jobID),
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

	var data annotationList
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Annotations, nil
}
