package travisci

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	acceptHeader = "application/vnd.travis-ci.2+json"
	userAgent    = "GoTravisCI/0.0.1"
	baseUrl      = "https://api.travis-ci.org"
)

type Client interface {
	GetRepo(string) (*Repo, error)
	GetJobByID(int) (*Job, error)
	LogByID(int) (*Log, error)
	ArchivedLogByJob(int) ([]byte, error)

	// Accounts
	Accounts(bool) ([]Account, error)

	// Annotations
	CreateAnnotation(int, Annotation) error
	ListAnnotations(int) ([]Annotation, error)

	// Branches
	ListBranches(string) ([]Branch, error)
	GetBranch(string, string) (Branch, error)

	// Broadcasts
	ListBroadcasts() ([]Broadcast, error)

	// Builds
	GetBuildByID(string) (*Build, error)
	ListBuildsForRepo(string) ([]*Build, error)
	CancelBuild(int) error
	RestartBuild(int) error
}

type client struct {
	travisToken string
	c           *http.Client
}

func NewClientFromGH(ghToken string) (Client, error) {
	req, err := newReq("POST", "/auth/github", map[string]string{
		"github_token": ghToken,
	})
	if err != nil {
		return nil, err
	}

	// TODO(ttacon): add appropriate round tripper
	c := &http.Client{}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	db, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data = make(map[string]string)
	err = json.Unmarshal(db, &data)
	if err != nil {
		return nil, err
	}

	tok, ok := data["access_token"]
	if !ok {
		return nil, errors.New("token was not processed")
	}

	return &client{
		travisToken: tok,
		c:           c,
	}, nil
}

func NewClientFromTravis(tok string) Client {
	return &client{
		travisToken: tok,
		c:           &http.Client{},
	}
}

func newReq(method, ep string, body interface{}) (*http.Request, error) {
	// this method is based off
	// https://github.com/google/go-github/blob/master/github/github.go:
	// NewRequest as it's a very nice way of doing this
	_, err := url.Parse(ep)
	if err != nil {
		return nil, err
	}

	// This is useful as this functionality works the same for the actual
	// BASE_URL and the download url (TODO(ttacon): insert download url)
	// this seems to be failing to work not RFC3986 (url resolution)
	//	resolvedUrl := c.BaseUrl.ResolveReference(parsedUrl)
	resolvedUrl, err := url.Parse(baseUrl + ep)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if body != nil {
		if err = json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, resolvedUrl.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Accept", acceptHeader)
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}
