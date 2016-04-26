// Package buffer provides pool of custom byte buffers.
//
// It is just ByteBuffer extracted from github.com/valyala/fasthttp.
package buffer

import (
	"sync"
)

const (
	defaultByteBufferSize = 128
)

// Buffer provides byte buffer, which can be used
// in order to minimize memory allocations.
//
// ByteBuffer may be used with functions appending data to the given []byte
// slice. See example code for details.
//
// Use AcquireByteBuffer for obtaining an empty byte buffer.
type Buffer struct {

	// B is a byte buffer to use in append-like workloads.
	// See example code for details.
	B []byte
}

// Write implements io.Writer - it appends p to ByteBuffer.B
func (b *Buffer) Write(p []byte) (int, error) {
	return b.Append(p), nil
}

// Append appends p to ByteBuffer.B and returns length of p
func (b *Buffer) Append(p []byte) int {
	b.B = append(b.B, p...)
	return len(p)
}

// Reset makes ByteBuffer.B empty.
func (b *Buffer) Reset() {
	b.B = b.B[:0]
}

// Acquire returns an empty byte buffer from the pool.
//
// Acquired byte buffer may be returned to the pool via ReleaseByteBuffer call.
// This reduces the number of memory allocations required for byte buffer
// management.
func Acquire() *Buffer {
	v := pool.Get()
	if v == nil {
		return &Buffer{
			B: make([]byte, 0, defaultByteBufferSize),
		}
	}
	return v.(*Buffer)
}

// Release returns byte buffer to the pool.
//
// ByteBuffer.B mustn't be touched after returning it to the pool.
// Otherwise data races occur.
func Release(b *Buffer) {
	b.B = b.B[:0]
	pool.Put(b)
}

var (
	pool sync.Pool
)
