package main

import "testing"

var sample = `

[all]
box2 =
box3  =
box4 =

[0.1.0]
habarimaalum=

[0.2.1]
mambojambo=
`

func TestConfig(t *testing.T) {
	cfg := &config{hosts: make(map[string]string)}
	err := cfg.load([]byte(sample))
	if err != nil {
		t.Fatal(err)
	}
	vers := []struct {
		h, v string
	}{
		{"box4", "all"},
		{"habarimaalum", "0.1.0"},
		{"mambojambo", "0.2.1"},
	}
	for _, v := range vers {
		s, ok := cfg.Get(v.h)
		if !ok {
			t.Errorf("missing host %s", v.h)
			continue
		}
		if s != v.v {
			t.Errorf("expected %s got %s", v.v, s)
		}
	}
}
