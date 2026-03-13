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

func ErrDTO(message error) string {
	eDTO := EDTO{
		CreatedAt: time.Now(),
		Message:   message.Error(),
	}

	return eDTO.ToString()
}

func (e *EDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		log.Println("couldn't marshal EDTO:", err)
		return "{}"
	}

	return string(b)
}
