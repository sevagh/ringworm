package ringbuffer_test

import (
	"encoding/binary"
	"testing"

	"github.com/sevagh/ringworm/ringbuffer1"
)

func BenchmarkManyRingbuffersBillionsOfIntegers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ringbuf := ringbuffer.NewRingbuffer(64)
		writeBuf := make([]byte, 4)
		readBuf := make([]byte, 4)

		num := uint32(0)
		for j := uint32(0); j < 1000000000; j++ {
			num = num*17 + 255

			binary.LittleEndian.PutUint32(writeBuf, num)

			ringbuf.Write(writeBuf)
			ringbuf.Read(readBuf)

			read := binary.LittleEndian.Uint32(readBuf)

			if num != read {
				b.Errorf("num not same after round trip")
			}
		}
	}
}
