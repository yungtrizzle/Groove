package app

import (
	"sync"
)

//imported from stackoverflow
//thanks to tux21b

type Task interface {
	Execute()
}

type WorkPool struct {
	mu    sync.Mutex
	size  int
	tasks chan Task
	kill  chan struct{}
	wg    sync.WaitGroup
}

func NewPool(size int) *WorkPool {
	pool := &WorkPool{
		tasks: make(chan Task, 128),
		kill:  make(chan struct{}),
	}
	pool.Resize(size)
	return pool
}

func (p *WorkPool) worker() {
	defer p.wg.Done()
	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			task.Execute()
		case <-p.kill:
			return
		}
	}
}

func (p *WorkPool) Resize(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for p.size < n {
		p.size++
		p.wg.Add(1)
		go p.worker()
	}
	for p.size > n {
		p.size--
		p.kill <- struct{}{}
	}
}

func (p *WorkPool) Close() {
	close(p.tasks)
}

func (p *WorkPool) Wait() {
	p.wg.Wait()
}

func (p *WorkPool) Exec(task Task) {
	p.tasks <- task
}
