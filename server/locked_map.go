package server

import (
	"sync"
	"time"
)

type LockedMap[Key comparable, T any] struct {
	allocator map[Key]*T
	lock      sync.Mutex
}

type LockedMapEntry[T any] struct {
	Value     *T
	UpdatedAt time.Time
}

func (mr *LockedMap[Key, T]) Size() (topResult int) {

	mr.lockMap(func() {
		topResult = len(mr.allocator)
	})

	return
}

func NewLockedMap[Key comparable, T any]() *LockedMap[Key, T] {

	reg := LockedMap[Key, T]{
		allocator: map[Key]*T{},
	}

	return &reg
}

func (mr *LockedMap[Key, T]) Put(key Key, val T) {

	mr.lockMap(func() {
		mr.allocator[key] = &val
	})

}

// put ref
func (mr *LockedMap[Key, T]) PutRef(key Key, val *T) {
	mr.lockMap(func() {
		mr.allocator[key] = val
	})
}

func (mr *LockedMap[Key, T]) Delete(key Key) {

	mr.lockMap(func() {
		delete(mr.allocator, key)
	})

}

func (mr *LockedMap[Key, T]) Copied() map[Key]T {

	copiedMap := map[Key]T{}

	mr.lockMap(func() {
		for key, val := range mr.allocator {
			copiedMap[key] = *val
		}
	})

	return copiedMap

}

func (mr *LockedMap[Key, T]) lockMap(cb func()) {
	mr.lock.Lock()
	defer mr.lock.Unlock()

	cb()
}

func (mr *LockedMap[Key, T]) Get(key Key) (topResult *T) {

	mr.lockMap(func() {
		oldVal, ok := mr.allocator[key]

		if ok {
			topResult = oldVal
		}
	})

	return
}

func (mr *LockedMap[Key, T]) CleanWithCb(cb func(it *T) bool) *LockedMap[Key, T] {

	newMap := map[Key]*T{}

	mr.lockMap(func() {

		for key, it := range mr.allocator {
			if cb(it) {
				newMap[key] = it
			}
		}

		mr.allocator = newMap

	})

	return mr
}

func (mr *LockedMap[Key, T]) GetOrCreate(key Key) (topResult *T) {

	mr.lockMap(func() {

		oldVal, ok := mr.allocator[key]

		if !ok {
			// create

			var obj T

			topResult = &obj
			mr.allocator[key] = &obj
		} else {
			topResult = oldVal
		}
	})

	return
}

func (mr *LockedMap[Key, T]) GetOrCreateWithFlag(key Key) (
	topResult *T,
	justCreated bool,
) {

	mr.lockMap(func() {

		oldVal, ok := mr.allocator[key]

		if !ok {
			// create

			var obj T

			topResult = &obj
			mr.allocator[key] = &obj
		} else {
			topResult = oldVal
		}

		justCreated = !ok
	})

	return
}
