package main

import (
	"fmt"
)

func Print(arguments ...interface{}) {
	fmt.Printf(fmt.Sprintf("logmon: %s\n", arguments[0]), arguments[1:]...)
}
