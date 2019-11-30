package ringbuffer_test

import (
	"bytes"
	"testing"

	"github.com/flyingmutant/rapid"
	"github.com/sevagh/ringworm/ringbuffer1"
)

type ringbufferMachine struct {
	r     ringbuffer.Ringbuffer
	n     int
	state [][]byte
}

func (m *ringbufferMachine) Init(t *rapid.T) {
	n := rapid.IntsRange(1, 20).Draw(t, "n").(int)
	ringbufSizePowerOfTwo := 1 << n
	m.r, _ = ringbuffer.NewRingbuffer(ringbufSizePowerOfTwo)

	t.Logf("Created ringbuffer with size 2^%d = %d\n", n, ringbufSizePowerOfTwo)
	m.n = n
}

func (m *ringbufferMachine) Get(t *rapid.T) {
	if m.r.Size() == 0 {
		t.Skip("ringbuffer empty")
	}

	currState := m.state[0]
	readBuf := make([]byte, len(currState))

	m.r.Read(readBuf)
	if !bytes.Equal(readBuf, currState) {
		t.Fatalf("got invalid value: %v vs expected %v", readBuf, currState)
	}
	m.state = m.state[1:]
}

func (m *ringbufferMachine) Put(t *rapid.T) {
	if m.r.Size() == (1 << m.n) {
		t.Skip("ringbuffer full")
	}

	writeBuf := rapid.SlicesOfN(rapid.Bytes(), 0, m.n).Draw(t, "writeSlice").([]byte)

	err := m.r.Write(writeBuf)
	if len(writeBuf) != 0 {
		if err == nil {
			m.state = append(m.state, writeBuf)
		} else {
			t.Logf("ringbuffer full, discarded write %d for remaining capacity %d", writeBuf, m.r.Capacity()-m.r.Size())
		}
	}
}

func (m *ringbufferMachine) Check(t *rapid.T) {
	stateSum := 0
	for _, state := range m.state {
		stateSum += len(state)
	}

	if m.r.Size() != stateSum {
		t.Fatalf("ringbuffer size mismatch: %v vs expected %v", m.r.Size(), stateSum)
	}
}

func TestRingbufferProperty(t *testing.T) {
	rapid.Check(t, rapid.StateMachine(&ringbufferMachine{}))
}
