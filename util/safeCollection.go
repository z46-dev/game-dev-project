package util

import "sync"

type Identifiable interface {
	GetID() uint64
}

type SafeStorage[T Identifiable] struct {
	storage           map[uint64]T
	mu                sync.RWMutex
	enqueuedAdditions []T
	enqueuedRemovals  []uint64
}

func NewSafeStorage[T Identifiable]() *SafeStorage[T] {
	return &SafeStorage[T]{
		storage:           make(map[uint64]T),
		enqueuedAdditions: make([]T, 0),
		enqueuedRemovals:  make([]uint64, 0),
	}
}

// Adds an item to the storage. In this instance, it is enqueued for addition until next flush.
func (ss *SafeStorage[T]) Add(item T) {
	ss.enqueuedAdditions = append(ss.enqueuedAdditions, item)
}

// Removes an item from the storage. In this instance, it is enqueued for removal until next flush.
func (ss *SafeStorage[T]) Remove(item T) {
	ss.enqueuedRemovals = append(ss.enqueuedRemovals, item.GetID())
}

// Flushes all enqueued additions and removals to the storage.
func (ss *SafeStorage[T]) Flush() {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	// Add enqueued items
	for _, item := range ss.enqueuedAdditions {
		ss.storage[item.GetID()] = item
	}

	ss.enqueuedAdditions = ss.enqueuedAdditions[:0]

	// Remove enqueued items
	for _, id := range ss.enqueuedRemovals {
		delete(ss.storage, id)
	}

	ss.enqueuedRemovals = ss.enqueuedRemovals[:0]
}

// Retrieves an item from the storage by its ID. Returns nil if the item does not exist.
func (ss *SafeStorage[T]) Get(id uint64) T {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	return ss.storage[id]
}

// Iterates over all items in the storage and applies the provided function to each item.
func (ss *SafeStorage[T]) ForEach(f func(T)) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	for _, item := range ss.storage {
		f(item)
	}
}

// Returns the number of items in the storage.
func (ss *SafeStorage[T]) Size() int {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	return len(ss.storage)
}

// Clears all items from the storage.
func (ss *SafeStorage[T]) Clear() {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	for key := range ss.storage {
		delete(ss.storage, key)
	}
}
