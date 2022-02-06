package main

import (
	"github.com/mkg0/bouldering/internal/persist"
)

func main() {
	persist.Load(&global)
	runCli()
}
