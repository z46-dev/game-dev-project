package game

func newSafeStorage[T GameObject]() (storage *SafeStorage[T]) {
	storage = &SafeStorage[T]{}
	storage.storage = make(map[uint64]T)
	return
}

func (s *SafeStorage[T]) Add(obj T) {
	s.enqueuedAdditions = append(s.enqueuedAdditions, obj)
}

func (s *SafeStorage[T]) Remove(id uint64) {
	s.enqueuedRemovals = append(s.enqueuedRemovals, id)
}

func (s *SafeStorage[T]) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, obj := range s.enqueuedAdditions {
		s.storage[obj.ID()] = obj
	}

	s.enqueuedAdditions = s.enqueuedAdditions[:0]

	for _, id := range s.enqueuedRemovals {
		delete(s.storage, id)
	}

	s.enqueuedRemovals = s.enqueuedRemovals[:0]
}

func (s *SafeStorage[T]) Get(id uint64) (obj T, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	obj, ok = s.storage[id]
	return
}

func (s *SafeStorage[T]) ForEach(f func(T)) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, obj := range s.storage {
		f(obj)
	}
}

func (s *SafeStorage[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage = make(map[uint64]T)
	s.enqueuedAdditions = s.enqueuedAdditions[:0]
	s.enqueuedRemovals = s.enqueuedRemovals[:0]
}
