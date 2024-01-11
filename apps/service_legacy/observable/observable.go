package observable

import "fmt"

type OnChangeHandler[T any] func(value T)

// Observable wraps a single map[string]interface with change listeners
type Observable[T any] struct {
	value T

	observer OnChangeHandler[T]
}

func New[T any](initialValue T) Observable[T] {
	return Observable[T]{
		value:    initialValue,
		observer: nil,
	}
}

func (o *Observable[T]) Get() T {
	return o.value
}

func (o *Observable[T]) Set(value T) {
	o.value = value
	o.Notify()
}

func (o *Observable[T]) Notify() {
	if o.observer == nil {
		return
	}

	fmt.Printf("Calling observer with value: %v\n", o.Get())
	o.observer(o.Get())
}

// OnChange registers a callback for changes.
//
// NOTE: the callback will be called immediately with the current value, similar to an rxjs BehaviorSubject
func (o *Observable[T]) OnChange(observer OnChangeHandler[T]) {
	o.observer = observer
	o.Notify()
}
