package aoi

import (
	"reflect"
	"sort"
)

type Int32Slice []uint32

func (p Int32Slice) Len() int           { return len(p) }
func (p Int32Slice) Less(i, j int) bool { return int64(p[i]) < int64(p[j]) }
func (p Int32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p Int32Slice) Sort() { sort.Sort(p) }

type Set interface {
	Add(uint32)
	Remove(uint32)
	Contains(uint32) bool
	Equals(Set) bool
	Length() int
	Values() []uint32
	Copy() Set
	Sum(Set) Set
	Sub(Set) Set
}

func NewUnsafeSet(values ...uint32) *unsafeSet {
	set := &unsafeSet{make(map[uint32]struct{})}
	for _, v := range values {
		set.Add(v)
	}
	return set
}

type unsafeSet struct {
	d map[uint32]struct{}
}

// Add adds a new value to the set (no-op if the value is already present)
func (us *unsafeSet) Add(value uint32) {
	us.d[value] = struct{}{}
}

// Remove removes the given value from the set
func (us *unsafeSet) Remove(value uint32) {
	delete(us.d, value)
}

// Contains returns whether the set contains the given value
func (us *unsafeSet) Contains(value uint32) (exists bool) {
	_, exists = us.d[value]
	return
}

// ContainsAll returns whether the set contains all given values
func (us *unsafeSet) ContainsAll(values []uint32) bool {
	for _, s := range values {
		if !us.Contains(s) {
			return false
		}
	}
	return true
}

// Equals returns whether the contents of two sets are identical
func (us *unsafeSet) Equals(other Set) bool {
	v1 := Int32Slice(us.Values())
	v2 := Int32Slice(other.Values())
	v1.Sort()
	v2.Sort()
	return reflect.DeepEqual(v1, v2)
}

// Length returns the number of elements in the set
func (us *unsafeSet) Length() int {
	return len(us.d)
}

// Values returns the values of the Set in an unspecified order.
func (us *unsafeSet) Values() (values []uint32) {
	values = make([]uint32, 0)
	for val := range us.d {
		values = append(values, val)
	}
	return
}

// Copy creates a new Set containing the values of the first
func (us *unsafeSet) Copy() Set {
	cp := NewUnsafeSet()
	for val := range us.d {
		cp.Add(val)
	}

	return cp
}

// Sub removes all elements in other from the set
func (us *unsafeSet) Sub(other Set) Set {
	oValues := other.Values()
	result := us.Copy().(*unsafeSet)

	for _, val := range oValues {
		if _, ok := result.d[val]; !ok {
			continue
		}
		delete(result.d, val)
	}

	return result
}

// 并集
func (self *unsafeSet) Sum(other Set) Set {
	result := self.Copy().(*unsafeSet)

	arr := other.Values()
	for _, val := range arr {
		result.Add(val)
	}

	return result
}
