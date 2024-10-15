package main

import (
	"os"
	"os/signal"
	"syscall"
)

func handleSysCalls(notifyChannel chan<- os.Signal) {
	signal.Notify(notifyChannel, syscall.SIGINT, syscall.SIGTERM)
}
