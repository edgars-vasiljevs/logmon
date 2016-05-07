package main

func main() {
	// Main channel for logs
	logs := make(chan LogMessage)

	config, err := NewConfig()
	if err != nil {
		Print(err)
		return
	}

	// Exit if new config was generated
	if config.new {
		Print("Empty configuration file was generated (config.json).")
		Print("Edit the file and re-run this command.")
		return
	}

	Print("Using " + config.configFile + " configuration file.")

	// Start file monitoring
	go NewFileMonitor(config, logs)

	// Start HTTP server
	go NewHTTPServer()

	// Start a TCP socket server
	go NewWebSocketServer(config, logs)

	// Block
	select {}
}
