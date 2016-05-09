package main

import (
	"flag"
	"fmt"
)

type FlagsConfig struct {
	Config, Host, Port string
}

func (c *FlagsConfig) ParseFlags() {
	config := flag.String("config", "config.json", "Specify configuration file to use")
	host := flag.String("host", "127.0.0.1", "Local server address")
	port := flag.String("port", "8080", "Local server port")

	flag.Parse()

	c.Config = *config
	c.Host = *host
	c.Port = *port
}

func NewFlagsConfig() FlagsConfig {
	flags := FlagsConfig{}
	flags.ParseFlags()
	return flags
}

func Print(arguments ...interface{}) {
	fmt.Printf(fmt.Sprintf("logmon: %s\n", arguments[0]), arguments[1:]...)
}
