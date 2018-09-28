package targetprocess

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
)

type ClientOpt struct {
	User        string
	Password    string
	Token       string
	AccessToken string
	Url         string

	isBasicAuth   bool
	isAccessToken bool
	isToken       bool
}

type TPClient struct {
	httpClient *http.Client
	opt        ClientOpt

	bugs  *BugsService
	story *UserStoryService
}

func NewClient(opt ClientOpt) TPClient {
	opt.isBasicAuth = len(opt.User) != 0 && len(opt.Password) != 0
	opt.isAccessToken = len(opt.AccessToken) != 0

	cli := TPClient{
		opt:        opt,
		httpClient: http.DefaultClient,
	}

	cli.bugs = &BugsService{
		client: &cli,
	}

	cli.story = &UserStoryService{
		client: &cli,
	}

	return cli
}

func (t *TPClient) Bugs() *BugsService {
	return t.bugs
}

func (t *TPClient) Story() *UserStoryService {
	return t.story
}

func (t TPClient) Url() string {
	return t.opt.Url
}

func (t *TPClient) do(req *http.Request) (*http.Response, error) {
	if t.opt.isBasicAuth {
		req.SetBasicAuth(t.opt.User, t.opt.Password)
	}

	return t.httpClient.Do(req)
}

func (t *TPClient) prepare(values *url.Values) {
	if t.opt.isAccessToken {
		values.Add(`access_token`, t.opt.AccessToken)
	}
}

func (t *TPClient) extractError(response *http.Response) error {
	if response.StatusCode == http.StatusUnauthorized {
		return Unauthorized
	}

	dec := xml.NewDecoder(response.Body)
	var tpErr Error
	if err := dec.Decode(&tpErr); err != nil {
		return err
	}

	return tpErr
}

func (t *TPClient) get(endpoint string, values url.Values, receiver interface{}) error {
	t.prepare(&values)
	url := fmt.Sprintf("%s%s?%s", t.Url(), endpoint, values.Encode())

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := t.do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return t.extractError(resp)
	}
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(receiver)
}
