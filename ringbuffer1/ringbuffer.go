/*
Package ringbuffer implements a simple circular buffer.

The indexing strategy is taken from the conversations here:
https://www.snellman.net/blog/archive/2016-12-13-ring-buffers/

The read and write pointers are uint32s that are stored modulo 2*capacity,
and are moduloed with capacity to index the underlying []byte storage.

Here are some of the characteristics:
- Size must be power of two (for efficient modulo)
- SPSC (single producer single consumer)
- Lock-free using sync.atomic
- Fixed size, no growing

It operates on []byte, which could make it usable for a variety of different
applications by using encoding/gob or similar.
*/
package ringbuffer

import (
	"fmt"
	"sync/atomic"
)

// A Ringbuffer is a struct that allows users to store and read []byte data.
type Ringbuffer struct {
	read  uint32
	write uint32
	buf   []byte
}

func (r *Ringbuffer) mask(ptr uint32) uint32 {
	return ptr & (uint32(len(r.buf)) - 1)
}

func (r *Ringbuffer) mask2(ptr uint32) uint32 {
	return ptr & (2*uint32(len(r.buf)) - 1)
}

func (r *Ringbuffer) writePtr() uint32 {
	return atomic.LoadUint32(&r.write)
}

func (r *Ringbuffer) readPtr() uint32 {
	return atomic.LoadUint32(&r.read)
}

// Size returns the size (bytes written by the user) of the ringbuffer.
// This is the distance between the write and read pointers.
func (r *Ringbuffer) Size() int {
	return int(r.mask2(r.writePtr() - r.readPtr()))
}

// Empty returns true if the ringbuffer is empty, false otherwise.
func (r *Ringbuffer) Empty() bool {
	return r.readPtr() == r.writePtr()
}

// Full returns true if the ringbuffer is full, false otherwise.
func (r *Ringbuffer) Full() bool {
	return r.Size() == r.Capacity()
}

// Capacity returns the capacity of the underlfying []byte buf.
// This is the same capacity that the user initialized the ringbuffer with.
func (r *Ringbuffer) Capacity() int {
	return len(r.buf)
}

// NewRingbuffer creates a ringbuffer with the specified capacity.
// Note that the capacity must be a power of two or an error is returned.
func NewRingbuffer(capacity int) (Ringbuffer, error) {
	if (capacity == 0) || ((capacity & (capacity - 1)) != 0) {
		return Ringbuffer{}, fmt.Errorf("please use a power-of-two size")
	}
	buf := make([]byte, capacity)
	return Ringbuffer{
		read:  0,
		write: 0,
		buf:   buf,
	}, nil
}

// Write copies all the bytes in the provided []byte slice into the ringbuffer.
// Data is copied to storage[write:], and the write pointer is advanced by n
// bytes written.
//
// If there isn't enough space for the entire write, it returns an Error.
func (r *Ringbuffer) Write(buf []byte) error {
	emptyCount := r.Capacity() - r.Size()

	if len(buf) > emptyCount {
		return fmt.Errorf("write %d is too big for remaining capacity %d", len(buf), emptyCount)
	}
	desiredWrite := uint32(len(buf))

	capacity := uint32(len(r.buf))
	writeIdx := r.mask(r.writePtr())

	copy(r.buf[writeIdx:], buf)

	if writeIdx+desiredWrite > capacity {
		// wraparound
		remain := capacity - writeIdx
		copy(r.buf, buf[remain:])
	}

	atomic.AddUint32(&r.write, desiredWrite)
	atomic.SwapUint32(&r.write, r.mask2(r.writePtr()))

	return nil
}

// Read fills the provided []byte slice with as much data as can fit. Data is
// copied from  the ringbuffer's storage[read:], and the read pointer is
// advanced by n bytes read.
//
// There are no errors - care must be taken for partial reads if the ringbuffer
// doesn't have enough data. Best check the ringbuffer Size() method beforehand.
func (r *Ringbuffer) Read(buf []byte) {
	if r.Empty() {
		return
	}

	size := r.Size()
	readCountTmp := len(buf)
	if size < readCountTmp {
		readCountTmp = r.Size()
	}
	readCount := uint32(readCountTmp)

	capacity := uint32(len(r.buf))
	readIdx := r.mask(r.readPtr())

	var remain uint32 = 0
	var firstChunk uint32 = 0
	possibleFirstChunk := capacity - readIdx

	if readCount > possibleFirstChunk {
		remain = readCount - possibleFirstChunk
		firstChunk = possibleFirstChunk
	} else {
		firstChunk = readCount
	}

	copy(buf, r.buf[readIdx:readIdx+firstChunk])
	copy(buf[firstChunk:], r.buf[:remain])

	atomic.AddUint32(&r.read, uint32(readCount))
	atomic.SwapUint32(&r.read, r.mask2(r.readPtr()))
}

// Drain creates and returns a []byte slice containing all data in the
// ringbuffer. This empties the ringbuffer.
func (r *Ringbuffer) Drain() []byte {
	buf := make([]byte, r.Size())
	r.Read(buf)
	return buf
}
