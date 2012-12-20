package tagheap_test

import (
	"fmt"
	"testing"

	"github.com/soniakeys/tagheap"
)

func ExampleNew() {
	type s struct {
		N int `heap:"min"`
		X int `heap:"index"`
	}
	var s0 []*s
	tagheap.New(`heap`, &s0)
}

func ExampleTagHeap_Remove() {
	type s struct {
		N int `heap:"min"`
		X int `heap:"index"`
	}
	var s0 []*s
	h, err := tagheap.New(`heap`, &s0)
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

func TestSingle(t *testing.T) {
	// test some settings to contrast with the example above.
	type s struct {
		vert string // an extra field, no index field
		Tent int    `pq:"max"` // alternate key, max heap instead of min
	}

	// and a non-empty initial slice.
	s0 := []*s{
		&s{"a", 3},
		&s{"b", 1},
	}
	th, err := tagheap.New(`pq`, &s0)
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

func TestMultiple(t *testing.T) {
	// combination of example and test single
	type s struct {
		Tent int `pq:"max"`
		N    int `heap:"min"`
		X    int `heap:"index" pq:"index"`
	}

	// initialize heap1
	var s1 []*s
	h1, err := tagheap.New(`heap`, &s1)
	if err != nil {
		t.Fatal(err)
	}

	// initialize heap2
	s2 := []*s{
		&s{Tent: 3},
		&s{Tent: 1},
	}
	h2, err := tagheap.New(`pq`, &s2)
	if err != nil {
		t.Fatal(err)
	}

	// some heap1 operations
	h1.Push(&s{N: 3})
	h1.Push(&s{N: 1})
	four := &s{N: 4}
	h1.Push(four)

	// a heap2 operation
	h2.Push(&s{Tent: 4})

	// heap1 ops
	h1.Push(&s{N: 1})
	h1.Push(&s{N: 5})

	// heap2 ops.  max heap should count down 4, 3, 1
	t2 := h2.Pop().(*s)
	if t2.Tent != 4 {
		t.Fatalf("%#v, expected 4", t2)
	}
	t2 = h2.Pop().(*s)
	if t2.Tent != 3 {
		t.Fatalf("%#v, expected 3", t2)
	}

	// heap1 ops.  min heap should count up 1, 1, 3, 5
	t1 := h1.Pop().(*s)
	if t1.N != 1 {
		t.Fatalf("%#v, expected 1", t1)
	}
	h1.Remove(four)
	t1 = h1.Pop().(*s)
	if t1.N != 1 {
		t.Fatalf("%#v, expected 1", t1)
	}
	t1 = h1.Pop().(*s)
	if t1.N != 3 {
		t.Fatalf("%#v, expected 3", t1)
	}
	t1 = h1.Pop().(*s)
	if t1.N != 5 {
		t.Fatalf("%#v, expected 5", t1)
	}

	// one last heap2 op
	t2 = h2.Pop().(*s)
	if t2.Tent != 1 {
		t.Fatalf("%#v, expected 1", t2)
	}
}
