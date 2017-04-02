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
	a.Version = "0.1.0"
	a.Name = "coaster"
	a.Usage = "manages syncing of playbook releases"
	a.Commands = []cli.Command{
		{
			Name:   "sync",
			Usage:  "syncs new or specified version of playbook",
			Action: rsync,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config",
					Usage: "path to the configuration file",
				},
				cli.StringFlag{
					Name:  "host",
					Usage: "the name of the host machine",
				},
				cli.StringSliceFlag{
					Name:  "ssh",
					Usage: "ssh connection url",
				},
				cli.BoolFlag{
					Name:  "force",
					Usage: "will force the operation",
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
	cFile := ctx.String("config")
	b, err := ioutil.ReadFile(cFile)
	if err != nil {
		return err
	}
	cfg, err := newConfig(b)
	if err != nil {
		return err
	}
	v, ok := cfg.Get(host)
	if !ok {
		fmt.Printf("can't find the host %sin the host file, falling back to latest\n", host)
		v = "latest"
	}
	if cfg.hasPlay(v) {
		if v == "latest" {
			v = cfg.latestPlay()
		}
		return rsync(cfg, v, ssh)
	}
	return fmt.Errorf("no play for %s found", v)
}

func rsync(cfg *config, ver, ssh string) error {
	src := filepath.Join(cfg.LocalPlaybookDir, ver)
	dest := filepath.Join(cfg.RemotePlaybookDir, ver)
	cmd := exec.Command(
		"rsync", "-vzh", "ssh", ssh, src, dest,
	)
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}
