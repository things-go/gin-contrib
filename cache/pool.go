package cache

import (
	"encoding"
	"net/http"
	"sync"
)

var cachePool = &sync.Pool{
	New: func() any { return &BodyCache{Header: make(http.Header)} },
}

// Get implement Pool interface
func poolGet() *BodyCache {
	return cachePool.Get().(*BodyCache)
}

// Put implement Pool interface
func poolPut(c *BodyCache) {
	c.Data = c.Data[:0]
	c.Header = make(http.Header)
	c.encoding = nil
	cachePool.Put(c)
}

// BodyCache body cache store
type BodyCache struct {
	Status   int
	Header   http.Header
	Data     []byte
	encoding Encoding
}

var _ encoding.BinaryMarshaler = (*BodyCache)(nil)
var _ encoding.BinaryUnmarshaler = (*BodyCache)(nil)

func (b *BodyCache) MarshalBinary() ([]byte, error) {
	return b.encoding.Marshal(b)
}

func (b *BodyCache) UnmarshalBinary(data []byte) error {
	return b.encoding.Unmarshal(data, b)
}
