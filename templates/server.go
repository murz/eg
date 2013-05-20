package templates

import (
	"bytes"
	"compress/gzip"
	"io"
)

// Server returns the raw, uncompressed file data data.
func Server() []byte {
	gz, err := gzip.NewReader(bytes.NewBuffer([]byte{
0x1f,0x8b,0x08,0x00,0x00,0x09,0x6e,0x88,0x00,0xff,0x6c,0x91,
0xc1,0x4e,0xf3,0x30,0x10,0x84,0xcf,0xf1,0x53,0x58,0xf9,0x2f,
0x8e,0x54,0xd9,0xf7,0x5f,0xe2,0x80,0x40,0x5c,0x90,0x8a,0x54,
0x10,0x97,0x8a,0x83,0x6b,0x36,0xa9,0x45,0x12,0x5b,0xf6,0xba,
0x50,0x2c,0xbf,0x3b,0x76,0x52,0xd1,0x50,0xf5,0x98,0xdd,0xd9,
0x6f,0x26,0x63,0x2b,0xd5,0x87,0xec,0x80,0x0e,0x52,0x8f,0x84,
0xe8,0xc1,0x1a,0x87,0x94,0x91,0xaa,0xee,0x34,0xee,0xc3,0x8e,
0x2b,0x33,0x88,0x21,0xb8,0x6f,0x01,0x9d,0xa9,0xaf,0x8f,0xc5,
0x1e,0xd1,0x96,0x5d,0x8c,0x74,0x2d,0x07,0xa0,0x29,0x09,0x69,
0xad,0x50,0x66,0x44,0x67,0xfa,0x1e,0x9c,0xbf,0xd8,0xe6,0x4d,
0x5b,0x46,0x0e,0xda,0x1e,0x14,0xd6,0xa4,0x21,0xa4,0x0d,0xa3,
0x9a,0x52,0xb0,0x86,0x46,0x52,0xc5,0xf8,0xef,0x56,0xa1,0x36,
0xa3,0x4f,0x89,0x54,0xc5,0x81,0x6f,0xa0,0xd3,0x1e,0xc1,0xcd,
0x73,0x56,0x80,0x77,0xbf,0x16,0x19,0xcb,0xcf,0x0e,0xf5,0x8a,
0x9e,0xd8,0xfc,0xe5,0x68,0xe1,0xa9,0x65,0x8b,0x30,0xfc,0xf2,
0x30,0xa6,0x66,0x45,0xb7,0x6f,0x1e,0x9d,0x1e,0xbb,0xec,0x5d,
0xcc,0x8b,0x00,0xbe,0xf0,0x11,0x8e,0x53,0x80,0x6a,0xca,0xff,
0x2a,0xfb,0x30,0xe3,0x27,0x91,0xb8,0x10,0xa5,0x55,0xfe,0x01,
0xbb,0x9d,0x39,0x7f,0x71,0x0f,0x1a,0xfa,0xf7,0x05,0x29,0x9f,
0x14,0xce,0x7f,0x7a,0x0d,0x7b,0x16,0xa7,0xa6,0x34,0x21,0x16,
0x4d,0x94,0xea,0xf8,0xc6,0x04,0x04,0xcf,0x9a,0xd3,0xe7,0xbd,
0x44,0xb9,0x93,0x7e,0x9e,0x1c,0xa4,0xa3,0x9e,0xde,0xd0,0xfc,
0x30,0x7c,0x0d,0x9f,0xcf,0xe0,0x0e,0xe0,0xd8,0xa2,0xfd,0x3a,
0x8b,0x3c,0xdf,0x84,0x5c,0x34,0x49,0x3f,0x01,0x00,0x00,0xff,
0xff,0x5d,0x98,0x04,0xa9,0xff,0x01,0x00,0x00,
	}))

	if err != nil {
		panic("Decompression failed: " + err.Error())
	}

	var b bytes.Buffer
	io.Copy(&b, gz)
	gz.Close()

	return b.Bytes()
}