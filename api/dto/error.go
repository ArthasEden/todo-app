package dto

import (
	"encoding/json"
	"log"
	"time"
)

type EDTO struct {
	CreatedAt time.Time
	Message   string
}

func eDTOnew(message string) *EDTO {
	return &EDTO{
		CreatedAt: time.Now(),
		Message:   message,
	}
}

func (e *EDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		log.Println("couldn't marshal EDTO:", err)
		return "{}"
	}

	return string(b)
}
