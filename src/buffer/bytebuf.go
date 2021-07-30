package buffer

import "math"

type ByteBuffer struct {
	Data []byte
}

type ByteBuf struct {
	Data []byte
	ReaderIndex uint32
	WriterIndex uint32
	Capacity uint32
	ByteOrder ByteOrder
}

func New(capacity uint32, order ByteOrder) *ByteBuf {
	var buf = ByteBuf{}
	buf.Data = make([]byte, capacity)
	buf.Capacity = capacity
	buf.ByteOrder = order

	return &buf
}

func (b *ByteBuf) ReadableBytes() uint32 {
	return b.WriterIndex - b.ReaderIndex
}

func (b *ByteBuf) WritableBytes() uint32 {
	return b.Capacity - b.WriterIndex
}

func (b *ByteBuf) WriteInt(value int32) {
	b.ensureWritable(4)

	b.ByteOrder.PutUint32(b.Data, b.WriterIndex, uint32(value))
	b.WriterIndex += 4
}

func (b *ByteBuf) WriteLong(value int64) {
	b.ensureWritable(8)

	b.ByteOrder.PutUint64(b.Data, b.WriterIndex, uint64(value))
	b.WriterIndex += 8
}

func (b *ByteBuf) WriteFloat(value float32) {
	b.ensureWritable(4)

	b.ByteOrder.PutUint32(b.Data, b.WriterIndex, math.Float32bits(value))
	b.WriterIndex += 4
}

func (b *ByteBuf) WriteDouble(value float64) {
	b.ensureWritable(8)

	b.ByteOrder.PutUint64(b.Data, b.WriterIndex, math.Float64bits(value))
	b.WriterIndex += 8
}

func (b *ByteBuf) WriteByte(value byte) {
	b.ensureWritable(1)

	b.Data[b.WriterIndex] = value
	b.WriterIndex += 1
}

func (b *ByteBuf) WriteBytes(value []byte) {
	b.WriteBytesByLen(value, 0, uint32(len(value)))
}

func (b *ByteBuf) WriteBytesByLen(value []byte, index uint32, len uint32) {
	b.ensureWritable(len)

	copy(b.Data[b.WriterIndex:], value[index:index + len])
	b.WriterIndex += len
}

func (b *ByteBuf) WriteByteBuf(value ByteBuf) {
	b.WriteBytesByLen(value.Data, value.ReaderIndex, value.ReadableBytes())
}

func (b *ByteBuf) WriteBool(value bool) {
	if value {
		b.WriteByte(1)
	} else {
		b.WriteByte(0)
	}
}

func (b *ByteBuf) ReadBool() bool {
	return b.ReadByte() == 1
}

func (b *ByteBuf) ReadByte() byte {
	var value = b.GetByte()
	b.ReaderIndex++
	return value
}

func (b *ByteBuf) ReadInt32() int32 {
	var value = b.GetInt32()
	b.ReaderIndex += 4
	return value
}

func (b *ByteBuf) ReadUInt32() uint32 {
	var value = b.GetUInt32()
	b.ReaderIndex += 4
	return value
}

func (b *ByteBuf) ReadInt64() int64 {
	var value = b.GetInt64()
	b.ReaderIndex += 8
	return value
}

func (b *ByteBuf) ReadUInt64() uint64 {
	var value = b.GetUInt64()
	b.ReaderIndex += 8
	return value
}

func (b *ByteBuf) ReadFloat() float32 {
	return math.Float32frombits(b.ReadUInt32())
}

func (b *ByteBuf) ReadDouble() float64 {
	return math.Float64frombits(b.ReadUInt64())
}

func (b *ByteBuf) ReadBytes(len uint32) []byte {
	v := b.GetBytes(len)
	b.ReaderIndex += len
	return v
}

func (b *ByteBuf) GetByte() byte {
	return b.Data[b.ReaderIndex]
}

func (b *ByteBuf) GetInt32() int32 {
	v := b.ByteOrder.Uint32(b.Data, b.ReaderIndex)
	return int32(v)
}

func (b *ByteBuf) GetUInt32() uint32 {
	v := b.ByteOrder.Uint32(b.Data, b.ReaderIndex)
	return v
}

func (b *ByteBuf) GetInt64() int64 {
	v := b.ByteOrder.Uint64(b.Data, b.ReaderIndex)
	return int64(v)
}

func (b *ByteBuf) GetUInt64() uint64 {
	v := b.ByteOrder.Uint64(b.Data, b.ReaderIndex)
	return v
}

func (b *ByteBuf) GetBytes(len uint32) []byte {
	v := b.Data[b.ReaderIndex : b.ReaderIndex + len]
	return v
}

func (b *ByteBuf) ensureWritable(minWritableBytes uint32) {
	if minWritableBytes <= b.WritableBytes() {
		return
	}

	var minNewCapacity = b.WriterIndex + minWritableBytes
	var newCapacity = b.calculateNewCapacity(minNewCapacity)
	var expandData = make([]byte, newCapacity - b.Capacity)
	b.Data = append(b.Data, expandData...)
}

func (b *ByteBuf) calculateNewCapacity(minNewCapacity uint32) uint32 {
	var newCapacity uint32 = 64
	for newCapacity < minNewCapacity {
		newCapacity <<= 1
	}
	if newCapacity > math.MaxInt32 {
		return math.MaxInt32
	}
	return newCapacity
}

func (b *ByteBuf) SkipBytes(len uint32) {
	b.ReaderIndex += len
	if b.ReaderIndex > b.Capacity {
		b.ReaderIndex = b.Capacity
	}
}

func (b *ByteBuf) Reset() {
	b.WriterIndex = 0
	b.ReaderIndex = 0
}