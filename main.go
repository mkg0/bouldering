package main

import (
	"bouldering-auto-book/internal/persist"
)

func main() {
	persist.Load(&global)
	runCli()
}
