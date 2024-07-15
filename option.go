package glycogen

type Option[T any] struct {
	none bool
	val  T
}

func (o *Option[T]) OnSome(f func(T)) {
	if !o.none {
		f(o.val)
	}
}

func (o *Option[T]) OnNone(f func()) {
	if o.none {
		f()
	}
}

func (o *Option[T]) OnBoth(some func(T), none func()) {
	if o.none {
		none()
	} else {
		some(o.val)
	}
}

func (o *Option[T]) Unwrap() T {
	if o.none {
		panic("called Unwrap on a None option")
	}
	return o.val
}

func (o *Option[T]) UnwrapOr(def T) T {
	if o.none {
		return def
	}
	return o.val
}

func Some[T any](v T) Option[T] {
	return Option[T]{val: v}
}

func None[T any]() Option[T] {
	return Option[T]{none: true}
}
