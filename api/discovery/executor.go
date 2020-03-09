package discovery

import (
	"context"
	"sync"
)

// executor provides shared logic to TargetExecutor and ClientExecutor.
type executor struct {
	m     sync.Mutex
	tasks map[interface{}]task
}

type task struct {
	Cancel context.CancelFunc
	Done   chan struct{}
}

func (e *executor) start(
	parent context.Context,
	key interface{},
	fn func(context.Context),
) {
	e.m.Lock()
	defer e.m.Unlock()

	if e.tasks == nil {
		e.tasks = map[interface{}]task{}
	} else if _, ok := e.tasks[key]; ok {
		return
	}

	if parent == nil {
		parent = context.Background()
	}

	ctx, cancel := context.WithCancel(parent)
	done := make(chan struct{})

	e.tasks[key] = task{
		cancel,
		done,
	}

	go func() {
		defer close(done)
		fn(ctx)
	}()
}

func (e *executor) stop(key interface{}) {
	if task, ok := e.remove(key); ok {
		task.Cancel()
		<-task.Done
	}
}

func (e *executor) remove(key interface{}) (task, bool) {
	e.m.Lock()
	defer e.m.Unlock()

	task, ok := e.tasks[key]
	if ok {
		delete(e.tasks, key)
	}

	return task, ok
}
