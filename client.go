package targetprocess

import (
	"encoding/xml"
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
}

func NewClient(opt ClientOpt) TPClient {
	opt.isBasicAuth = len(opt.User) != 0 && len(opt.Password) != 0
	opt.isAccessToken = len(opt.AccessToken) != 0

	return TPClient{
		opt:        opt,
		httpClient: http.DefaultClient,
	}
}

func (t *TPClient) Bugs() *BugsService {
	return &BugsService{
		client: t,
	}
}

func (t *TPClient) Url() string {
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
