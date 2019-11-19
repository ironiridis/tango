package rpc

// These are not real UUIDs, but are close enough. Notably they are not compliant with
// RFC 4122 because they do not correctly encode their version or variant bits.

import (
	"bytes"
	"crypto/rand"
	"fmt"
)

// A UUID is a pseudo-unique identifier, very similar to the ones described in RFC 4122.
type UUID [16]byte

// NewUUID returns a securely generated UUID, ready to use. NewUUID will panic if it is
// unable to read enough random data to generate an identifier.
func NewUUID() *UUID {
	var u UUID
	_, err := rand.Read(u[:])
	if err != nil {
		//
		panic(fmt.Errorf("failed to generate random UUID: %w", err))
	}

	return &u
}

// EqualTo returns true if the UUID c is bit-for-bit identical to the UUID u.
func (u *UUID) EqualTo(c *UUID) bool {
	return (bytes.Equal(u[:], c[:]))
}

func (u *UUID) String() string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// MarshalText implements encoding.TextMarshaler.
func (u *UUID) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *UUID) UnmarshalText(b []byte) error {
	var t [32]byte
	var i int
	for _, v := range b {
		switch v {
		case 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39:
			t[i] = v - 0x30
			i++
		case 0x41, 0x42, 0x43, 0x44, 0x45, 0x46:
			t[i] = v - 0x37
			i++
		case 0x61, 0x62, 0x63, 0x64, 0x65, 0x66:
			t[i] = v - 0x57
			i++
		}
		if i == 32 {
			break
		}
	}
	if i != 32 {
		return fmt.Errorf("need 32 hex bytes; read %d", i)
	}
	for i = range u[:] {
		u[i] = t[2*i]<<4 | t[2*i+1]
	}
	return nil
}
