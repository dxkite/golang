package pngio

import (
	"fmt"
	"image"
	"testing"
)

func TestImagePack_Write(t *testing.T) {
	type fields struct {
		Image  *image.RGBA
		offset int
	}

	tests := []struct {
		name    string
		i       *ImagePack
		b       []byte
		wantN   int
		wantErr bool
	}{
		{
			name:    "hello world",
			i:       NewPack(2),
			b:       []byte("hello world"),
			wantN:   len([]byte("hello world")),
			wantErr: false,
		},
		{
			name:    "hello world 12",
			i:       NewPack(12),
			b:       []byte("hello world"),
			wantN:   len([]byte("hello world")),
			wantErr: false,
		},
		{
			name:    "hello world 1024",
			i:       NewPack(1024),
			b:       []byte("hello world"),
			wantN:   len([]byte("hello world")),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := tt.i.Write(tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Write() gotN = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestImagePack_Read(t *testing.T) {
	d := []byte("hello world")
	pack := NewDataPack(d)
	buf := make([]byte, 5)
	if n, err := pack.Read(buf); err == nil {
		fmt.Println("read", string(buf[:n]))
		if string(buf[:n]) != "hello" {
			t.Errorf("Read() got = %v, want hello", string(buf[:n]))
		}
	}

	buf2 := make([]byte, 6)
	if n, err := pack.Read(buf2); err == nil {
		fmt.Println("read", string(buf2[:n]))
		if string(buf2[:n]) != " world" {
			t.Errorf("Read() got = %v, want \" world\"", string(buf2[:n]))
		}
	}
}
