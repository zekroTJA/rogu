package rogu

import (
	"sync"
)

type resetObject interface {
	Reset()
}

type safePool[T resetObject] struct {
	p *sync.Pool
}

func newSafePool[T resetObject](create func() T) safePool[T] {
	return safePool[T]{
		p: &sync.Pool{
			New: func() any { return create() },
		},
	}
}

func (t *safePool[T]) Get() T {
	return t.p.Get().(T)
}

func (t *safePool[T]) Put(v T) {
	v.Reset()
	t.p.Put(v)
}
