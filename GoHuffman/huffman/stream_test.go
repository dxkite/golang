package huffman

import (
	"fmt"
	"testing"
)

func TestByteStream_WriteByte(t *testing.T) {
	b := NewByteStream([]byte{0xAC}, 8)
	data := byte(0xCA)
	b.WriteByte(data)
	if b.data[0] == 0xAC && b.data[1] == data {
	} else {
		t.Error("write byte error")
	}
}

func TestByteStream_WriteByte2(t *testing.T) {
	b := NewByteStream([]byte{0xAC}, 8)
	// Write bit as little-endian
	// Write 1100 to stream
	b.WriteBit(1)
	b.WriteBit(1)
	b.WriteBit(0)
	b.WriteBit(0)
	// Actually store as big-endian
	if b.data[0] == 0xAC && b.data[1] == 0b0011 {
	} else {
		t.Error("write byte error", b)
	}
}

func TestByteStream_WriteBitString(t *testing.T) {
	b := NewByteStream([]byte{0xAC}, 8)
	// Write bit as little-endian
	// Write 1100 to stream
	b.WriteBitString("1100")
	b.WriteBitString("00111")
	// Actually stored is big-endian
	// After b.String() is like 00110101 11000011 1
	// Actually is 10101100 11000011 1 (Big-Endian)
	// 0011 -> 1100 -> C
	// 0101 -> 1010 -> A
	if b.data[0] == 0xAC && b.data[1] == 0b11000011 && b.data[2] == 0b1 {
	} else {
		t.Error("write byte error", b)
	}
}

func TestByteStream_Append(t *testing.T) {
	b := NewByteStream([]byte{0b11111101}, 3)
	c := NewByteStream([]byte{0b11001111, 0b0001011}, 13)
	b.Append(c)
	// Hex is 101,11110 011,11010
	if b.data[0] == 0x7D && b.data[1] == 0x5E {
	} else {
		t.Error("write byte error", b)
	}
}

func TestByteStream_WriteUint16(t *testing.T) {
	b := NewEmptyByteStream()
	b.WriteUint16(0b1100111100010110)
	// Hex is 01101000 11110011
	if b.data[0] == 22 && b.data[1] == 207 {
	} else {
		t.Error("write byte error", b)
	}
}

func TestByteStream_WriteUint32(t *testing.T) {
	b := NewEmptyByteStream()
	b.WriteUint32(0b11001111000101111001111000101100)
	// Hex is 0011 0100 0111 1001 1110 1000 1111 0011
	if b.data[0] == 44 && b.data[1] == 158 && b.data[2] == 23 && b.data[3] == 207 {
	} else {
		t.Error("write byte error", b)
	}
}

func TestByteStream_WriteUint64(t *testing.T) {
	b := NewEmptyByteStream()
	b.WriteUint64(0b1100111100010111100111100010110011001111000101111001111000101100)
	// Hex is 0011 0100 0111 1001 1110 1000 1111 0011 0011 0100 0111 1001 1110 1000 1111 0011
	start32 := b.data[0] == 44 && b.data[1] == 158 && b.data[2] == 23 && b.data[3] == 207
	stop32 := b.data[4] == 44 && b.data[5] == 158 && b.data[6] == 23 && b.data[7] == 207
	if start32 && stop32 {
	} else {
		t.Error("write byte error", b)
	}
}

func TestByteStream_ReadByte(t *testing.T) {
	type fields struct {
		data   []byte
		length uint64
		offset uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		{
			name: "simple offset 0",
			fields: struct {
				data   []byte
				length uint64
				offset uint64
			}{data: []byte{0b10010111}, length: 8, offset: 0},
			want: 0b10010111,
		},
		{
			name: "simple offset 1",
			fields: struct {
				data   []byte
				length uint64
				offset uint64
			}{data: []byte{0b10010111, 0b11101110}, length: 8, offset: 1},
			want: 0b01001011,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ByteStream{
				data:   tt.fields.data,
				length: tt.fields.length,
				offset: tt.fields.offset,
			}
			if got := b.ReadByte(); got != tt.want {
				t.Errorf("ReadByte() = %08b, want %08b", got, tt.want)
			}
		})
	}
}

func TestByteStream_ReadUint16(t *testing.T) {
	tests := [] uint16{0xccdd, 0xeeff, 0x00ed}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("0x%04x", tt), func(t *testing.T) {
			b := NewEmptyByteStream()
			b.WriteUint16(tt)
			if got := b.ReadUint16(); got != tt {
				t.Errorf("ReadByte() = %016b, want %016b", got, tt)
			}
		})
	}
}

func TestByteStream_ReadUint16_WithSomeStart(t *testing.T) {
	tests := [] uint16{0xccdd, 0xefce, 0xcece}
	for index, tt := range tests {
		t.Run(fmt.Sprintf("0x%04x", tt), func(t *testing.T) {
			b := NewEmptyByteStream()
			val := uint8(index % 2)
			b.WriteBit(val)
			b.WriteUint16(tt)
			b.WriteBit(val)
			read := b.ReadBit()
			if got := b.ReadUint16(); got != tt && read == val {
				t.Errorf("ReadUint16() = %016b, want %016b", got, tt)
			}
		})
	}
}
