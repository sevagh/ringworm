package ringbuffer

import "fmt"

type ringBuffer struct {
        storage *[]interface{}
        max     int
        head    uint8
        tail    uint8
}

func NewRingBuffer() *ringBuffer {
        max := int(^uint8(0))
        storage := make([]int, max+1)

        return &ringBuffer{
                storage: &storage,
                max: max,
                head: 0,
                tail: 0,
        }
}

func (r *ringBuffer) InsertWithError(val int) error {
        if r.Full() {
                return fmt.Errorf("Ringbuffer is full!")
        }
        (*r.storage)[r.head] = val
        r.head += 1
        return nil
}

func (r *ringBuffer) Insert(val int) {
        (*r.storage)[r.head] = val
        r.head += 1
}

func (r *ringBuffer) Pop() (int, error) {
        if r.Empty() {
                return -1, fmt.Errorf("Ringbuffer is empty!")
        }
        ret := (*r.storage)[r.tail]
        r.tail += 1
        return ret, nil
}

func (r *ringBuffer) Peek() int {
        ret := (*r.storage)[r.tail]
        return ret
}

func (r *ringBuffer) Max() int {
        return r.max
}

func (r *ringBuffer) Empty() bool {
        return r.head == r.tail
}

func (r *ringBuffer) Size() int {
        return int(r.head - r.tail)
}

func (r *ringBuffer) Full() bool {
        return r.Size() == r.max
}
