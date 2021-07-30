package buffer


// A ByteOrder specifies how to convert byte sequences into
// 16-, 32-, or 64-bit unsigned integers.
type ByteOrder interface {
	Uint16([]byte, uint32) uint16
	Uint32([]byte, uint32) uint32
	Uint64([]byte, uint32) uint64
	PutUint16([]byte, uint32, uint16)
	PutUint32([]byte, uint32, uint32)
	PutUint64([]byte, uint32, uint64)
	String() string
}

// LittleEndian is the little-endian implementation of ByteOrder.
var LittleEndian littleEndian

// BigEndian is the big-endian implementation of ByteOrder.
var BigEndian bigEndian

type littleEndian struct{}

func (littleEndian) Uint16(b []byte, index uint32) uint16 {
	_ = b[index + 1] // bounds check hint to compiler; see golang.org/issue/14808
	return uint16(b[index]) | uint16(b[index + 1])<<8
}

func (littleEndian) PutUint16(b []byte, index uint32, v uint16) {
	_ = b[index + 1] // early bounds check to guarantee safety of writes below
	b[index] = byte(v)
	b[index + 1] = byte(v >> 8)
}

func (littleEndian) Uint32(b []byte, index uint32) uint32 {
	_ = b[index + 3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[index]) | uint32(b[index + 1])<<8 | uint32(b[index + 2])<<16 | uint32(b[index + 3])<<24
}

func (littleEndian) PutUint32(b []byte, index uint32, v uint32) {
	_ = b[index + 3] // early bounds check to guarantee safety of writes below
	b[index] = byte(v)
	b[index + 1] = byte(v >> 8)
	b[index + 2] = byte(v >> 16)
	b[index + 3] = byte(v >> 24)
}

func (littleEndian) Uint64(b []byte, index uint32) uint64 {
	_ = b[index + 7] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[index]) | uint64(b[index + 1])<<8 | uint64(b[index + 2])<<16 | uint64(b[index + 3])<<24 |
		uint64(b[index + 4])<<32 | uint64(b[index + 5])<<40 | uint64(b[index + 6])<<48 | uint64(b[index + 7])<<56
}

func (littleEndian) PutUint64(b []byte, index uint32, v uint64) {
	_ = b[index + 7] // early bounds check to guarantee safety of writes below
	b[index] = byte(v)
	b[index + 1] = byte(v >> 8)
	b[index + 2] = byte(v >> 16)
	b[index + 3] = byte(v >> 24)
	b[index + 4] = byte(v >> 32)
	b[index + 5] = byte(v >> 40)
	b[index + 6] = byte(v >> 48)
	b[index + 7] = byte(v >> 56)
}

func (littleEndian) String() string { return "LittleEndian" }

func (littleEndian) GoString() string { return "binary.LittleEndian" }

type bigEndian struct{}

func (bigEndian) Uint16(b []byte, index uint32) uint16 {
	_ = b[index + 1] // bounds check hint to compiler; see golang.org/issue/14808
	return uint16(b[index + 1]) | uint16(b[index])<<8
}

func (bigEndian) PutUint16(b []byte, index uint32, v uint16) {
	_ = b[index + 1] // early bounds check to guarantee safety of writes below
	b[index] = byte(v >> 8)
	b[index + 1] = byte(v)
}

func (bigEndian) Uint32(b []byte, index uint32) uint32 {
	_ = b[index + 3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[index + 3]) | uint32(b[index + 2])<<8 | uint32(b[index + 1])<<16 | uint32(b[index])<<24
}

func (bigEndian) PutUint32(b []byte, index uint32, v uint32) {
	_ = b[index + 3] // early bounds check to guarantee safety of writes below
	b[index] = byte(v >> 24)
	b[index + 1] = byte(v >> 16)
	b[index + 2] = byte(v >> 8)
	b[index + 3] = byte(v)
}

func (bigEndian) Uint64(b []byte, index uint32) uint64 {
	_ = b[index + 7] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[index + 7]) | uint64(b[index + 6])<<8 | uint64(b[index + 5])<<16 | uint64(b[index + 4])<<24 |
		uint64(b[index + 3])<<32 | uint64(b[index + 2])<<40 | uint64(b[index + 1])<<48 | uint64(b[index])<<56
}

func (bigEndian) PutUint64(b []byte, index uint32, v uint64) {
	_ = b[index + 7] // early bounds check to guarantee safety of writes below
	b[index] = byte(v >> 56)
	b[index + 1] = byte(v >> 48)
	b[index + 2] = byte(v >> 40)
	b[index + 3] = byte(v >> 32)
	b[index + 4] = byte(v >> 24)
	b[index + 5] = byte(v >> 16)
	b[index + 6] = byte(v >> 8)
	b[index + 7] = byte(v)
}

func (bigEndian) String() string { return "BigEndian" }

func (bigEndian) GoString() string { return "binary.BigEndian" }