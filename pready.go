package main

import (
	"time"
)

func run() {
	for {
		checkRepositories()
		time.Sleep(checkInterval * time.Millisecond)
	}
}
