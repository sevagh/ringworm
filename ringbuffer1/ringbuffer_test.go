package ringbuffer_test

import (
	"bytes"
	"sync"
	"testing"

	"github.com/google/gofuzz"
	"github.com/sevagh/ringworm/ringbuffer1"
)

func TestRingbufferNonPowerOfTwoSize(t *testing.T) {
	_, err := ringbuffer.NewRingbuffer(3)
	if err == nil {
		t.Errorf("Expected error when creating non-power-of-two capacity ringbuffer")
	}
}

func TestRingbufferEmptyReadDoesNothing(t *testing.T) {
	ringbuf, _ := ringbuffer.NewRingbuffer(4)

	emptyBytes := make([]byte, 4)
	readBuf := make([]byte, 4)
	ringbuf.Read(readBuf)

	if !bytes.Equal(readBuf, emptyBytes) {
		t.Errorf("Expected read to do nothing")
	}
}

func TestRingbufferReadTooMuchOnlyDrains(t *testing.T) {
	ringbuf, _ := ringbuffer.NewRingbuffer(4)

	writeBuf := []byte{0, 1, 2, 3}
	err := ringbuf.Write(writeBuf)
	if err != nil {
		t.Errorf("Didn't expect error when writing to ringbuf: %+v\n", err)
	}

	if !ringbuf.Full() {
		t.Errorf("Expected a write of 4 bytes to fill ringbuffer of 4 capacity\n")
	}

	readBuf := make([]byte, 16)

	ringbuf.Read(readBuf)

	for i, x := range readBuf {
		if i == 0 && x != 0 {
			t.Errorf("wrong value at %d, expected %d, got %d", i, x, writeBuf[i])
		} else if i == 1 && x != 1 {
			t.Errorf("wrong value at %d, expected %d, got %d", i, x, writeBuf[i])
		} else if i == 2 && x != 2 {
			t.Errorf("wrong value at %d, expected %d, got %d", i, x, writeBuf[i])
		} else if i == 3 && x != 3 {
			t.Errorf("wrong value at %d, expected %d, got %d", i, x, writeBuf[i])
		} else if i >= 4 && x != 0 {
			t.Errorf("expected empty at %d, got %d", i, readBuf[i])
		}
	}

	writeBuf = []byte{4, 5, 6}
	err = ringbuf.Write(writeBuf)
	if err != nil {
		t.Errorf("Didn't expect error when writing to ringbuf: %+v\n", err)
	}

	ringbuf.Read(readBuf[4:])

	for i, x := range readBuf[4:] {
		if i == 0 && x != 4 {
			t.Errorf("wrong value at %d, expected %d, got %d", i+4, x, writeBuf[i])
		} else if i == 1 && x != 5 {
			t.Errorf("wrong value at %d, expected %d, got %d", i+4, x, writeBuf[i])
		} else if i == 2 && x != 6 {
			t.Errorf("wrong value at %d, expected %d, got %d", i+4, x, writeBuf[i])
		} else if i >= 3 && x != 0 {
			t.Errorf("expected empty at %d, got %d", i+7, readBuf[i+4])
		}
	}

	writeBuf = []byte{7, 8, 9}
	err = ringbuf.Write(writeBuf)
	if err != nil {
		t.Errorf("Didn't expect error when writing to ringbuf: %+v\n", err)
	}

	ringbuf.Read(readBuf[7:])

	for i, x := range readBuf[7:] {
		if i == 0 && x != 7 {
			t.Errorf("wrong value at %d, expected %d, got %d", i+7, x, writeBuf[i])
		} else if i == 1 && x != 8 {
			t.Errorf("wrong value at %d, expected %d, got %d", i+7, x, writeBuf[i])
		} else if i == 2 && x != 9 {
			t.Errorf("wrong value at %d, expected %d, got %d", i+7, x, writeBuf[i])
		} else if i >= 3 && x != 0 {
			t.Errorf("expected empty at %d, got %d", i+7, readBuf[i+7])
		}
	}
}

func TestRingbufferFillCount(t *testing.T) {
	ringbuf, _ := ringbuffer.NewRingbuffer(128)
	if !ringbuf.Empty() {
		t.Errorf("Expected ringbuf to be empty")
	}
	if ringbuf.Capacity() != 128 {
		t.Errorf("Expected ringbuf to have 128 capacity")
	}
	if ringbuf.Size() != 0 {
		t.Errorf("Expected ringbuf to have 0 size")
	}

	testDat := []byte{'d', 'e', 'a', 'd', 'b', 'e', 'e', 'f'}
	ringbuf.Write(testDat)

	if ringbuf.Size() != len(testDat) {
		t.Errorf("Expected ringbuf to have %d size", len(testDat))
	}

	ringbuf.Read(testDat)

	if ringbuf.Size() != 0 {
		t.Errorf("Expected ringbuf to have 0 size")
	}
}

func TestRingbufferWrite(t *testing.T) {
	ringbuf, _ := ringbuffer.NewRingbuffer(32)

	testData := "hello, world!"
	dataBuf := []byte(testData)

	err := ringbuf.Write(dataBuf)
	if err != nil {
		t.Errorf("Didn't expect error when writing to ringbuf: %+v\n", err)
	}
	dataBuf = nil

	if len(dataBuf) != 0 {
		t.Errorf("Expected dataBuf to be empty after clearing it")
	}

	dataBuf = make([]byte, len(testData))
	ringbuf.Read(dataBuf)

	ret := string(dataBuf)
	if ret != testData {
		t.Errorf("Expected dataBuf to contain correct data after reading into it\n\texp: %+v\n\tgot: %+v\n", testData, ret)
	}
}

func TestRingbufferWriteTooMuch(t *testing.T) {
	ringbuf, _ := ringbuffer.NewRingbuffer(4)

	testData := "hello, world!"
	dataBuf := []byte(testData)

	err := ringbuf.Write(dataBuf)
	if err == nil {
		t.Errorf("Expected an error here")
	}
}

func TestRingbufferWriteMultiple(t *testing.T) {
	ringbuf, _ := ringbuffer.NewRingbuffer(32)

	testData := []string{
		"aaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbb",
		"cccccccccccccccc",
		"dddddddddddddddd",
	}

	var dataBuf []byte
	readBuf := make([]byte, 16)

	for i := 0; i < len(testData); i++ {
		dataBuf = []byte(testData[i])

		err := ringbuf.Write(dataBuf)
		if err != nil {
			t.Errorf("Didn't expect error when writing to ringbuf: %+v\n", err)
		}

		if ringbuf.Capacity() != 32 {
			t.Errorf("Ringbuffer shouldn't be growing here")
		}

		ringbuf.Read(readBuf)

		ret := string(readBuf)
		if ret != testData[i] {
			t.Errorf("Expected dataBuf to contain correct data after reading into it\n\texp: %+v\n\tgot: %+v\n", testData[i], ret)
		}
	}
}

func TestRingbufferWriteConcurrent(t *testing.T) {
	ringbuf, _ := ringbuffer.NewRingbuffer(512)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		dataBuf := []byte("aaaaaaaaaaaaaaaa")
		var err error

		for i := 0; i < 20; i++ {
			err = ringbuf.Write(dataBuf)
			if err != nil {
				t.Errorf("Didn't expect error when writing to ringbuf: %+v\n", err)
			}
		}

		wg.Done()
	}(&wg)

	go func(wg *sync.WaitGroup) {
		readBuf := make([]byte, 16)
		for i := 0; i < 20; i++ {
			ringbuf.Read(readBuf)
		}

		wg.Done()
	}(&wg)

	wg.Wait()
}

func TestRingbufferDrain(t *testing.T) {
	ringbuf, _ := ringbuffer.NewRingbuffer(32)

	testData := []string{
		"aaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbb",
		"cccccccccccccccc",
		"dddddddddddddddd",
	}

	for i := 0; i < len(testData); i++ {
		dataBuf := []byte(testData[i])

		err := ringbuf.Write(dataBuf)
		if err != nil {
			t.Errorf("Didn't expect error when writing to ringbuf: %+v\n", err)
		}

		if ringbuf.Capacity() != 32 {
			t.Errorf("Ringbuffer shouldn't be growing here")
		}

		ret := string(ringbuf.Drain())
		if ret != testData[i] {
			t.Errorf("Expected dataBuf to contain correct data after reading into it\n\texp: %+v\n\tgot: %+v\n", testData[i], ret)
		}
	}
}

func TestRingbufferWriteptrAdvances(t *testing.T) {
	ringbuf, _ := ringbuffer.NewRingbuffer(64)

	testData := []string{
		"aaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbb",
		"cccccccccccccccc",
		"dddddddddddddddd",
		"eeeeeeeeeeeeeeee",
		"ffffffffffffffff",
		"gggggggggggggggg",
		"hhhhhhhhhhhhhhhh",
	}

	expectError := false
	for i := 0; i < len(testData); i++ {
		dataBuf := []byte(testData[i])

		expectError = ((ringbuf.Capacity() - ringbuf.Size()) == 0)

		if i < 4 && expectError {
			t.Errorf("first 4 writes should've been ok. on write: %d", i)
		}

		if i >= 4 && !expectError {
			t.Errorf("last 4 writes should not be ok. on write: %d", i)
		}

		err := ringbuf.Write(dataBuf)
		if !expectError && err != nil {
			t.Errorf("Didn't expect error when writing size %d to ringbuf with size %d, capacity %d: %+v\n", err, len(dataBuf), ringbuf.Size(), ringbuf.Capacity())
		}
		if expectError && err == nil {
			t.Errorf("Expected error when writing size %d to ringbuf with size %d, capacity %d\n", len(dataBuf), ringbuf.Size(), ringbuf.Capacity())
		}
	}
}

func TestRingbufferFuzzBytes(t *testing.T) {
	ringbuf, _ := ringbuffer.NewRingbuffer(128)

	writeBuf := make([]byte, 64)

	f := fuzz.New().NumElements(1, 64)
	f.Fuzz(&writeBuf)

	err := ringbuf.Write(writeBuf)
	if err != nil {
		t.Errorf("error when writing to ringbuf: %+v\n", err)
	}

	bytesWritten := ringbuf.Size()

	readBuf := make([]byte, 64)
	ringbuf.Read(readBuf)

	if !bytes.Equal(readBuf[:bytesWritten], writeBuf) {
		t.Errorf("expected to read same as what i wrote\n")
	}
}
