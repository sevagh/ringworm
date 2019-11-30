Two ringbuffers written in Go.

### ringbuffer1

[ringbuffer1](./ringbuffer1) stores bytes. It implements an indexing strategy learned from https://www.snellman.net/blog/archive/2016-12-13-ring-buffers/

Dependencies for ringbuffer1 are for running tests:

* https://github.com/flyingmutant/rapid
* https://github.com/google/gofuzz

### ringbuffer2

[ringbuffer2](./ringbuffer2) stores interfaces. The API is inspired by https://github.com/armon/circbuf
