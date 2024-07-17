package glycogen

type Option[T any] struct {
	none bool
	val  T
}

func (o Option[T]) OnSome(f func(T)) {
	if !o.none {
		f(o.val)
	}
}

func (o Option[T]) OnNone(f func()) {
	if o.none {
		f()
	}
}

func (o Option[T]) OnBoth(some func(T), none func()) {
	if o.none {
		none()
	} else {
		some(o.val)
	}
}

func (o Option[T]) UnwrapOrPanic() T {
	if o.none {
		panic("called Unwrap on a None option")
	}
	return o.val
}

func (o Option[T]) UnwrapOrDefault(def T) T {
	if o.none {
		return def
	}
	return o.val
}

func (o Option[T]) UnwrapOrReturn(ctx *ContextExecution[Option[T]]) T {
	if o.none {
		ctx.Return(o)
	}
	return o.val
}

func (o Option[T]) UnwrapOrHandle(f func()) T {
	if o.none {
		f()
	}
	return o.val
}

func (o Option[T]) UnwrapBoth() (T, bool) {
	return o.val, !o.none
}

func Some[T any](v T) Option[T] {
	return Option[T]{val: v}
}

func None[T any]() Option[T] {
	return Option[T]{none: true}
}

func MapAccess[K comparable, V any](m map[K]V, k K) Option[V] {
	v, ok := m[k]
	if !ok {
		return None[V]()
	}
	return Some(v)
}

func FuncOk0[R any](f func() (R, bool)) Option[R] {
	v, ok := f()
	if !ok {
		return None[R]()
	}
	return Some(v)
}

func FuncOk1[P, R any](f func(P) (R, bool), p P) Option[R] {
	v, ok := f(p)
	if !ok {
		return None[R]()
	}
	return Some(v)
}

func FuncOk2[P1, P2, R any](f func(P1, P2) (R, bool), p1 P1, p2 P2) Option[R] {
	v, ok := f(p1, p2)
	if !ok {
		return None[R]()
	}
	return Some(v)
}

func FuncOk3[P1, P2, P3, R any](f func(P1, P2, P3) (R, bool), p1 P1, p2 P2, p3 P3) Option[R] {
	v, ok := f(p1, p2, p3)
	if !ok {
		return None[R]()
	}
	return Some(v)
}
