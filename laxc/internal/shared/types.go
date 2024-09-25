package shared

type LocalSymVar int
type SymReg int

type Reg string

type Option[T any] struct {
	Value T
	IsSet bool
}

func Some[T any](value T) Option[T] {
	return Option[T]{
		Value: value,
		IsSet: true,
	}
}

func None[T any]() Option[T] {
	return Option[T]{}
}
