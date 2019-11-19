package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
)

// ConfigJSON represents a configuration node stored as JSON.
type ConfigJSON struct {
	sync.RWMutex
	j []byte
}

// NewConfigJSON returns an empty ConfigJSON.
func NewConfigJSON() *ConfigJSON {
	return &ConfigJSON{j: make([]byte, 0)}
}

// Copy copies from the provided buffer after validating it with json.Valid(b).
func (c *ConfigJSON) Copy(b []byte) error {
	if !json.Valid(b) {
		return fmt.Errorf("invalid json")
	}
	c.Lock()
	c.j = make([]byte, len(b))
	copy(c.j, b)
	c.Unlock()
	return nil
}

// Bytes copies to a new byte slice.
func (c *ConfigJSON) Bytes() []byte {
	c.RLock()
	b := make([]byte, len(c.j))
	copy(b, c.j)
	c.RUnlock()
	return b
}

// DecodeInto unmarshals the stored JSON into `v`.
func (c *ConfigJSON) DecodeInto(v interface{}) error {
	c.RLock()
	defer c.RUnlock()
	return json.Unmarshal(c.j, v)
}

// EncodeFrom marshals `v` into JSON and stores it.
func (c *ConfigJSON) EncodeFrom(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	c.Lock()
	c.j = b // encoding/json/encode.go shows this buffer is up for grabs
	c.Unlock()
	return nil
}

func (c *ConfigJSON) String() string {
	r := new(bytes.Buffer)
	json.Indent(r, c.j, "", "  ")
	return r.String()
}
