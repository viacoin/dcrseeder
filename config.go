// Copyright (c) 2018 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/viacoin/viad/chaincfg"
	"github.com/viacoin/viautil"
	flags "github.com/jessevdk/go-flags"
)

const (
	defaultConfigFilename = "viaseeder.conf"
)

var (
	// Default network parameters
	activeNetParams = &chaincfg.MainNetParams

	// Default configuration options
	defaultConfigFile = filepath.Join(defaultHomeDir, defaultConfigFilename)
	defaultHomeDir    = viautil.AppDataDir("viaseeder", false)
)

// config defines the configuration options for hardforkdemo.
//
// See loadConfig for details on the configuration load process.
type config struct {
	Host       string `short:"H" long:"host" description:"hostname"`
	Nameserver string `short:"n" long:"nameserver" description:"hostname of nameserver"`
	Seeder     string `short:"s" long:"default seeder" description:"IP address of a  working node"`
	TestNet    bool   `long:"testnet" description:"Use the test network"`
}

func loadConfig() (*config, error) {
	err := os.MkdirAll(defaultHomeDir, 0700)
	if err != nil {
		// Show a nicer error message if it's because a symlink is
		// linked to a directory that does not exist (probably because
		// it's not mounted).
		if e, ok := err.(*os.PathError); ok && os.IsExist(err) {
			if link, lerr := os.Readlink(e.Path); lerr == nil {
				str := "is symlink %s -> %s mounted?"
				err = fmt.Errorf(str, e.Path, link)
			}
		}

		str := "failed to create home directory: %v"
		err := fmt.Errorf(str, err)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	// Default config.
	cfg := config{}

	preCfg := cfg
	preParser := flags.NewParser(&preCfg, flags.Default)
	_, err = preParser.Parse()
	if err != nil {
		e, ok := err.(*flags.Error)
		if ok && e.Type == flags.ErrHelp {
			os.Exit(0)
		}
		preParser.WriteHelp(os.Stderr)
		return nil, err
	}

	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	usageMessage := fmt.Sprintf("Use %s -h to show usage", appName)

	// Load additional config from file.
	parser := flags.NewParser(&cfg, flags.Default)
	err = flags.NewIniParser(parser).ParseFile(defaultConfigFile)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			fmt.Fprintf(os.Stderr, "Error parsing config "+
				"file: %v\n", err)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, err
		}
	}

	// Parse command line options again to ensure they take precedence.
	_, err = parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
		}
		return nil, err
	}

	if len(cfg.Host) == 0 {
		str := "Please specify a hostname"
		err := fmt.Errorf(str)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	if len(cfg.Nameserver) == 0 {
		str := "Please specify a nameserver"
		err := fmt.Errorf(str)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	if len(cfg.Seeder) == 0 {
		str := "Please specify a seeder"
		err := fmt.Errorf(str)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	if cfg.TestNet {
		activeNetParams = &chaincfg.TestNet3Params
	}

	return &cfg, nil
}
