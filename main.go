package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func defaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".scpull.yaml"
	}
	return filepath.Join(home, ".scpull.yaml")
}

func run() error {
	configPath := flag.String("config", defaultConfigPath(), "path to the configuration file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: scpull [-config PATH] PROFILE HOST [HOST...]\n\n")
		fmt.Fprintf(os.Stderr, "Fetch the files defined in PROFILE from each HOST via scp.\n")
		fmt.Fprintf(os.Stderr, "Files are saved under ./HOST/ preserving the remote path structure.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		return fmt.Errorf("expected a profile name and at least one host")
	}

	profileName := args[0]
	hosts := args[1:]

	cfg, err := loadConfig(*configPath)
	if err != nil {
		return err
	}

	p, err := cfg.profile(profileName)
	if err != nil {
		return err
	}

	for _, host := range hosts {
		for _, remote := range p.Paths {
			if err := fetchFile(p.User, host, remote); err != nil {
				return err
			}
		}
	}

	return nil
}
