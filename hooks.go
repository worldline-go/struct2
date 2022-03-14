package struct2

import (
	"errors"
	"reflect"
)

// ErrContinueHook usable with HookFunc.
var ErrContinueHook = errors.New("continue to decode")

// Hooker interface for structs.
type Hooker interface {
	Struct2Hook() interface{}
}

// HookFunc get reflect.Value to modify custom in decoder.
type HookFunc func(reflect.Value) (interface{}, error)
