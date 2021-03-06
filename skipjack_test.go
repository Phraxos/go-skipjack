package skipjack

import (
	"bytes"
	"testing"
)

// http://csrc.nist.gov/groups/ST/toolkit/documents/skipjack/skipjack.pdf
var skipjackTestVectors = []struct {
	key    []byte
	plain  []byte
	cipher []byte
}{
	{
		[]byte{0x00, 0x99, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11},
		[]byte{0x33, 0x22, 0x11, 0x00, 0xdd, 0xcc, 0xbb, 0xaa},
		[]byte{0x25, 0x87, 0xca, 0xe2, 0x7a, 0x12, 0xd3, 0x00},
	},
}

// http://csrc.nist.gov/publications/nistpubs/800-17/800-17.pdf
var skipjackVariablePlaintextValidation = []struct {
	plain  []byte
	cipher []byte
}{
	{[]byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x9A, 0x90, 0xBC, 0x0B, 0x75, 0xC7, 0x37, 0x03}},
	{[]byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xCC, 0x68, 0x43, 0x59, 0x8C, 0x73, 0x2B, 0xBE}},
	{[]byte{0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x13, 0x72, 0x95, 0x35, 0x09, 0xB3, 0xC1, 0x4C}},
	{[]byte{0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x70, 0xAA, 0xAA, 0x84, 0x18, 0xE4, 0x89, 0x30}},
	{[]byte{0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xE4, 0xB0, 0xB4, 0xA1, 0x39, 0xE8, 0x54, 0x6E}},
	{[]byte{0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x70, 0x18, 0xF7, 0x13, 0x66, 0x14, 0x6E, 0xAF}},
	{[]byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xB3, 0x8F, 0x3D, 0x7E, 0x4F, 0x2D, 0x25, 0x3D}},
	{[]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xD6, 0x4B, 0xA2, 0x06, 0x51, 0x13, 0xD9, 0x1E}},
	{[]byte{0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xF9, 0x5B, 0x92, 0x2F, 0x14, 0x27, 0xA9, 0xF2}},
	{[]byte{0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x6B, 0x64, 0x2F, 0xDE, 0x40, 0x85, 0x85, 0x86}},
	{[]byte{0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x6C, 0xF5, 0x2D, 0x5E, 0x61, 0x69, 0x52, 0x17}},
	{[]byte{0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xBC, 0x0F, 0x6B, 0xCA, 0x62, 0xE1, 0x39, 0xA6}},
	{[]byte{0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x6A, 0xD5, 0x03, 0xDC, 0x2A, 0xB0, 0xBF, 0xE2}},
	{[]byte{0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xAF, 0xAD, 0xD7, 0xCA, 0xB6, 0x72, 0x35, 0x16}},
	{[]byte{0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x42, 0x1B, 0x89, 0x5A, 0xF5, 0xC0, 0x0A}},
	{[]byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xCA, 0xD0, 0x45, 0x6C, 0xF8, 0x6C, 0xD5, 0x98}},
	{[]byte{0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x16, 0xF4, 0x1C, 0x8F, 0x8A, 0x6A, 0x5B, 0x79}},
	{[]byte{0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x4C, 0xE7, 0x71, 0xC7, 0x51, 0xBA, 0x27, 0x60}},
	{[]byte{0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x72, 0xC9, 0x02, 0xE5, 0x8C, 0xE5, 0x5B, 0x87}},
	{[]byte{0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x6D, 0x37, 0x8C, 0x66, 0x64, 0xD0, 0x01, 0x10}},
	{[]byte{0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xAC, 0x27, 0xB8, 0x5B, 0x0A, 0x75, 0xE8, 0xBA}},
	{[]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x54, 0xDF, 0x3A, 0x75, 0x5B, 0x00, 0x63, 0xD2}},
	{[]byte{0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x31, 0x4F, 0x4D, 0x28, 0x6D, 0xB4, 0x90, 0x58}},
	{[]byte{0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x88, 0xAE, 0x06, 0x66, 0xB2, 0xA0, 0x78, 0x46}},
	{[]byte{0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00}, []byte{0xD8, 0x60, 0xA8, 0xD9, 0xA0, 0x2C, 0xBC, 0xE8}},
	{[]byte{0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00}, []byte{0x37, 0xCE, 0x5E, 0xEA, 0x53, 0x13, 0x53, 0x5D}},
	{[]byte{0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00}, []byte{0x73, 0x3A, 0xF9, 0x2D, 0xA1, 0xC1, 0x80, 0x26}},
	{[]byte{0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00}, []byte{0x34, 0x1C, 0x23, 0x5F, 0x6E, 0x32, 0x98, 0x1D}},
	{[]byte{0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00}, []byte{0xC6, 0xA6, 0x56, 0x14, 0x47, 0xD9, 0xE0, 0x96}},
	{[]byte{0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00}, []byte{0xC5, 0x50, 0x66, 0xA8, 0xD8, 0x39, 0xE5, 0xFA}},
	{[]byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00}, []byte{0x65, 0x86, 0x4B, 0x48, 0x79, 0x11, 0xA1, 0x0C}},
	{[]byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}, []byte{0x87, 0x29, 0x07, 0xE2, 0xD3, 0x36, 0x33, 0x2A}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00}, []byte{0xAF, 0x03, 0x76, 0x88, 0xE7, 0xA5, 0x24, 0x9C}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00}, []byte{0xC1, 0xFC, 0xD1, 0xB4, 0xDC, 0xC2, 0xAC, 0xBB}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00}, []byte{0x40, 0x48, 0x48, 0x80, 0x2D, 0x69, 0x3D, 0xDA}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00}, []byte{0xB2, 0xDC, 0xCE, 0xE3, 0x3B, 0x15, 0x6D, 0xB6}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00}, []byte{0xE6, 0x20, 0xF4, 0x2A, 0x7F, 0xA9, 0x01, 0x0B}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00}, []byte{0x7C, 0xF0, 0x67, 0xF3, 0xBD, 0x3E, 0xC3, 0x53}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00}, []byte{0x06, 0x37, 0x78, 0x1F, 0x1A, 0x34, 0x72, 0x81}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, []byte{0x47, 0x41, 0xF1, 0x46, 0x4B, 0x71, 0x70, 0x8E}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00}, []byte{0xED, 0xAD, 0x33, 0xF4, 0x56, 0xF5, 0x14, 0xDF}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00}, []byte{0xED, 0x81, 0x27, 0x48, 0xB7, 0xF5, 0x23, 0xE9}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00}, []byte{0x83, 0x8C, 0x9C, 0xC3, 0x83, 0xD4, 0x62, 0x97}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00}, []byte{0xFB, 0x2B, 0xC0, 0xFC, 0xC9, 0x2F, 0x9B, 0x24}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00}, []byte{0xE5, 0x9A, 0xA1, 0x12, 0x2A, 0x65, 0x44, 0x32}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00}, []byte{0xD4, 0xC8, 0xEF, 0x7E, 0x06, 0x43, 0x12, 0x53}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00}, []byte{0x32, 0xED, 0x63, 0x28, 0x14, 0xC2, 0xA8, 0x56}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00}, []byte{0x5D, 0xC2, 0x9F, 0x7D, 0xE9, 0x6E, 0xE5, 0x2C}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00}, []byte{0x68, 0xA0, 0x7C, 0x7E, 0x8E, 0xAD, 0xD5, 0x61}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00}, []byte{0xB2, 0x70, 0x68, 0xF2, 0xD6, 0xB3, 0x37, 0xE2}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00}, []byte{0x1A, 0xF5, 0x1E, 0x9C, 0x29, 0xBF, 0xDC, 0x7B}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00}, []byte{0x92, 0x1D, 0xBD, 0x9B, 0x1C, 0x6B, 0xEA, 0xEB}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x00}, []byte{0x5B, 0x6A, 0x60, 0x22, 0x35, 0x94, 0x35, 0xD2}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00}, []byte{0xD7, 0x74, 0xC6, 0x23, 0x74, 0xB2, 0x3B, 0x09}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00}, []byte{0xFD, 0x9F, 0x05, 0x27, 0x59, 0x4C, 0xE3, 0x7B}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00}, []byte{0x67, 0x86, 0x01, 0xC8, 0xB3, 0x64, 0xA7, 0x94}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}, []byte{0xD5, 0x18, 0x22, 0x8D, 0x5B, 0x0B, 0xE3, 0xD7}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40}, []byte{0xA4, 0x5F, 0xEE, 0x6B, 0xDD, 0x1F, 0x73, 0x4A}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20}, []byte{0xD1, 0xBA, 0x95, 0x51, 0xDF, 0x7C, 0xD5, 0x68}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10}, []byte{0xAE, 0xA3, 0x3D, 0x09, 0xDC, 0x9D, 0x13, 0x10}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08}, []byte{0x96, 0xB4, 0x91, 0xC1, 0xFE, 0x44, 0x3E, 0x9A}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04}, []byte{0xD0, 0xE0, 0x14, 0xCF, 0xEE, 0x94, 0x58, 0x9D}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02}, []byte{0x0B, 0x9E, 0x44, 0xB5, 0x37, 0xAF, 0x28, 0x79}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, []byte{0x22, 0xF4, 0x28, 0xE3, 0xEC, 0x49, 0x1E, 0x60}},
}

var skipjackVariableKeyValidation = []struct {
	key    []byte
	cipher []byte
}{
	{[]byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x7A, 0x00, 0xE4, 0x94, 0x41, 0x46, 0x1F, 0x5A}},
	{[]byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xA1, 0x4F, 0xF8, 0xBC, 0xD1, 0xBC, 0x9E, 0xF9}},
	{[]byte{0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xD7, 0xE8, 0x10, 0x38, 0x5A, 0x42, 0xAA, 0xEA}},
	{[]byte{0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x28, 0xFE, 0x2C, 0x33, 0x32, 0xAA, 0xBD, 0x35}},
	{[]byte{0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x3F, 0xC0, 0xF0, 0x5E, 0xE6, 0xCE, 0x78, 0x8A}},
	{[]byte{0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x44, 0x3D, 0xD0, 0xCB, 0x75, 0x26, 0xF7, 0x4B}},
	{[]byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xAD, 0x81, 0x9E, 0x67, 0x7C, 0xF9, 0x03, 0x05}},
	{[]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x98, 0x91, 0x75, 0x5E, 0x5E, 0xBA, 0x5B, 0x1D}},
	{[]byte{0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x0E, 0x64, 0xB4, 0x94, 0x63, 0x3B, 0xF2, 0xCB}},
	{[]byte{0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x63, 0x38, 0x1A, 0x08, 0xA4, 0x7F, 0xC4, 0x8D}},
	{[]byte{0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xF4, 0x10, 0x8B, 0x09, 0x9B, 0x04, 0x70, 0x40}},
	{[]byte{0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x74, 0x02, 0x16, 0x61, 0x4E, 0xD0, 0xE2, 0x5B}},
	{[]byte{0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x80, 0x00, 0x91, 0x7B, 0x2E, 0x16, 0xB9, 0x2A}},
	{[]byte{0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xA9, 0x76, 0x9B, 0x62, 0xB3, 0xA0, 0xBE, 0x4E}},
	{[]byte{0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x42, 0xFD, 0xB8, 0x72, 0xEA, 0x31, 0x41, 0x21}},
	{[]byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x1D, 0x67, 0x2B, 0xA0, 0x15, 0x6A, 0xB3, 0x9D}},
	{[]byte{0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xF4, 0x44, 0x41, 0xD7, 0xC7, 0x77, 0xF0, 0x57}},
	{[]byte{0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xEA, 0x48, 0x7D, 0xDC, 0x36, 0x0D, 0x15, 0x94}},
	{[]byte{0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x32, 0x4B, 0x0E, 0x78, 0x5F, 0xF2, 0xB9, 0x08}},
	{[]byte{0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x1A, 0xF5, 0x9E, 0xC2, 0xB9, 0xD6, 0x4C, 0x4F}},
	{[]byte{0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x81, 0x9B, 0x7E, 0x10, 0x2E, 0x76, 0xA0, 0xEE}},
	{[]byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x0B, 0x0B, 0xFE, 0x0D, 0x4A, 0x37, 0xAA, 0x9E}},
	{[]byte{0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x12, 0xB4, 0x3E, 0x37, 0x60, 0xD3, 0x0D, 0xA6}},
	{[]byte{0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x31, 0x77, 0x25, 0x6C, 0x46, 0x15, 0x41, 0xEE}},
	{[]byte{0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x36, 0x00, 0xEB, 0x92, 0x83, 0x6C, 0xA0, 0x26}},
	{[]byte{0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x75, 0xA4, 0x35, 0xAD, 0x22, 0xEC, 0xF7, 0x93}},
	{[]byte{0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x71, 0x90, 0xAA, 0x99, 0x13, 0xC1, 0xF9, 0xEC}},
	{[]byte{0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xAB, 0xA7, 0x18, 0xB1, 0x85, 0xA1, 0x1D, 0xD0}},
	{[]byte{0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x40, 0xF6, 0x7A, 0xBF, 0xCC, 0x3B, 0x87, 0x3C}},
	{[]byte{0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x38, 0xA0, 0xA5, 0x8F, 0xB0, 0x97, 0x28, 0xF2}},
	{[]byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xCA, 0x70, 0x2E, 0x49, 0xBF, 0x6F, 0xA6, 0x45}},
	{[]byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x45, 0x5D, 0x93, 0xF0, 0x39, 0xEA, 0x08, 0x60}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x53, 0x47, 0x64, 0x3F, 0xE8, 0x03, 0x88, 0x3F}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xF4, 0x0F, 0xF1, 0xDC, 0xBA, 0x2B, 0xC1, 0xE5}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x57, 0x4A, 0x48, 0x48, 0x36, 0x9D, 0x41, 0x2E}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0xB2, 0xBE, 0x93, 0x6E, 0x36, 0x67, 0x06, 0x36}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x5C, 0x88, 0x51, 0x7D, 0x27, 0x42, 0xE6, 0x19}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x99, 0x3C, 0x89, 0xD0, 0x9A, 0x2F, 0xE5, 0x56}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x1A, 0x3F, 0x72, 0xDA, 0x69, 0x4C, 0x9F, 0xC7}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte{0x96, 0x59, 0xD5, 0x22, 0x8F, 0x4C, 0xB1, 0x51}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00}, []byte{0x7C, 0x13, 0xF4, 0x9E, 0x75, 0x0F, 0x5C, 0x30}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00}, []byte{0x35, 0x00, 0xBD, 0x40, 0x7B, 0xCD, 0x01, 0xF6}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00}, []byte{0x85, 0xC5, 0x8E, 0x3C, 0x49, 0x44, 0x20, 0x28}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00}, []byte{0x84, 0x13, 0x84, 0x0A, 0x2D, 0x48, 0xAB, 0xEA}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00}, []byte{0x83, 0x28, 0x50, 0xE6, 0xE5, 0xC4, 0xAE, 0x5A}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00}, []byte{0x29, 0xE9, 0x7F, 0x0D, 0x9F, 0x0F, 0xDC, 0x5F}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00}, []byte{0x2C, 0x45, 0x23, 0x04, 0x37, 0xFF, 0x2E, 0x04}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}, []byte{0x10, 0xC4, 0x09, 0xFB, 0x87, 0x2A, 0x98, 0x4F}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00}, []byte{0x14, 0x69, 0x3B, 0x30, 0xC3, 0xAF, 0x74, 0x70}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00}, []byte{0x91, 0x3A, 0x90, 0x50, 0xD5, 0x85, 0xBA, 0xB9}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00}, []byte{0x5B, 0xFB, 0x0F, 0x83, 0xAB, 0x0C, 0x6E, 0xEA}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00}, []byte{0x6C, 0x0C, 0xA7, 0x28, 0x4D, 0x83, 0x6A, 0xAE}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00}, []byte{0xAC, 0x57, 0x27, 0xD6, 0x12, 0xE1, 0x85, 0xE8}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00}, []byte{0x38, 0xD7, 0xD5, 0x96, 0xA3, 0xD2, 0x9D, 0x90}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00}, []byte{0x78, 0xBA, 0xDA, 0xD3, 0xBC, 0x43, 0x6C, 0xA2}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}, []byte{0xE4, 0x05, 0x77, 0x87, 0x41, 0xB0, 0x4B, 0xA0}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00}, []byte{0x72, 0xFF, 0xE4, 0x3D, 0xEA, 0x02, 0xAF, 0xA5}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00}, []byte{0x52, 0xE9, 0x31, 0xDF, 0x24, 0x8C, 0xE4, 0xC7}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00}, []byte{0x4B, 0xB1, 0x65, 0xFD, 0xB3, 0xBF, 0xF6, 0x5C}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00}, []byte{0x7C, 0xFA, 0xFA, 0x68, 0x61, 0xD7, 0xB4, 0x7D}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00}, []byte{0x48, 0xD1, 0x75, 0x52, 0x31, 0xF8, 0x7A, 0x2A}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00}, []byte{0x41, 0x32, 0x07, 0xDA, 0x1C, 0x9B, 0x6A, 0xB5}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00}, []byte{0x63, 0xF8, 0x18, 0xE9, 0x38, 0x2A, 0x27, 0x78}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00}, []byte{0xED, 0xAF, 0x2B, 0x85, 0xFC, 0x30, 0xEB, 0x09}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00}, []byte{0x11, 0xFC, 0x59, 0x93, 0x82, 0x07, 0x63, 0xF7}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00}, []byte{0xE5, 0x39, 0xC3, 0x96, 0x99, 0x15, 0x09, 0x2F}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00}, []byte{0x50, 0x6F, 0x6A, 0x1E, 0x83, 0x4A, 0xD8, 0xF7}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00}, []byte{0x8B, 0x15, 0xBA, 0x30, 0x47, 0xFA, 0x31, 0x95}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x00}, []byte{0x13, 0x0B, 0xE1, 0x5C, 0x39, 0x3E, 0x4B, 0x7A}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00}, []byte{0x88, 0x95, 0xEC, 0x31, 0x04, 0xCA, 0x10, 0x41}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00}, []byte{0xE4, 0x40, 0xAC, 0xDF, 0x4B, 0x64, 0xC9, 0xC9}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00}, []byte{0xC2, 0x32, 0x80, 0xEB, 0xE0, 0x93, 0xF0, 0x02}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}, []byte{0x52, 0x64, 0xA6, 0x57, 0x41, 0xFE, 0x78, 0xE3}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40}, []byte{0x80, 0x89, 0x2E, 0x76, 0x85, 0x47, 0xCE, 0x61}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20}, []byte{0x09, 0x11, 0x41, 0x2D, 0x72, 0x09, 0x34, 0x75}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10}, []byte{0x9F, 0x21, 0xAA, 0x76, 0x47, 0x83, 0xE6, 0x49}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08}, []byte{0x4C, 0xA9, 0xFA, 0xBE, 0xAD, 0x2C, 0x02, 0xC6}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04}, []byte{0x59, 0xCE, 0x10, 0x97, 0x3A, 0x7B, 0x1F, 0xD5}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02}, []byte{0x68, 0x3B, 0x29, 0x34, 0xE0, 0xCC, 0xBE, 0xAA}},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, []byte{0x74, 0xD0, 0xE7, 0xC2, 0xE3, 0xB4, 0x50, 0xA8}},
}

func reverse(b []byte) []byte {

	r := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		r[len(b)-1-i] = b[i]
	}

	return r
}

// make sure we can encrypt to produce our test vectors, and decrypt to produce the original plaintext.
func TestSkipjackEncrypt(t *testing.T) {

	for _, v := range skipjackTestVectors {
		h, _ := New(v.key)

		var c, p [8]byte

		h.Encrypt(c[:], v.plain)

		if bytes.Compare(v.cipher, c[:]) != 0 {
			t.Errorf("skipjack encrypt failed: got %#v wanted %#v\n", c, v.cipher)
		}

		h.Decrypt(p[:], c[:])

		if bytes.Compare(v.plain, p[:]) != 0 {
			t.Errorf("skipjack decrypt failed: got %#v wanted %#v\n", p, v.plain)
		}
	}

	// NOTE: In the SKIPJACK standard, vectors are presented MSB, while in
	// the validation document they're presented LSB.  Rather than
	// rewriting the vectors, I'm just calling reverse() before use.

	// validation vectors all encrypted with the 0 key
	for _, v := range skipjackVariablePlaintextValidation {

		var z [10]byte
		h, _ := New(z[:])

		var c, p [8]byte

		h.Encrypt(c[:], reverse((v.plain)))

		if bytes.Compare(reverse(v.cipher), c[:]) != 0 {
			t.Errorf("skipjack validation variable plaintext encrypt failed: got %#v wanted %#v\n", c, v.cipher)
		}

		h.Decrypt(p[:], c[:])

		if bytes.Compare(reverse(v.plain[:]), p[:]) != 0 {
			t.Errorf("skipjack validation variable plaintext decrypt failed: got %#v wanted %#v \n", p, v.plain)
		}
	}

	// validation vectors all encrypt the zero block
	for _, v := range skipjackVariableKeyValidation {

		h, _ := New(reverse(v.key))

		var c, p [8]byte

		var z [8]byte
		h.Encrypt(c[:], z[:])

		if bytes.Compare(reverse(v.cipher), c[:]) != 0 {
			t.Errorf("skipjack validation variable key failed: got %#v wanted %#v\n", c, v.cipher)
		}

		h.Decrypt(p[:], c[:])

		if bytes.Compare(z[:], p[:]) != 0 {
			t.Errorf("skipjack validation decrypt failed: got %#v wanted all zero\n", p)
		}

	}
}
