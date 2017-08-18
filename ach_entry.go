package rtn

import (
	"fmt"
	"log"
	"strings"
)

// Format: https://www.frbservices.org/EPaymentsDirectory/achFormat.html
type achEntry struct {
	City          string `json:"city"`
	State         string `json:"state"`
	RoutingNumber string `json:"routingNumber"`
	Name          string `json:"name"`
	LastModified  string `json:"lastModified"`
}

func makeFrom(line string) achEntry {
	if len(line) < 150 {
		log.Println("Unexpected line length.")
		return achEntry{}
	}

	lm := strings.TrimSpace(line[20:26])
	if len(lm) == 6 {
		lm = fmt.Sprintf("%v-%v-%v", lm[0:2], lm[2:4], lm[4:6])
	}

	return achEntry{
		RoutingNumber: line[0:9],
		Name:          strings.TrimSpace(line[35:71]),
		State:         strings.TrimSpace(line[127:129]),
		City:          strings.TrimSpace(line[107:127]),
		LastModified:  lm,
	}
}
