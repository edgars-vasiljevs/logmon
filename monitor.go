package main

import (
	"github.com/hpcloud/tail"
	"os"
	"regexp"
)

type LogMessage struct {
	Name    string `json:"n"`
	Content string `json:"c"`
}

// Local file monitoring
func LocalFileMonitor(item []string, logs chan<- LogMessage) {

	var err error

	_, err = os.Stat(item[1])
	if os.IsNotExist(err) {
		Print("File " + item[1] + " does not exist")
		return
	}

	file, err := tail.TailFile(item[1], tail.Config{
		Follow: true,
		//Location: &tail.SeekInfo{0, 2},
		Logger: tail.DiscardingLogger,
	})

	if err != nil {
		Print("Could not open: " + item[1] + ". " + err.Error())
		return
	}

	for line := range file.Lines {
		logs <- LogMessage{item[0], line.Text}
	}
}

func RemoteFileMonitor(item []string, logs chan<- LogMessage) {
	// root:password@host:232:/filename
}

func NewFileMonitor(config Config, logs chan<- LogMessage) {
	for _, item := range config.content {
		// Check if SSH connection
		if ok, _ := regexp.MatchString("@", item[1]); ok {
			go RemoteFileMonitor(item, logs)
		} else {
			go LocalFileMonitor(item, logs)
		}
	}
}
