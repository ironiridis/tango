package rpc_test

import (
	"testing"

	"github.com/ironiridis/tango/rpc"
)

func TestUUIDNoDuplicates(t *testing.T) {
	a := rpc.NewUUID()
	b := rpc.NewUUID()

	if a.EqualTo(b) {
		t.Error("Generated two duplicate UUIDs")
	}
}

func TestUUIDZeroEqual(t *testing.T) {
	a := rpc.ZeroUUID()
	b := rpc.ZeroUUID()

	if !a.EqualTo(b) {
		t.Error("Generated zero UUIDs were not equal")
	}
}

func TestUUIDZeroMarshalRoundtrip(t *testing.T) {
	a := rpc.ZeroUUID()
	if !a.IsZero() {
		t.Errorf("a.IsZero() returned false")
	}
	buf, err := a.MarshalText()
	if err != nil {
		t.Errorf("Failed to marshal zero UUID: %v", err)
	}
	b := new(rpc.UUID)
	err = b.UnmarshalText(buf)
	if err != nil {
		t.Errorf("Failed to unmarshal zero UUID: %v", err)
	}
	if !b.IsZero() {
		t.Errorf("b.IsZero() returned false")
	}
}

func TestUUIDMarshalRoundtrip(t *testing.T) {
	a := rpc.NewUUID()
	buf, err := a.MarshalText()
	if err != nil {
		t.Errorf("Failed to marshal UUID: %v", err)
	}
	b := new(rpc.UUID)
	err = b.UnmarshalText(buf)
	if err != nil {
		t.Errorf("Failed to unmarshal UUID: %v", err)
	}
	if !a.EqualTo(b) {
		t.Errorf("Roundtrip did not produce an equal UUID (%x -> %q -> %x)", a, buf, b)
	}
}

func TestUUIDCommutativeEquality(t *testing.T) {
	var err error
	subjects := []string{
		"f5138183-990e-421b-8363-7428d0bf5436",
		"F7B9F29D-E9C9-40FF-B97C-C00897471078",
		"d6d4177c-3943-4c72-b6cb-7a214818dc55",
		"025d4507-577f-430d-848e-20183c207410",
		"00000000-0000-0000-0000-000000000000",
	}
	x := new(rpc.UUID)
	err = x.UnmarshalText([]byte("12341234-1234-1234-1234-123412341234"))
	if err != nil {
		t.Errorf("Failed to unmarshal contrived UUID: %v", err)
	}

	for i := range subjects {
		a := new(rpc.UUID)
		err = a.UnmarshalText([]byte(subjects[i]))
		if err != nil {
			t.Errorf("Failed to unmarshal %q: %v", subjects[i], err)
		}
		b := new(rpc.UUID)
		err = b.UnmarshalText([]byte(subjects[i]))
		if err != nil {
			t.Errorf("Failed to unmarshal %q: %v", subjects[i], err)
		}
		if !a.EqualTo(b) {
			t.Errorf("a.EqualTo(b) returned false, but both unmarshalled from %q", subjects[i])
		}
		if !b.EqualTo(a) {
			t.Errorf("b.EqualTo(a) returned false, but both unmarshalled from %q", subjects[i])
		}
		if x.EqualTo(a) {
			t.Errorf("x.EqualTo(a) returned true")
		}
		if a.EqualTo(x) {
			t.Errorf("a.EqualTo(x) returned true")
		}
		if x.EqualTo(b) {
			t.Errorf("x.EqualTo(b) returned true")
		}
		if b.EqualTo(x) {
			t.Errorf("b.EqualTo(x) returned true")
		}
	}
}

func TestUUIDUnmarshalTolerance(t *testing.T) {
	var err error
	subjects := []string{
		"12341234abcdabcd1234123412341234",
		"12341234:ABCD:ABCD:1234:123412341234",
		"12341234 aBcD AbCd 1234 123412341234",
		"{12341234-ABcd-abCD-1234-123412341234}",
		"1 2 3 4 1 2 3 4 A B C D a b c d 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4",
	}
	x := new(rpc.UUID)
	err = x.UnmarshalText([]byte("12341234-abcd-ABCD-1234-123412341234"))
	if err != nil {
		t.Errorf("Failed to unmarshal contrived UUID: %v", err)
	}

	for i := range subjects {
		a := new(rpc.UUID)
		err = a.UnmarshalText([]byte(subjects[i]))
		if err != nil {
			t.Errorf("Failed to unmarshal %q: %v", subjects[i], err)
		}
		if !x.EqualTo(a) {
			t.Errorf("x.EqualTo(a) returned false")
		}
		if !a.EqualTo(x) {
			t.Errorf("a.EqualTo(x) returned false")
		}
	}
}

func TestUUIDUnmarshalIntolerance(t *testing.T) {
	var err error
	subjects := []string{
		"",
		"12341234-1234-1234-1234-12341234123",
	}

	for i := range subjects {
		a := new(rpc.UUID)
		err = a.UnmarshalText([]byte(subjects[i]))
		if err == nil {
			t.Errorf("Unmarshal of %q succeded unexpectedly", subjects[i])
		}
	}
}
