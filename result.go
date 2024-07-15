package glycogen

import "errors"

type Result[T any] struct {
	err error
	val T
}

func (r *Result[T]) OnOk(f func(T)) {
	if r.err == nil {
		f(r.val)
	}
}

func (r *Result[T]) OnErr(f func(error)) {
	if r.err != nil {
		f(r.err)
	}
}

func (r *Result[T]) OnBoth(ok func(T), err func(error)) {
	if r.err != nil {
		err(r.err)
	} else {
		ok(r.val)
	}
}

func (r *Result[T]) Unwrap() T {
	if r.err != nil {
		panic(errors.Join(errors.New("called Unwrap on an Err result"), r.err))
	}
	return r.val
}

func (r *Result[T]) UnwrapErr() error {
	if r.err == nil {
		panic("called UnwrapErr on an Ok result")
	}
	return r.err
}

func (r *Result[T]) UnwrapOr(def T) T {
	if r.err != nil {
		return def
	}
	return r.val
}

func (r *Result[T]) UnwrapBoth() (T, error) {
	return r.val, r.err
}

func Ok[T any](v T) Result[T] {
	return Result[T]{val: v}
}

func Err[T any](e error) Result[T] {
	return Result[T]{err: e}
}
