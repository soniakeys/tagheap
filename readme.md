Tagheap
=======
Interface-free heap API.

Internally, the implementation uses container/heap from the standard library
but it wraps this in an API that uses struct tags to direct heap ordering.
The constructor argument specifies a struct that represents heap nodes.
This struct must be annotated with struct tags.

To install
----------
go get github.com/soniakeys/tagheap
