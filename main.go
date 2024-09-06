package main

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/martin.kluth1/fserve/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Err(err).Msg("")
	}
}
