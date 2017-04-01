package main

import (
	"gopkg.in/ini.v1"
)

type config struct {
	hosts map[string]string
}

func newConfig() *config {
	return &config{hosts: make(map[string]string)}
}

func (c *config) load(b []byte) error {
	f, err := ini.Load(b)
	if err != nil {
		return err
	}
	for _, n := range f.SectionStrings() {
		s, err := f.GetSection(n)
		if err != nil {
			return err
		}
		for _, kn := range s.KeyStrings() {
			c.hosts[kn] = n
		}
	}
	return nil
}

func (c *config) Get(host string) (string, bool) {
	v, ok := c.hosts[host]
	return v, ok
}
