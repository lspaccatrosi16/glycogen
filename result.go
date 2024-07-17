package glycogen

import "errors"

type Result[T any] struct {
	err error
	val T
}

func (r Result[T]) OnOk(f func(T)) {
	if r.err == nil {
		f(r.val)
	}
}

func (r Result[T]) OnErr(f func(error)) {
	if r.err != nil {
		f(r.err)
	}
}

func (r Result[T]) OnBoth(ok func(T), err func(error)) {
	if r.err != nil {
		err(r.err)
	} else {
		ok(r.val)
	}
}

func (r Result[T]) UnwrapOrPanic() T {
	if r.err != nil {
		panic(errors.Join(errors.New("called Unwrap on an Err result"), r.err))
	}
	return r.val
}

func (r Result[T]) UnwrapErrOrPanic() error {
	if r.err == nil {
		panic(errors.New("called UnwrapErr on an Ok result"))
	}
	return r.err
}

func (r Result[T]) UnwrapOrDefault(def T) T {
	if r.err != nil {
		return def
	}
	return r.val
}

func (r Result[T]) UnwrapOrHandle(f func(error)) T {
	if r.err != nil {
		f(r.err)
	}
	return r.val
}

func (r Result[T]) UnwrapBoth() (T, error) {
	return r.val, r.err
}

func (r Result[T]) UnwrapOrReturn(ctx *ContextExecution[Result[T]]) T {
	if r.err != nil {
		ctx.Return(r)
	}
	return r.val
}

func (r Result[T]) UnwrapOrReturnErr(ctx *ContextExecution[error], tag string) T {
	if r.err != nil {
		ctx.Return(ctx.WrapAndTag(r.err, tag))
	}
	return r.val
}

func Ok[T any](v T) Result[T] {
	return Result[T]{val: v}
}

func Err[T any](e error) Result[T] {
	return Result[T]{err: e}
}

func ErrCtx[T any](ctx *Context, msg string) Result[T] {
	return Err[T](ctx.Error(msg))
}

func ErrCtxf[T any](ctx *Context, format string, args ...interface{}) Result[T] {
	return Err[T](ctx.Errorf(format, args...))
}

func FuncErr0[R any](f func() (R, error)) Result[R] {
	v, err := f()
	if err != nil {
		return Err[R](err)
	}
	return Ok(v)
}

func FuncErr1[P, R any](f func(P) (R, error), p P) Result[R] {
	v, err := f(p)
	if err != nil {
		return Err[R](err)
	}
	return Ok(v)
}

func FuncErr2[P1, P2, R any](f func(P1, P2) (R, error), p1 P1, p2 P2) Result[R] {
	v, err := f(p1, p2)
	if err != nil {
		return Err[R](err)
	}
	return Ok(v)
}

func FuncErr3[P1, P2, P3, R any](f func(P1, P2, P3) (R, error), p1 P1, p2 P2, p3 P3) Result[R] {
	v, err := f(p1, p2, p3)
	if err != nil {
		return Err[R](err)
	}
	return Ok(v)
}

func CEHandler[T any](ctx *ContextExecution[Result[T]], tag string) func(error) {
	return func(err error) {
		ctx.Return(Err[T](ctx.WrapAndTag(err, tag)))
	}
}
