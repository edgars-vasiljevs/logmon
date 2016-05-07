package main

import (
	"fmt"
)

func Print(msg interface{}) {
	fmt.Printf("[logmon] %s\n", msg)
}