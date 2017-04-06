package main

import (
	"testing"
)

func TestUnauthoized(t *testing.T) {
	cfg := &config{hosts: make(map[string]hostProp)}
	err := cfg.load([]byte(sample))
	if err != nil {
		t.Fatal(err)
	}
	err = uauthorized(cfg, "habari", "someuniqueid")
	if err != nil {
		t.Fatal(err)
	}
}
