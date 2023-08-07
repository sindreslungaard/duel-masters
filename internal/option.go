package internal

type Option[T any] struct {
	value *T
}

func NewOption[T any]() *Option[T] {
	return &Option[T]{
		value: nil,
	}
}

func (o *Option[T]) Unwrap() (value *T, ok bool) {
	if o.value == nil {
		return nil, false
	}

	return o.value, true
}

func (o *Option[T]) Set(value *T) *Option[T] {
	o.value = value
	return o
}

func (o *Option[T]) Clear() {
	o.value = nil
}

func (o *Option[T]) Some() bool {
	return o.value != nil
}
