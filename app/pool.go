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

	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			task.Execute()
		case <-p.kill:
			return
		
		default:
			p.mu.Lock()
			wt := len(p.tasks)
			p.mu.Unlock()

			if wt == 0 {
				p.size--
				return
			}
		}
	}
}

func (p *WorkPool) Resize(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for p.size < n {
		p.size++
		go p.worker()
	}
	for p.size > n {
		p.size--
		p.kill <- struct{}{}
	}
}

func (p *WorkPool) Close() {
	close(p.tasks)
	close(p.kill)
}

func (p *WorkPool) Wait() {
	for {
		p.mu.Lock()
		size := p.size
		p.mu.Unlock()

		if size == 0 {
			return
		}
	}
}

func (p *WorkPool) Exec(task Task) {
	p.tasks <- task
}
