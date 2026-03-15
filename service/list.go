package service

import "sync"

type List struct {
	tasks map[string]Task
	mtx   sync.RWMutex
}

func NewList() *List {
	return &List{
		tasks: make(map[string]Task),
		mtx:   sync.RWMutex{},
	}
}

func (l *List) Add(task Task) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if _, ok := l.tasks[task.Title]; ok {
		return ErrTaskAlreadyExist
	}

	l.tasks[task.Title] = task

	return nil
}

func (l *List) Delete(title string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if _, ok := l.tasks[title]; !ok {
		return ErrTaskNotFound
	}

	delete(l.tasks, title)

	return nil
}

func (l *List) Complete(title string) (Task, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	task, ok := l.tasks[title]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	task.Complete()
	l.tasks[title] = task

	return task, nil
}

func (l *List) Uncomplete(title string) (Task, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	task, ok := l.tasks[title]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	task.Uncomplete()
	l.tasks[title] = task

	return task, nil
}

func (l *List) GetOne(title string) (Task, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	task, ok := l.tasks[title]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	return task, nil
}

func (l *List) GetAll() map[string]Task {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	return l.tasks
}
