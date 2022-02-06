package main

import (
	"github.com/mkg0/bouldering-auto-book/internal/persist"
)

func main() {
	persist.Load(&global)
	runCli()
}
