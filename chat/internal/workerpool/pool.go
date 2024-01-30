package workerpool

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrIsStopping = errors.New("pool is stopping")
)

type argsWrapper[T any] struct {
	ctx  context.Context
	args T
}

type balancerFunc[T any] func(args T) int64

type Worker[T, U any] func(ctx context.Context, args T) U

type Pool[T, U any] struct {
	numOfWorkers int
	ins          []chan argsWrapper[T]
	out          chan U
	workerFunc   Worker[T, U]
	balancerFunc balancerFunc[T]
	wg           sync.WaitGroup
	mu           sync.RWMutex
	isStopping   bool

	ignoreResult bool
}

func NewPool[T, U any](numOfWorkers int, worker Worker[T, U], balancerFunc balancerFunc[T], opts ...WorkerPoolOption[T, U]) *Pool[T, U] {
	ins := make([]chan argsWrapper[T], numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		ins[i] = make(chan argsWrapper[T], numOfWorkers/2)
	}

	pool := &Pool[T, U]{
		numOfWorkers: numOfWorkers,
		workerFunc:   worker,
		balancerFunc: balancerFunc,
		ins:          ins,
		out:          make(chan U, numOfWorkers),
		wg:           sync.WaitGroup{},
		mu:           sync.RWMutex{},
		isStopping:   false,
		ignoreResult: false,
	}
	for _, opt := range opts {
		opt(pool)
	}

	return pool
}

func (p *Pool[T, U]) Start() {
	for i := 0; i < p.numOfWorkers; i++ {
		p.wg.Add(1)
		go p.handleTask(i)
	}
}

func (p *Pool[T, U]) handleTask(workerID int) {
	defer p.wg.Done()
	for taskArgs := range p.ins[workerID] {
		if p.ignoreResult {
			p.workerFunc(taskArgs.ctx, taskArgs.args)
		} else {
			p.out <- p.workerFunc(taskArgs.ctx, taskArgs.args)
		}
	}
}

func (p *Pool[T, U]) AddTask(ctx context.Context, args T) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.isStopping {
		return ErrIsStopping
	}

	p.ins[p.balancerFunc(args)%int64(p.numOfWorkers)] <- argsWrapper[T]{
		ctx:  ctx,
		args: args,
	}

	return nil
}

func (p *Pool[T, U]) Result() <-chan U {
	return p.out
}

func (p *Pool[T, U]) Stop(ctx context.Context) error {
	p.mu.Lock()
	p.isStopping = true
	p.mu.Unlock()

	for _, inCh := range p.ins {
		close(inCh)
	}
	stopCh := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(stopCh)
	}()

	select {
	case <-stopCh:
		close(p.out)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
