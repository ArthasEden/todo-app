package service

import "errors"

var (
	ErrTaskAlreadyExist = errors.New("task already exist")
	ErrTaskNotFound     = errors.New("task not found")
)
