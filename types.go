package targetprocess

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	jsonNull = []byte(`null`)

	// \/Date(1477663200000-0500)\/
	dateRx = regexp.MustCompile(`\\/Date\((\d+)[-+]\d+\)\\/`)
)

type DateTime time.Time

func (dt *DateTime) UnmarshalJSON(data []byte) error {
	// i don't know how it do
	if bytes.Compare(data, jsonNull) != 0 {
		found := dateRx.FindStringSubmatch(string(data))
		if len(found) == 0 {
			return fmt.Errorf("Failed to parse date: %s", string(data))
		}

		milliseconds, _ := strconv.ParseUint(found[1], 10, 64)
		var unixTime time.Time

		unixTime = time.Unix(int64(milliseconds/1000), 0)
		*dt = DateTime(unixTime)
	}
	return nil
}

type EntityState struct {
	Id       int64
	Type     string `json:"ResourceType"`
	Name     string
	Priority float32 `json:"NumericPriority"`
}

type EntityType struct {
	Id   int64
	Type string `json:"ResourceType"`
	Name string
}

type Project struct {
	Id      int64
	Type    string `json:"ResourceType"`
	Name    string
	Process Process
}

type Process struct {
	Id   int64
	Type string `json:"ResourceType"`
	Name string
}

type Priority struct {
	Id         int64
	Type       string `json:"ResourceType"`
	Name       string
	Importance int
}

type Bug struct {
	Id          int64
	Name        string
	Description string
	StartDate   DateTime `json:"StartDate,omitempty"`
	EndDate     DateTime `json:"EndDate,omitempty"`
	CreateDate  DateTime `json:"CreateDate,omitempty"`

	Tags     string  `json:"Tags,omitempty"`
	Priority float32 `json:"NumericPriority"`
	Effort   float32 `json:"Effort,omitempty"`
	Units    string
	Type     EntityType
	Project  Project
}

type UseStory struct {
	Id          int64
	Name        string
	Description string
	StartDate   DateTime `json:"StartDate,omitempty"`
	EndDate     DateTime `json:"EndDate,omitempty"`
	CreateDate  DateTime `json:"CreateDate,omitempty"`

	Tags     string  `json:"Tags,omitempty"`
	Priority float32 `json:"NumericPriority"`
	Effort   float32 `json:"Effort,omitempty"`
	Units    string
	Type     EntityType
	Project  Project
}

type Assignable struct {
	Id          int64
	Name        string
	Description string
}
