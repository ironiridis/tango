package main

import "os"

type ConfigProviderSetup map[string]string

// ConfigProvider describes a configuration storage backend
type ConfigProvider interface {
	// Init is called prior to any other operations.
	Init(setup ConfigProviderSetup)

	// Have returns true if `key` exists in this backend.
	Have(key string) bool

	// Get returns a ConfigJSON for `key`. If it didn't exist, the ConfigJSON is empty but ready to use.
	Get(key string) *ConfigJSON

	// Put writes the ConfigJSON to the storage backend as `key`, returning an error on failure.
	Put(key string, c *ConfigJSON) error

	// Purge removes any stored configuration for `key`, returning an error on failure.
	Purge(key string) error

	// Sync prompts the storage backend to flush pending writes, returning an error on failure.
	Sync() error
}

func getConfigProviderByName(n string) (p ConfigProvider) {

	switch n {
	case "mem":
		p = &ConfigProviderMemory{}
		p.Init(ConfigProviderSetup{})
	}
	return
}

func senseConfigProviderArgument() ConfigProvider {
	// read arguments passed to the binary, looking for -conf=
	return nil
}

func senseConfigProviderEnvironment() ConfigProvider {
	cpenv, ok := os.LookupEnv("TANGOCONF")
	if ok {
		return getConfigProviderByName(cpenv)
	}
	return nil
}

func senseConfigProviderBoot() ConfigProvider {
	// detect if we're on Linux, if /proc is mounted, and if /proc/cmdline contains `tangoconf=`
	return nil
}

func senseConfigProviderFS() ConfigProvider {
	// detect a sentinel file in our current working directory or (if Linux) in /
	return nil
}

// DefaultConfigProvider returns a heuristically determined ConfigProvider.
func DefaultConfigProvider() ConfigProvider {
	var p ConfigProvider
	p = senseConfigProviderArgument()
	if p == nil {
		p = senseConfigProviderEnvironment()
	}
	if p == nil {
		p = senseConfigProviderBoot()
	}
	if p == nil {
		p = senseConfigProviderFS()
	}
	if p == nil {
		emitWarning("Unable to determine a configuration backend, using an in-memory store. All changes will be lost.")
		p = getConfigProviderByName("mem")
	}
	return p
}
