package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"os"
)

const defaultConfig string = "config.json"
const defaultConfigTemplate string = "assets/config-template.json"

type Config struct {
	configFile     string
	configTemplate string
	content        [][]string
	new            bool
}

// Create default json config file from template
func (c Config) createDefaultConfig() error {
	template, err := Asset(c.configTemplate)
	if err != nil {
		return errors.New("Could not open config template file: " + err.Error())
	}
	err = ioutil.WriteFile(c.configFile, template, 0777)
	if err != nil {
		return errors.New("Could not generate config file: " + err.Error())
	}
	return nil
}

// Parses cli flag and overrides default config name
func (c *Config) parseFlags() bool {
	customConfig := flag.String("config", "", "Specify configuration file to use")
	// Parse cli flags
	flag.Parse()
	// Override default config file
	if *customConfig != "" {
		c.configFile = *customConfig
		return true
	}
	return false
}

// Load configuration from json file
func (c *Config) load() error {
	// Check if default config file is overwritten in cli flags
	if c.parseFlags() == false {
		// Check if default config file exists.
		// If not, create empty one from template
		if _, err := os.Stat(c.configFile); os.IsNotExist(err) {
			err := c.createDefaultConfig()
			if err != nil {
				return err
			}
			// Mark config as new so it's not loaded right away
			c.new = true
		}
	}

	data, err := ioutil.ReadFile(c.configFile)
	if err != nil {
		return errors.New("Could not read config file: " + err.Error())
	}

	err = json.Unmarshal(data, &c.content)
	if err != nil {
		return errors.New("Config parse error: " + err.Error())
	}

	return nil
}

// Create new configuration object
func NewConfig() (Config, error) {
	config := Config{configFile: defaultConfig, configTemplate: defaultConfigTemplate}
	return config, config.load()
}
