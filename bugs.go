package targetprocess

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	netUrl "net/url"
	"strconv"
)

const (
	bugEndpoint = `/api/v1/Bugs`
)

var (
	emptyTeams = []byte(`{}`)
)

type BugsService struct {
	client *TPClient
}

type BugsGetReply struct {
	Prev  string
	Next  string
	Items []Bug
}

func (b *BugsService) Search(query *Query) (BugsGetReply, error) {
	if query == nil {
		query = defaultQuery
	}

	var reply BugsGetReply
	values := query.values()
	b.client.prepare(&values)

	url := fmt.Sprintf("%s%s?%s", b.client.Url(), bugEndpoint, values.Encode())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return reply, err
	}

	resp, err := b.client.do(req)
	if err != nil {
		return reply, err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&reply)

	return reply, nil
}

func (b *BugsService) Get(id int64) (Bug, error) {
	values := defaultQuery.values()
	b.client.prepare(&values)

	var bug Bug

	url := fmt.Sprintf("%s%s/%d?%s", b.client.Url(), bugEndpoint, id, values.Encode())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return bug, err
	}

	resp, err := b.client.do(req)
	if err != nil {
		return bug, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&bug)
	return bug, err
}

func (b *BugsService) Next(url string) (BugsGetReply, error) {
	values := netUrl.Values{}
	b.client.prepare(&values)
	if len(values) != 0 {
		url = fmt.Sprintf("%s&%s", url, values.Encode())
	}

	var reply BugsGetReply
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return reply, err
	}

	resp, err := b.client.do(req)
	if err != nil {
		return reply, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&reply)
	return reply, err
}

type ProjectDescription struct {
	Id int
}

type Teams []int

type BugDescription struct {
	//Id          uint64 `json:"Id,omitempty"`
	Name        string
	Description string
	Project     ProjectDescription
	Teams       Teams `json:"AssignedTeams, omitempty"`
}

func (b *BugsService) Create(description BugDescription) (Bug, error) {
	values := defaultQuery.values()
	b.client.prepare(&values)

	var bug Bug
	buffer := bytes.Buffer{}
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(&description)
	if err != nil {
		return bug, err
	}

	url := fmt.Sprintf("%s%s?%s", b.client.Url(), bugEndpoint, values.Encode())
	req, err := http.NewRequest(http.MethodPost, url, &buffer)
	req.Header.Add("Content-type", "application/json; charset=utf-8")

	resp, err := b.client.do(req)
	if err != nil {
		return bug, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&bug)
	return bug, err
}

func (m Teams) MarshalJSON() ([]byte, error) {
	if len(m) == 0 {
		return emptyTeams, nil
	}

	buff := bytes.Buffer{}
	buff.WriteByte('[')
	for i, teamId := range m {
		if i != 0 {
			buff.WriteByte(',')
		}
		buff.WriteString(`{"Team":{"Id":`)
		buff.WriteString(strconv.Itoa(teamId))
		buff.WriteString(`}}`)
	}

	buff.WriteByte(']')
	return buff.Bytes(), nil
}
