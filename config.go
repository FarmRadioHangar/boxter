package main

import (
	"encoding/json"
	"log"

	"io/ioutil"

	"path/filepath"

	"sort"

	"github.com/blang/semver"
	"gopkg.in/ini.v1"
)

type config struct {
	SerialDir         string `json:"boxidDir"`
	HostFile          string `json:"hostsFile"`
	LocalPlaybookDir  string `json:"localPlaybookDIr"`
	RemotePlaybookDir string `json:"remotePlaybookDIr"`
	hosts             map[string]hostProp
	plays             playList
}

type hostProp struct {
	Name    string
	Serial  string
	version string
}

type playList []string

func (p playList) Len() int      { return len(p) }
func (p playList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p playList) Less(i, j int) bool {
	vi, err := semver.Make(p[i])
	if err != nil {
		log.Fatal(err)
	}
	vj, err := semver.Make(p[j])
	if err != nil {
		log.Fatal(err)
	}
	return vi.LT(vj)
}

func newConfig(src []byte) (*config, error) {
	c := &config{hosts: make(map[string]hostProp)}
	err := json.Unmarshal(src, c)
	if err != nil {
		return nil, err
	}
	err = c.loadHosts()
	if err != nil {
		return nil, err
	}
	err = c.loadPlays()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *config) loadHosts() error {
	b, err := ioutil.ReadFile(c.HostFile)
	if err != nil {
		return err
	}
	return c.load(b)
}
func (c *config) loadPlays() error {
	dirs, err := ioutil.ReadDir(c.LocalPlaybookDir)
	if err != nil {
		return err
	}
	var p playList
	for _, dir := range dirs {
		if dir.IsDir() {
			b := filepath.Base(dir.Name())
			_, err = semver.Make(b)
			if err != nil {
				return err
			}
			p = append(p, b)
		}
	}
	sort.Sort(p)
	c.plays = p
	return nil

}

func (c *config) hasPlay(ver string) bool {
	for _, v := range c.plays {
		if v == ver {
			return true
		}
	}
	if ver == "latest" && len(c.plays) > 0 {
		return true
	}
	return false
}
func (c *config) latestPlay() string {
	return c.plays[len(c.plays)-1]
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
			c.hosts[kn] = hostProp{Name: kn, version: n, Serial: s.Key(kn).String()}
		}
	}
	return nil
}

func (c *config) Get(host string) (hostProp, bool) {
	v, ok := c.hosts[host]
	return v, ok
}
