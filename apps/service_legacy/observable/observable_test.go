package observable_test

import (
	"testing"

	"github.com/bensivo/salad-bowl/observable"
	"github.com/stretchr/testify/assert"
)

func TestSetGet_ReturnsValue(t *testing.T) {
	o := observable.New(map[string]interface{}{})

	assert.Equal(t, o.Get(), map[string]interface{}{})

	value := map[string]interface{}{
		"hello": "world",
	}
	o.Set(value)

	assert.Equal(t, o.Get(), value)
}

func TestAddObserver_CallsWithCurrentValue(t *testing.T) {
	// Given an observable
	o := observable.New(map[string]interface{}{})

	called := false
	onChange := func(value map[string]interface{}) {
		called = true
	}

	// When I add an observer
	o.OnChange(onChange)

	// It is called
	assert.Equal(t, called, true)
}

func TestAddObserver_CalledOnSet(t *testing.T) {
	called := false

	// Given an observable, with an observer
	o := observable.New(map[string]interface{}{})
	onChange := func(value map[string]interface{}) {
		called = true
	}
	o.OnChange(onChange)
	called = false

	// When I set a value
	o.Set(map[string]interface{}{
		"hello": "world",
	})

	// Then the observer is called
	assert.Equal(t, called, true)
}
