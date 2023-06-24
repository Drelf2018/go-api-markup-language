package parser

type Handler[T any] struct {
	keys []func(T) bool
	pres []func(T)
	fs   []func(T)
}

func (hd *Handler[T]) Prepare(v func(T)) *Handler[T] {
	hd.pres = append(hd.pres, v)
	return hd
}

func (hd *Handler[T]) Add(k func(T) bool, v func(T)) *Handler[T] {
	hd.keys = append(hd.keys, k)
	hd.fs = append(hd.fs, v)
	return hd
}

func (hd *Handler[T]) Do(t T) {
	for _, pre := range hd.pres {
		pre(t)
	}
	for i, k := range hd.keys {
		if k(t) {
			hd.fs[i](t)
			break
		}
	}
}

func NewHandler[T any]() *Handler[T] {
	return &Handler[T]{
		make([]func(T) bool, 0),
		make([]func(T), 0),
		make([]func(T), 0),
	}
}
