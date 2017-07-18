package main

import "time"

// Change at will.
const minApprovals int = 2
const checkInterval time.Duration = 30000

// TODO: make this var at runtime
var repositories = [...]string{"_repo"}
