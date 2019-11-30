package ringbuffer_test

import (
	"testing"

	"github.com/sevagh/ringworm/ringbuffer1"
)

var CoolText []byte = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce tincidunt nisi tincidunt velit euismod gravida. Aenean at turpis nec lectus faucibus lobortis vitae at arcu. In hac habitasse platea dictumst. Praesent commodo tellus vitae massa maximus porta. Etiam tortor justo, pulvinar nec tempus sit amet, sodales sit amet libero. Nullam non massa mi. Aenean laoreet arcu nec efficitur interdum. Donec vel neque id enim ullamcorper congue nec cursus metus. In hac habitasse platea dictumst. Vestibulum gravida enim non luctus consectetur. Maecenas sit amet risus id orci ultrices viverra dapibus a justo. Nulla facilisi. Nulla velit justo, convallis vel rhoncus a, laoreet ac lectus. Duis ex augue, facilisis id tortor eu, vulputate consectetur felis. Curabitur vitae aliquam nisi. Fusce blandit placerat metus eu elementum. Nam sed lectus tellus. Etiam nec eros sed turpis pulvinar interdum in sed quam. Phasellus tempus lobortis justo. Nullam vulputate nisl sed felis porta, sit amet auctor lectus porttitor. Aenean neque quam, luctus viverra felis eu, feugiat venenatis diam. Nulla facilisi. Nullam luctus augue erat, sit amet suscipit tortor eleifend non. Maecenas congue luctus nisi, ac ullamcorper ligula tristique nec. Nullam non magna bibendum, pretium lorem eget, faucibus nisl. Praesent risus tellus, commodo vitae convallis ac, dictum at felis. Cras lacus risus, accumsan nec lectus sit amet, consectetur pharetra ante. Nulla eget bibendum erat, ut placerat sapien. Mauris nulla magna, fringilla sed turpis eget, sodales luctus orci. Vestibulum lectus sapien, facilisis nec eros vitae, dictum feugiat purus. Phasellus felis nunc, condimentum et dolor mattis, sagittis varius mauris. Proin diam turpis, efficitur vitae ipsum convallis, cursus imperdiet libero. Etiam porta sapien sapien, quis fringilla magna porttitor a. Maecenas ac urna vitae ante dapibus convallis. Pellentesque sit amet varius justo. Morbi tristique elementum pellentesque. Vivamus id sodales risus, eget condimentum libero. Sed maximus, metus vitae mollis ultricies, ligula purus sodales enim, ut accumsan felis nibh non enim. Phasellus in orci sem. Sed ac felis rutrum, aliquet felis eget, tristique sem. Quisque non justo pharetra, condimentum nibh sit amet, condimentum neque. Donec a odio at quam hendrerit eleifend at ac mi. Vivamus a iaculis sem, quis tempor ante. Vestibulum quis ligula nisi. Vestibulum commodo orci at lectus consequat, consectetur eleifend est tristique. Proin consectetur elementum metus eget fermentum. Sed facilisis sapien id orci cursus pellentesque.")

func BenchmarkRingbufferSingleByteManyRoundTrips(b *testing.B) {
	//ringbuffer of size one so that every read and write of 1 byte in my bench loop is
	// "maximally expensive" (wraparound logic etc.)? or is that crazy
	ringbuf, _ := ringbuffer.NewRingbuffer(1)
	oneByte := make([]byte, 1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, b := range CoolText {
			//just round-trip single bytes into the ringbuffer
			oneByte[0] = b
			ringbuf.Write(oneByte)
			ringbuf.Read(oneByte)
		}
	}
}

func BenchmarkRingbufferManyBytesManyRoundTrips(b *testing.B) {
	ringbuf, _ := ringbuffer.NewRingbuffer(64)
	readBuf := make([]byte, 32)
	var chunk int = 0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for chunk = 0; chunk < len(CoolText)-32; chunk += 32 {
			ringbuf.Write(CoolText[chunk : chunk+32])
			ringbuf.Read(readBuf)
		}
		//chunk -= 32
		remain := len(CoolText) - chunk
		ringbuf.Read(readBuf[:remain])
	}
}
