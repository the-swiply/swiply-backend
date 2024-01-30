package workerpool

type WorkerPoolOption[T, U any] func(p *Pool[T, U])

func WithIgnoreResult[T, U any]() WorkerPoolOption[T, U] {
	return func(p *Pool[T, U]) {
		p.ignoreResult = true
	}
}
