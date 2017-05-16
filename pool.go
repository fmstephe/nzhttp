package main

type workerPool struct {
	size int
	work chan func()
}

func newWorkerPool(size int) *workerPool {
	work := make(chan func(), size)
	for i := 0; i < size; i++ {
		go func() {
			for f := range work {
				f()
			}
		}()
	}
	return &workerPool{
		size: size,
		work: work,
	}
}

func (p *workerPool) submit(f func()) {
	p.work <- f
}

func (p *workerPool) stop() {
	close(p.work)
}
