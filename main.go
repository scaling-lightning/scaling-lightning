package main

import (
	"github.com/rs/zerolog"
	"github.com/scaling-lightning/scaling-lightning/cmd/scalinglightning"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	scalinglightning.Execute()
}
