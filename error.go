package targetprocess

import (
	"encoding/xml"
	"fmt"
)

type Error struct {
	XMLName xml.Name `xml:"Error"`
	Status  string
	Message string
	Type    string
	Details string
	Id      string `xml:"ErrorId"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s:%s", e.Status, e.Message)
}

var (
	Unauthorized = &Error{
		Status:  "Unauthorized",
		Message: "Check your credential",
	}
)
