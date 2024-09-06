package server

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/martin.kluth1/fserve/internal/config"
	"gitlab.com/martin.kluth1/fserve/internal/logging"
)

func pre_run() {
	var err error

	cfg, err = config.ReadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
	}

	if err := logging.Setup(&cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to init logging")
	}
}
