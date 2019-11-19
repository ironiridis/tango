package main

import "sync"

type ConfigProviderMemory struct {
	sync.RWMutex
	m map[string][]byte
}

func (m *ConfigProviderMemory) Init(setup ConfigProviderSetup) {
	// setup is ignored
	m.m = make(map[string][]byte)
}

func (m *ConfigProviderMemory) Have(key string) bool {
	m.RLock()
	_, ok := m.m[key]
	m.RUnlock()
	return ok
}

func (m *ConfigProviderMemory) Get(key string) *ConfigJSON {
	m.RLock()
	b, ok := m.m[key]
	m.RUnlock()
	c := NewConfigJSON()
	if ok {
		c.Copy(b)
	}
	return c
}

func (m *ConfigProviderMemory) Put(key string, c *ConfigJSON) error {
	m.Lock()
	m.m[key] = c.Bytes()
	m.Unlock()
	return nil
}

func (m *ConfigProviderMemory) Purge(key string) error {
	m.Lock()
	delete(m.m, key)
	m.Unlock()
	return nil
}

func (m *ConfigProviderMemory) Sync() error {
	return nil
}
