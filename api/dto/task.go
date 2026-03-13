package dto

import "errors"

type TDTO struct {
	title string
	text  string
}

func (t *TDTO) Validate() error {
	if t.title == "" {
		return errors.New("title is empty")
	}
	if t.text == "" {
		return errors.New("text is empty")
	}
	return nil
}
