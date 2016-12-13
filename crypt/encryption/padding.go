package encryption

import (
	"bytes"
	"errors"
	"fmt"
)

type Padding interface {
	Pad(p []byte) ([]byte, error)
	Unpad(p []byte) ([]byte, error)
}

var (
	// ErrInvalidBlockSize indicates hash blocksize <= 0.
	ErrInvalidBlockSize = errors.New("invalid blocksize")

	// ErrInvalidPKCS7Data indicates bad input to PKCS7 pad or unpad.
	ErrInvalidPKCS7Data = errors.New("invalid PKCS7 data (empty or not padded)")

	// ErrInvalidPKCS7Padding indicates PKCS7 unpad fails to bad input.
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

type Padder struct{ blockSize int }

// NewPkcs5Padding returns a PKCS5 padding type structure. The blocksize
// defaults to 8 bytes (64-bit).
// See https://tools.ietf.org/html/rfc2898 PKCS #5: Password-Based Cryptography.
// Specification Version 2.0
func NewPkcs5Padding() Padding {
	return &Padder{blockSize: 8}
}

// NewPkcs7Padding returns a PKCS7 padding type structure. The blocksize is
// passed as a parameter.
// See https://tools.ietf.org/html/rfc2315 PKCS #7: Cryptographic Message
// Syntax Version 1.5.
// For example the block size for AES is 16 bytes (128 bits).
func NewPkcs7Padding() Padding {
	return &Padder{blockSize: 16}
}

// Pad returns the byte array passed as a parameter padded with bytes such that
// the new byte array will be an exact multiple of the expected block size.
// For example, if the expected block size is 8 bytes (e.g. PKCS #5) and that
// the initial byte array is:
// 	[]byte{0x0A, 0x0B, 0x0C, 0x0D}
// the returned array will be:
// 	[]byte{0x0A, 0x0B, 0x0C, 0x0D, 0x04, 0x04, 0x04, 0x04}
// The value of each octet of the padding is the size of the padding. If the
// array passed as a parameter is already an exact multiple of the block size,
// the original array will be padded with a full block.
func (p *Padder) Pad(buf []byte) ([]byte, error) {
	bufLen := len(buf)
	padLen := p.blockSize - (bufLen % p.blockSize)
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(buf, padText...), nil
}

// Unpad removes the padding of a given byte array, according to the same rules
// as described in the Pad function. For example if the byte array passed as a
// parameter is:
// 	[]byte{0x0A, 0x0B, 0x0C, 0x0D, 0x04, 0x04, 0x04, 0x04}
// the returned array will be:
// 	[]byte{0x0A, 0x0B, 0x0C, 0x0D}
func (p *Padder) Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("Unpad Data is empty.")
	}
	if length%p.blockSize != 0 {
		return nil, errors.New("Unpad Data is not block-aligned.")
	}

	c := data[len(data)-1]
	n := int(c)
	if n == 0 || n > len(data) {
		return nil, fmt.Errorf("Invalid padding length : %d=%d", n, len(data))
	}

	for i := 0; i < n; i++ {
		if data[len(data)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return data[:len(data)-n], nil
}
