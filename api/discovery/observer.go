package discovery

import "sync"

// observerSet provides shared logic to the XXXObserverSet types.
type observerSet struct {
	m         sync.RWMutex
	observers map[interface{}]struct{}
	entities  map[interface{}]struct{}
}

// register adds an observer to the set.
func (s *observerSet) register(
	o interface{},
	fn func(e interface{}),
) {
	s.m.Lock()
	defer s.m.Unlock()

	if s.observers == nil {
		s.observers = map[interface{}]struct{}{}
	} else if _, ok := s.observers[o]; ok {
		return
	}

	s.observers[o] = struct{}{}
	s.notifyOne(fn)
}

// unregister removes an observer from the set.
func (s *observerSet) unregister(
	o interface{},
	fn func(e interface{}),
) {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.observers[o]; !ok {
		return
	}

	delete(s.observers, o)
	s.notifyOne(fn)
}

// add adds an entity to the set.
func (s *observerSet) add(
	e interface{},
	fn func(o interface{}),
) {
	s.m.Lock()
	defer s.m.Unlock()

	if s.entities == nil {
		s.entities = map[interface{}]struct{}{}
	} else if _, ok := s.entities[e]; ok {
		return
	}

	s.entities[e] = struct{}{}
	s.notifyAll(fn)
}

// remove removes an entity from the set.
func (s *observerSet) remove(
	e interface{},
	fn func(o interface{}),
) {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.entities[e]; !ok {
		return
	}

	delete(s.entities, e)
	s.notifyAll(fn)
}

// notifyAll notifies all observers about a change to one entity.
func (s *observerSet) notifyAll(
	fn func(o interface{}),
) {
	var g sync.WaitGroup

	g.Add(len(s.observers))

	for o := range s.observers {
		o := o // capture loop variable

		go func() {
			defer g.Done()
			fn(o)
		}()
	}

	g.Wait()
}

// notifyOne notifies one observer about a change to all entities.
func (s *observerSet) notifyOne(
	fn func(e interface{}),
) {
	var g sync.WaitGroup

	g.Add(len(s.entities))

	for e := range s.entities {
		e := e // capture loop variable

		go func() {
			defer g.Done()
			fn(e)
		}()
	}

	g.Wait()
}
