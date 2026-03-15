package service

import "time"

type Task struct {
	Title     string
	Text      string
	Completed bool

	CreatedAt   time.Time
	CompletedAt *time.Time
}

func NewTask(title, text string) *Task {
	return &Task{
		Title:     title,
		Text:      text,
		Completed: false,

		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}
}

func (t *Task) Complete() {
	timeComplete := time.Now()

	t.CompletedAt = &timeComplete
	t.Completed = true
}

func (t *Task) Uncomplete() {
	t.CompletedAt = nil
	t.Completed = false
}
