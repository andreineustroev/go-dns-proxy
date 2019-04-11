package main

import (
	"flag"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultListenAddress = "0.0.0.0"
	defaultDnsServer     = ""
	defaultListenPort    = "53"
	defaultConfig        = "./config.json"
	defaultLogLevel      = log.InfoLevel
	defaultDomains       = "{}"
)

func CreateConfig() error {
	path := flag.String("config", defaultConfig, "path to config file")
	flag.Parse()

	initDefaultConfig()

	viper.SetConfigFile(*path)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("can't read config: %v", err)
	}

	viper.WatchConfig()
	return nil
}

func initDefaultConfig() {
	viper.SetDefault("listen_host", defaultListenAddress)
	viper.SetDefault("listen_port", defaultListenPort)
	viper.SetDefault("default_dns", defaultDnsServer)
	viper.SetDefault("log_level", defaultLogLevel)
	viper.SetDefault("domains", defaultDomains)
}

func ConfigLogging() {
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "2006/01/02 15:04:05.000"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	logLevel := strings.ToLower(viper.GetString("log_level"))
	if level, err := log.ParseLevel(logLevel); err == nil {
		log.SetLevel(level)
	} else {
		log.Warnf("Incorrect config logLevel: '%v', use default '%v'", logLevel, defaultLogLevel)
		log.SetLevel(defaultLogLevel)
	}
}
