package generics

import (
	"log"
	"testing"
)

type Name struct {
	First string
	Last  string
}

var _ Sortable[Name] = Name{}

func (n Name) Less(member Name) bool {
	return n.First < member.First || n.First == member.First && n.Last < member.Last
}

func TestSortable(t *testing.T) {
	name1 := Name{"John", "Doe"}
	name2 := Name{"Jane", "Doe"}
	name3 := Name{"Frank", "Reynolds"}

	ss := NewSliceSet[Name]()
	ms := NewMapSet[Name]()

	sets := []SortableSet[Name]{ss, ms}
	for _, s := range sets {
		s.Add(name1)
		s.Add(name2)
		s.Add(name2)

		if s.Size() != 2 {
			log.Fatal("set size is not 2")
		}

		if !s.Contains(name1) {
			log.Fatal("set does not contain name1")
		}
		if !s.Contains(name2) {
			log.Fatal("set does not contain name1")
		}
		if s.Contains(name3) {
			log.Fatal("set contains name3")
		}

		sorted := ss.Sorted()
		expectedSorted := []Name{name2, name1}

		if len(sorted) != len(expectedSorted) {
			log.Fatal("sorted length does not match")
		}

		for i := range sorted {
			if sorted[i] != expectedSorted[i] {
				log.Fatal("sorted does not match")
			}
		}
	}
}
