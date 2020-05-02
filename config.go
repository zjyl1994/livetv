package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-ini/ini"
)

type Config struct {
	BaseURL     string
	ListenOn    string
	ProxyStream bool
	M3UTemplate string
	ChannelFile string
	YtdlCmd     string
	YtdlArgs    string
	PreloadCron string
}

var (
	cfg            *Config
	configFilepath string
	pidFilepath    string
)

const (
	httpClientTimeout = 10 * time.Second
	versionString     = "LiveTV ALPHA Edition"
)

func initProc() error {
	var (
		h bool // help
		v bool // version
	)
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.BoolVar(&h, "help", false, "show this help")
	flag.BoolVar(&v, "version", false, "get version info")
	flag.StringVar(&configFilepath, "config", "config.ini", "specify config file path")
	flag.StringVar(&pidFilepath, "pidfile", "", "specify PID file path")
	flag.Parse()
	if h {
		flag.Usage()
		os.Exit(0)
	}
	if v {
		fmt.Println(versionString)
		os.Exit(0)
	}
	if pidFilepath != "" {
		if _, err := os.Stat(pidFilepath); os.IsNotExist(err) {
			ioutil.WriteFile(pidFilepath, []byte(strconv.Itoa(os.Getpid())), 0644)
		} else {
			log.Fatalln("The PID file of LiveTV exists and the current instance will not start.")
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
