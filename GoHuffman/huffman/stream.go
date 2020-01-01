// ByteSteam
// Use to save bit stream
package huffman

import (
	"fmt"
	"io"
	"strings"
)

type ByteStream struct {
	data   []byte
	length uint64
	offset uint64
}

func NewByteStream(b []byte, l uint64) *ByteStream {
	return &ByteStream{
		data:   b,
		length: l,
	}
}

func NewEmptyByteStream() *ByteStream {
	return &ByteStream{
		data:   []byte{},
		length: 0,
	}
}

// Write as Little-Endian from bit
// Actually store is Big-Endian
func (b *ByteStream) WriteBit(bit byte) {
	var offset = b.length / 8
	var offsetBit = b.length % 8
	b.length++
	if offset >= uint64(len(b.data)) {
		b.data = append(b.data, 0)
	}
	if bit == 1 {
		b.data[offset] = b.data[offset] | (1 << (offsetBit))
	} else {
		b.data[offset] = b.data[offset] & (^(1 << (offsetBit)))
	}
	//fmt.Println("write",bit,"at",offset,"+",offsetBit,b)
}

// Write raw
func (b *ByteStream) WriteByte(data byte) {
	b.WriteByteLen(data, 8)
}

// Write raw
func (b *ByteStream) WriteUint16(data uint16) {
	b.WriteByte(byte(data))      // 写第一字节
	b.WriteByte(byte(data >> 8)) // 写第二字节
}

func (b *ByteStream) WriteUint32(data uint32) {
	b.WriteUint16(uint16(data))         // 写第一字节
	b.WriteUint16(uint16((data >> 16))) // 写第二字节
}

func (b *ByteStream) WriteUint64(data uint64) {
	b.WriteUint32(uint32(data))         // 写第一字节
	b.WriteUint32(uint32((data >> 32))) // 写第二字节
}

func (b *ByteStream) ReadBit() byte {
	var offset = b.offset / 8
	var offsetBit = b.offset % 8
	if offset >= uint64(len(b.data)) {
		return 0
	}
	b.offset++
	if BitIsOne(b.data[offset], int(offsetBit)) {
		return 1
	}
	return 0
}

func (b *ByteStream) ReadByte() byte {
	var offset = b.offset / 8
	var offsetBit = b.offset % 8
	if offsetBit == 0 {
		b.offset += 8
		return b.data[offset]
	}
	return b.ReadByteLen(8)
}

func (b *ByteStream) ReadByteLen(l int) byte {
	var val byte
	for i := 0; i < l; i++ {
		val += b.ReadBit() << i
	}
	return val
}

func (b *ByteStream) ReadBitString(l int) string {
	var val string
	for i := 0; i < l; i++ {
		if b.ReadBit() == 1 {
			val += "1"
		} else {
			val += "0"
		}
	}
	return val
}

func (b *ByteStream) ReadUint16() uint16 {
	return uint16(b.ReadByte()) + (uint16(b.ReadByte()) << 8)
}

func (b *ByteStream) ReadUint32() uint32 {
	return uint32(b.ReadUint16()) + (uint32(b.ReadUint16()) << 16)
}

func (b *ByteStream) ReadUint64() uint64 {
	return uint64(b.ReadUint32()) + (uint64(b.ReadUint32()) << 32)
}

// Write byte raw
func (b *ByteStream) WriteByteLen(data byte, len int) {
	for i := 0; i < len; i++ {
		if BitIsOne(data, i) {
			b.WriteBit(1)
		} else {
			b.WriteBit(0)
		}
	}
}

// Check bit as Little-Endian if is 1
func BitIsOne(d byte, i int) bool {
	return d&(1<<i) > 0
}

// Write as Little-Endian from string bit (eg. 00101)
// Actually store as Big-Endian
func (b *ByteStream) WriteBitString(data string) {
	for _, bit := range data {
		if bit == '0' {
			b.WriteBit(0)
		} else {
			b.WriteBit(1)
		}
	}
}

// Write to stream
func (b *ByteStream) WriteStream(writer io.Writer) (n int, err error) {
	last := b.length % 8
	if last%8 == 0 {
		n, err = writer.Write(b.data)
		b.data = []byte{}
		b.length = 0
	} else {
		n, err = writer.Write(b.data[0 : len(b.data)-1])
		b.data = []byte{b.data[len(b.data)-1]}
		b.length = last
	}
	return
}

// io.Writer
func (b *ByteStream) Write(data []byte) (n int, err error) {
	for _, item := range data {
		b.WriteByte(item)
	}
	return
}

func (b *ByteStream) Read(data []byte) (n int, err error) {
	for index, _ := range data {
		data[index] = b.ReadByte()
	}
	less := int((b.length - b.offset) / 8)
	if less > len(data) {
		return len(data), nil
	}
	return less, nil
}

// 添加到末尾
func (b *ByteStream) Append(stream *ByteStream) {
	last := stream.length % 8
	if last == 0 {
		for _, data := range stream.data {
			b.WriteByte(data)
		}
	} else {
		for _, data := range stream.data[0 : len(stream.data)-1] {
			b.WriteByte(data)
		}
		b.WriteByteLen(stream.data[len(stream.data)-1], int(last))
	}
}

func (b *ByteStream) AppendByte(data []byte)  {
	for _, item := range data {
		b.WriteByte(item)
	}
	return
}

// Little-Endian Like Stream
func (b *ByteStream) String() string {
	var lastShow = b.length % 8
	size := len(b.data)
	var data string
	for index, b := range b.data {
		if size == index+1 && lastShow > 0 {
			data += getBitString(b, int(lastShow)) + " "
		} else {
			data += getBitString(b, 8) + " "
		}
	}
	return strings.TrimSpace(data)
}

// Little-Endian Like Stream with hex dump
func (b *ByteStream) HexString() string {
	var lastShow = b.length % 8
	size := len(b.data)
	var data = ""
	for index, b := range b.data {
		if size == index+1 && lastShow > 0 {
			data += fmt.Sprintf("%02X", b) + "[" + getBitString(b, int(lastShow)) + "] "
		} else {
			data += fmt.Sprintf("%02X", b) + " "
		}
	}
	return strings.TrimSpace(data)
}

// Write byte raw
func getBitString(data byte, len int) (hex string) {
	for i := 0; i < len; i++ {
		if BitIsOne(data, i) {
			hex += "1"
		} else {
			hex += "0"
		}
	}
	return
}


