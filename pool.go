// Go-routine safe pool/cache of byte slices
package cache

type Pool struct {
	c chan []byte
	s int
}

// Create a new pool/cache
func New(slicesize int, nitems int) Pool {
	return Pool{make(chan []byte, nitems), slicesize}
}

// Initialize or reset a pool/cache. This is *not* atomic.
func (c *Pool) Init(slicesize, nitems int) {
	c.c = make(chan []byte, nitems)
	c.s = slicesize
}

// Put slice into pool/cache without zeroing it.
func (c *Pool) FreeNoZero(bs []byte) {
	select {
	case c.c <- bs:
	default:
	}
}

// Put slice into pool/cache zeroing it.
func (c *Pool) FreeZeroing(bs []byte) {
	for i := range bs {
		bs[i] = 0
	}
	c.FreeNoZero(bs)
}

// Allocate a new item from pool/cache or alloc it if the cache is empty.
func (c *Pool) Alloc() (bs []byte) {
	select {
	case bs = <-c.c:

	default:
		bs = make([]byte, c.s)
	}
	return
}
