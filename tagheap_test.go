package tagheap_test

import (
	"fmt"
	"testing"

	"tagheap"
)

func ExampleNew() {
	type s struct {
		N int `heap:"min"`
		X int `heap:"index"`
	}
	var s0 []*s
	tagheap.New(&s0)
}

func ExampleTagHeap_Remove() {
	type s struct {
		N int `heap:"min"`
		X int `heap:"index"`
	}
	var s0 []*s
	h, err := tagheap.New(&s0)
	if err != nil {
		fmt.Println(err)
		return
	}

	h.Push(&s{N: 3})
	h.Push(&s{N: 1})
	four := &s{N: 4}
	h.Push(four)
	h.Push(&s{N: 1})
	h.Push(&s{N: 5})
	fmt.Println(h.Pop().(*s).N)
	h.Remove(four)
	fmt.Println(h.Pop().(*s).N)
	fmt.Println(h.Pop().(*s).N)
	fmt.Println(h.Pop().(*s).N)
	// Output:
	// 1
	// 1
	// 3
	// 5
}

func TestInit(t *testing.T) {
	// test some settings to contrast with the example above.
	type s struct {
		vert string // an extra field, no index field
		Tent int    `heap:"max"` // max heap instead of min
	}

	// and a non-empty initial slice.
	s0 := []*s{
		&s{"a", 3},
		&s{"b", 1},
	}
	th, err := tagheap.New(&s0)
	if err != nil {
		t.Fatal(err)
	}
	th.Push(&s{"c", 4})

	s1 := th.Pop().(*s)
	if s1.Tent != 4 {
		t.Fatalf("%#v, expected 4", s1)
	}
	s1 = th.Pop().(*s)
	if s1.Tent != 3 {
		t.Fatalf("%#v, expected 3", s1)
	}
	s1 = th.Pop().(*s)
	if s1.Tent != 1 {
		t.Fatalf("%#v, expected 1", s1)
	}
}
