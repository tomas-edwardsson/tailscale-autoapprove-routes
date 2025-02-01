package main

type Set[E comparable] map[E]struct{}

func NewSet[E comparable](vals ...E) Set[E] {
	s := make(Set[E])
	for _, val := range vals {
		s[val] = struct{}{}
	}
	return s
}

func (s Set[E]) Add(value E) {
	s[value] = struct{}{}
}

func (s Set[E]) Remove(value E) {
	delete(s, value)
}

func (s Set[E]) Contains(value E) bool {
	_, ok := s[value]
	return ok
}

func (s Set[E]) Union(other Set[E]) Set[E] {
	result := NewSet[E]()
	for value := range s {
		result.Add(value)
	}
	for value := range other {
		result.Add(value)
	}
	return result
}

func (s Set[E]) Intersection(other Set[E]) Set[E] {
	result := NewSet[E]()
	for value := range s {
		if other.Contains(value) {
			result.Add(value)
		}
	}
	return result
}

func (s Set[E]) Difference(other Set[E]) Set[E] {
	result := NewSet[E]()
	for value := range s {
		if !other.Contains(value) {
			result.Add(value)
		}
	}
	return result
}

func (s Set[E]) Len() int {
	return len(s)
}

func (s Set[E]) ToSlice() []E {
	vals := make([]E, 0, len(s))
	for val := range s {
		vals = append(vals, val)
	}
	return vals
}
