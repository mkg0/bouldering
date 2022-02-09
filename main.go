package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mkg0/bouldering/internal/persist"
)

func main() {
	persist.Load(&global)
	go shut()

	runCli()
}

func shut() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(1)
}
