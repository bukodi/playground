package generics

// https://www.dolthub.com/blog/2024-07-01-golang-generic-collections/

import "sort"

// Interface type used for vars
type Sortable[T comparable] interface {
	Less(member T) bool
}

// Type set used only for constraints, not vars
type SortableConstraint[T comparable] interface {
	comparable
	Sortable[T]
}

type SortableSet[T SortableConstraint[T]] interface {
	Add(member T)
	Size() int
	Contains(member T) bool
	Sorted() []T
}

type MapSet[T SortableConstraint[T]] struct {
	members map[T]struct{}
}

func NewMapSet[T SortableConstraint[T]]() SortableSet[T] {
	return MapSet[T]{
		members: make(map[T]struct{}),
	}
}

func (s MapSet[T]) Add(member T) {
	s.members[member] = struct{}{}
}

func (s MapSet[T]) Size() int {
	return len(s.members)
}

func (s MapSet[T]) Contains(member T) bool {
	_, found := s.members[member]
	return found
}

func (s MapSet[T]) Sorted() []T {
	sorted := make([]T, 0, len(s.members))
	for member := range s.members {
		sorted = append(sorted, member)
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Less(sorted[j])
	})

	return sorted
}

type SliceSet[T SortableConstraint[T]] struct {
	members []T
}

func NewSliceSet[T interface {
	Sortable[T]
	comparable
}]() SortableSet[T] {
	return &SliceSet[T]{
		members: make([]T, 0),
	}
}

func (s *SliceSet[T]) Add(member T) {
	if !s.Contains(member) {
		s.members = append(s.members, member)
	}
}

func (s SliceSet[T]) Size() int {
	return len(s.members)
}

func (s SliceSet[T]) Contains(member T) bool {
	for _, m := range s.members {
		if m == member {
			return true
		}
	}
	return false
}

func (s SliceSet[T]) Sorted() []T {
	sorted := make([]T, len(s.members))
	copy(sorted, s.members)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Less(sorted[j])
	})

	return sorted
}
