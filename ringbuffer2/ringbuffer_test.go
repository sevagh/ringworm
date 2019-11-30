package ringbuffer

import "testing"

func TestInsertions(t *testing.T) {
	rb := NewRingBuffer()
	t.Logf("Max capacity of ringbuffer: %d\n", rb.Max())

	for i := 0; i < 2*rb.Max(); i++ {
		err := rb.InsertWithError(i)
		if err != nil {
			t.Logf("At elem %d we have error: %s\n", i, err.Error())
		}
	}

	for i := 0; i < rb.Max(); i++ {
		pop, err := rb.Pop()
		if err != nil {
			t.Logf("Error: %s\n", err.Error())
		} else {
			t.Logf("pop: %d\n", pop)
		}
	}

	for i := 0; i < rb.Max(); i++ {
		pop, err := rb.Pop()
		if err != nil {
			t.Logf("Error: %s\n", err.Error())
		} else {
			t.Logf("pop: %d\n", pop)
		}
	}
}
