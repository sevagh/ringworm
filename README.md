Two ringbuffers written in Go.

### ringbuffer1

[ringbuffer1](./ringbuffer1) stores bytes. It features:

1. indexing strategy learned from https://www.snellman.net/blog/archive/2016-12-13-ring-buffers/
2. https://github.com/bmkessler/fastdiv for faster modulo on the read/write indices
2. https://github.com/google/gofuzz and https://github.com/flyingmutant/rapid for testing

### ringbuffer2

[ringbuffer2](./ringbuffer2) stores interfaces. The API is inspired by https://github.com/armon/circbuf
