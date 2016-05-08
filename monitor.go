package main

import (
	"bufio"
	"fmt"
	"github.com/hpcloud/tail"
	"golang.org/x/crypto/ssh"
	"os"
	"regexp"
	"strings"
)

type LogMessage [2]string

func LocalFileMonitor(item []string, logs chan<- LogMessage) {

	var err error

	_, err = os.Stat(item[1])
	if os.IsNotExist(err) {
		Print("File " + item[1] + " does not exist")
		return
	}

	file, err := tail.TailFile(item[1], tail.Config{
		Follow: true,
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

	sshSplit := regexp.MustCompile(`([^@]+)@([^:]+(?::\d+|)):(.+)$`).FindStringSubmatch(item[1])

	auth := strings.Split(sshSplit[1], ":")
	sshConfig := ssh.ClientConfig{}

	if len(auth) == 2 {
		sshConfig.User = auth[0]
		sshConfig.Auth = []ssh.AuthMethod{ssh.Password(auth[1])}
	} else {
		// TODO: pub key
		sshConfig.User = auth[0]
	}

	address := strings.Split(sshSplit[2], ":")

	// Add port if not defined
	if len(address) == 1 {
		address = append(address, "22")
	}

	// Connect to SSH server
	connection, err := ssh.Dial("tcp", strings.Join(address, ":"), &sshConfig)
	if err != nil {
		Print(fmt.Sprintf("Failed to dial SSH: %s", err))
		return
	}

	// Create new session
	session, err := connection.NewSession()
	if err != nil {
		Print(fmt.Sprintf("Failed to create SSH session: %s", err))
		return
	}

	// Close session once done
	defer session.Close()

	// Monitor ssh stdout
	pipe, _ := session.StdoutPipe()
	scanner := bufio.NewScanner(pipe)

	go func() {
		for scanner.Scan() {
			logs <- LogMessage{item[0], scanner.Text()}
		}
	}()

	// tail the file
	session.Run("tail -f " + sshSplit[3])
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
