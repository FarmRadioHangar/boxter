package main

import (
	"io/ioutil"
	"log"
	"os"

	"fmt"

	"path/filepath"

	"os/exec"

	"github.com/urfave/cli"
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
	if !ok {
		if ver != "" {
			v = ver
		} else {
			fmt.Printf("can't find the host %sin the host file, falling back to latest\n", host)
			v = "latest"
		}
	}
	if cfg.hasPlay(v) {
		if v == "latest" {
			v = cfg.latestPlay()
		}
		return rsync(cfg, v, rsh, ssh)
	}
	return fmt.Errorf("no play for %s found", v)
}

func rsync(cfg *config, ver, rsh, ssh string) error {
	src := filepath.Join(cfg.LocalPlaybookDir, ver)
	dest := filepath.Join(cfg.RemotePlaybookDir, ver)
	cmd := exec.Command(
		"rsync", "-z", "--rsh", rsh, src, fmt.Sprintf("%s:%s", ssh, dest),
	)
	b, err := cmd.CombinedOutput()
	fmt.Println(string(b))
	if err != nil {
		return err
	}
	return nil
}
