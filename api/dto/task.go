package dto

import "errors"

type TDTO struct {
	Title string
	Text  string
}

func (t *TDTO) Validate() error {
	if t.Title == "" {
		return errors.New("title is empty")
	}
	if t.Text == "" {
		return errors.New("text is empty")
	}
	return nil
}
