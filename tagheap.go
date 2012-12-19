// Tagheap provides an interface-free heap API.
//
// Internally, the implementation uses container/heap from the standard library
// but it wraps this in an API that uses struct tags to direct heap ordering.
// The constructor argument specifies a struct that represents heap nodes.
// This struct must be annotated with struct tags.
//
// Struct tags
//
// In the struct type, a field with the struct key `heap:"min"` is called
// a key field and specifies that the heap will be a min-heap based on
// the value of this field.  The field must be < comparable, that is,
// it must be a string, integer, or floating point type.
//
// The struct key `heap:"max"` similarly indicates a key field and
// specifies a max-heap.  There must be exactly one key field in the struct,
// either min or max.
//
// Optionally, another field may have the tag `heap:"index"`.  The tag
// specifies that tagHeap methods should maintain this field as an index
// that can be used by the Remove method.  The field type must be int.
package tagheap

import (
	"container/heap"
	"errors"
	"fmt"
	"reflect"
)

// TagHeap exports heap functions.
//
// The unexported underlying type implements container/heap.Interface.
type TagHeap tagHeap

// New constructs a new TagHeap object.
//
// For the argument, New takes a pointer to slice of pointer to struct.
// The struct type is the type you wish to store in the heap.
// It is valid to pass an empty slice, but it is not valid
// to pass an untyped nil.
//
// The heap is initialized to the contents of the slice and the slice
// is used without making a copy.  The slice and its contents are modified
// by various TagHeap methods.
//
// If the struct type is not properly specified or the argument is otherwise
// not usable for TagHeap, the error result explains why.
func New(arg interface{}) (*TagHeap, error) {
	t, err := newTagHeap(arg)
	if err != nil {
		return nil, err
	}
	return (*TagHeap)(t), nil
}

// Len returns the number of structs on the heap.
func (t TagHeap) Len() int { return tagHeap(t).Len() }

// Push performs a heap push operation, pushing a struct onto the heap.
func (t *TagHeap) Push(u interface{}) { heap.Push((*tagHeap)(t), u) }

// Pop performs a heap pop operation, popping the next struct in heap order
// from the heap.
func (t *TagHeap) Pop() interface{} { return heap.Pop((*tagHeap)(t)) }

// Remove performs a heap remove operation, removing the specified struct.
func (t *TagHeap) Remove(u interface{}) interface{} {
	th := (*tagHeap)(t)
	return heap.Remove(th,
		int(reflect.ValueOf(u).Elem().Field(th.indexFieldIndex).Int()))
}

// unexported type implementing heap.Interface
type tagHeap struct {
	s               reflect.Value // slice of ptr to struct
	minHeap         bool
	keyFieldIndex   int
	indexFieldIndex int
	less            func(vi, vj reflect.Value) bool
	swapTemp        reflect.Value
}

// constructor
func newTagHeap(arg interface{}) (*tagHeap, error) {
	at := reflect.TypeOf(arg)
	if at == nil {
		return nil, errors.New("argument cannot be untyped nil")
	}
	// create return value
	s := &tagHeap{
		keyFieldIndex:   -1,
		indexFieldIndex: -1,
	}
	if at.Kind() != reflect.Ptr {
		return nil, errors.New("argument must be pointer")
	}
	slct := at.Elem()
	if slct.Kind() != reflect.Slice {
		return nil, errors.New("argument must be pointer to slice")
	}
	pt := slct.Elem()
	if pt.Kind() != reflect.Ptr {
		return nil, errors.New("argument must be pointer to slice of pointer")
	}
	st := pt.Elem()
	if st.Kind() != reflect.Struct {
		return nil, errors.New("argument must be pointer to slice of pointer to struct")
	}
	// find and validate struct tags
	for i, n := 0, st.NumField(); i < n; i++ {
		sf := st.Field(i)
		switch tv := sf.Tag.Get("heap"); tv {
		case "":
			continue
		case "min", "max":
			if sf.PkgPath > "" {
				return nil, errors.New("key field must be exported")
			}
			if s.keyFieldIndex >= 0 {
				return nil, errors.New("struct tags specify multiple keys.")
			}
			switch k := sf.Type.Kind(); {
			case k == reflect.String:
				s.less = lessString
			case k >= reflect.Int && k <= reflect.Int64:
				s.less = lessInt
			case k >= reflect.Uint && k <= reflect.Uint64:
				s.less = lessUint
			case k == reflect.Float64 || k == reflect.Float32:
				s.less = lessFloat
			default:
				return nil, errors.New("key field must be " +
					"a string, integer, or floating point type")
			}
			s.keyFieldIndex = i
			if tv == "min" {
				s.minHeap = true
			}
		case "index":
			if sf.PkgPath > "" {
				return nil, errors.New("index field must be exported")
			}
			if s.indexFieldIndex >= 0 {
				return nil, errors.New("struct tags specify multiple indexes")
			}
			if sf.Type.Kind() != reflect.Int {
				return nil, errors.New("index field must have type int")
			}
			s.indexFieldIndex = i
		default:
			return nil, fmt.Errorf("invalid struct tag %q", tv)
		}
	}
	if s.keyFieldIndex < 0 {
		return nil, errors.New("struct must indicate key field")
	}
	// initialize s.s, swapTemp
	s.s = reflect.ValueOf(arg).Elem()
	s.swapTemp = reflect.New(pt).Elem()
	heap.Init(s)
	return s, nil
}

func lessString(vi, vj reflect.Value) bool { return vi.String() < vj.String() }
func lessInt(vi, vj reflect.Value) bool    { return vi.Int() < vj.Int() }
func lessUint(vi, vj reflect.Value) bool   { return vi.Uint() < vj.Uint() }
func lessFloat(vi, vj reflect.Value) bool  { return vi.Float() < vj.Float() }

// method of heap.Interface
func (s tagHeap) Len() int { return s.s.Len() }

// method of heap.Interface
func (s tagHeap) Less(i, j int) bool {
	return s.less(
		s.s.Index(i).Elem().Field(s.keyFieldIndex),
		s.s.Index(j).Elem().Field(s.keyFieldIndex)) == s.minHeap
}

// method of heap.Interface
func (s tagHeap) Swap(i, j int) {
	s.swapTemp.Set(s.s.Index(i))
	s.s.Index(i).Set(s.s.Index(j))
	s.s.Index(j).Set(s.swapTemp)
	if s.indexFieldIndex >= 0 {
		s.s.Index(i).Elem().Field(s.indexFieldIndex).SetInt(int64(i))
		s.s.Index(j).Elem().Field(s.indexFieldIndex).SetInt(int64(j))
	}
}

// method of heap.Interface
func (s *tagHeap) Push(u interface{}) {
	np := reflect.ValueOf(u)
	if s.indexFieldIndex >= 0 {
		np.Elem().Field(s.indexFieldIndex).SetInt(int64(s.s.Len()))
	}
	s.s.Set(reflect.Append(s.s, np))
}

// method of heap.Interface
func (s *tagHeap) Pop() interface{} {
	l := s.s.Len()
	if l == 0 {
		return nil
	}
	l--
	r := s.s.Index(l).Interface()
	s.s.Set(s.s.Slice(0, l))
	return r
}
