package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/go-ini/ini"
)

type Config struct {
	BaseURL     string
	ListenOn    string
	ProxyStream bool
	M3UTemplate string
	ChannelFile string
	YtdlCmd     string
}

var cfg *Config
var configFilepath string
var pidFilepath string

func initProc() error {
	var (
		h bool // help
		v bool // version
	)
	flag.BoolVar(&h, "help", false, "show this help")
	flag.BoolVar(&v, "version", false, "get version info")
	flag.StringVar(&configFilepath, "config", "config.ini", "config file")
	flag.StringVar(&pidFilepath, "pidfile", "", "PID file")
	flag.Parse()
	if h {
		flag.Usage()
		os.Exit(0)
	}
	if v {
		fmt.Println("LiveTV ALPHA Edition")
		os.Exit(0)
	}
	if pidFilepath != "" {
		if _, err := os.Stat(pidFilepath); os.IsNotExist(err) {
			ioutil.WriteFile(pidFilepath, []byte(strconv.Itoa(os.Getpid())), 0644)
		} else {
			fmt.Println("The PID file of LiveTV exists and the current instance will not start.")
			os.Exit(0)
		}
	}
	cfg = new(Config)
	return ini.MapTo(cfg, configFilepath)
}

func removePidFile() error {
	if pidFilepath != "" {
		return os.RemoveAll(pidFilepath)
	} else {
		return nil
	}
}
