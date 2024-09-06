package client

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/martin.kluth1/fserve/internal/config"
	"gitlab.com/martin.kluth1/fserve/internal/logging"
)

var (
	cfg config.Config
	opt getSignatureOptions
)

func init() {
	var err error

	cfg, err = config.ReadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
	}

	if err := logging.Setup(&cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to init logging")
	}

	ClientCmd.AddCommand(createSignatureCmd)
	addSignatureFlags(createSignatureCmd, &opt)
}
