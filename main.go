package main

import (
	"io/ioutil"
	"log"
	"os"

	"fmt"

	"path/filepath"

	"os/exec"

	"strings"

	"github.com/urfave/cli"
	"gopkg.in/ini.v1"
)

func main() {
	a := cli.NewApp()
	a.Version = "0.1.1"
	a.Name = "boxter"
	a.Usage = "manages syncing of playbook releases"
	a.Commands = []cli.Command{
		{
			Name:   "sync",
			Usage:  "syncs new or specified version of playbook",
			Action: sync,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config",
					Usage: "path to the configuration file",
				},
				cli.StringFlag{
					Name:  "host",
					Usage: "the name of the host machine",
				},
				cli.StringFlag{
					Name:  "boxid",
					Usage: "unique id of the box",
				},
				cli.StringFlag{
					Name:  "remote-playbook-dir",
					Usage: "the directory to sync playbooks to in a remote host",
				},
				cli.StringFlag{
					Name:  "ssh",
					Usage: "ssh connection url",
				},
				cli.StringFlag{
					Name:  "rsh",
					Usage: "passed to rsync",
				},
			},
		},
	}
	err := a.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func sync(ctx *cli.Context) error {
	host := ctx.String("host")
	serial := ctx.String("boxid")
	ssh := ctx.String("ssh")
	rsh := ctx.String("rsh")
	cFile := ctx.String("config")
	rDir := ctx.String("remote-playbook-dir")
	ver := ctx.Args().First()
	b, err := ioutil.ReadFile(cFile)
	if err != nil {
		return err
	}
	cfg, err := newConfig(b)
	if err != nil {
		return err
	}
	if rDir != "" {
		cfg.RemotePlaybookDir = rDir
	}
	v, ok := cfg.Get(host)
	if !ok || strings.TrimSpace(v.Serial) != strings.TrimSpace(serial) {
		return uauthorized(cfg, host, serial)
	}
	if ver != "" {
		v.version = ver
	}
	if cfg.hasPlay(v.version) {
		if v.version == "latest" {
			v.version = cfg.latestPlay()
		}
		return rsync(cfg, v, rsh, ssh)
	}
	return fmt.Errorf("no play for %s found", v)
}

func uauthorized(cfg *config, host, id string) error {
	f := filepath.Join(cfg.SerialDir, "unauthorized.ini")
	u := "unauthorized"
	fi, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer fi.Close()
	b, err := ioutil.ReadAll(fi)
	if err != nil {
		return err
	}
	var o *ini.File
	if len(b) > 0 {
		o, err = ini.Load(f)
		if err != nil {
			return err
		}
		s := o.Section(u)
		s.NewKey(host, id)

	} else {
		o = ini.Empty()
		s, err := o.NewSection(u)
		if err != nil {
			return err
		}
		s.NewKey(host, id)
	}
	_, err = o.WriteTo(fi)
	return err
}

func rsync(cfg *config, ver hostProp, rsh, ssh string) error {
	src := filepath.Join(cfg.LocalPlaybookDir, ver.version)
	dest := filepath.Join(cfg.RemotePlaybookDir, "playbook")
	fmt.Printf("syncing %s\n", ver)
	cmd := exec.Command(
		"rsync", "-z", "--rsh", rsh, src, fmt.Sprintf("%s:%s", ssh, dest),
	)
	b, err := cmd.CombinedOutput()
	fmt.Println(string(b))
	if err != nil {
		return err
	}
	fmt.Println("OK")
	fmt.Println("syncing mainfest")
	o := filepath.Join(cfg.SerialDir, ver.Serial)
	manifestFile := "voxbox-manifest.json"
	os.MkdirAll(o, 0755)
	hm := fmt.Sprintf("%s:%s", ssh, filepath.Join(cfg.RemotePlaybookDir, manifestFile))
	sm := filepath.Join(o, manifestFile)
	cmd = exec.Command(
		"rsync", "-z", "--rsh", rsh, hm, sm,
	)
	b, err = cmd.CombinedOutput()
	fmt.Println(string(b))
	if err != nil {
		return err
	}
	fmt.Println("OK")
	return nil
}
