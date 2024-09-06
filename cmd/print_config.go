package cmd

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gitlab.com/martin.kluth1/fserve/internal/config"
	"gitlab.com/martin.kluth1/fserve/internal/logging"
)

var printConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "print the configuration",
	Run:   printConfig,
}

func printConfig(cmd *cobra.Command, args []string) {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
	}

	if err := logging.Setup(&cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to init logging")
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Err(err).Msg("failed to marshal config")
	}
	log.Info().Msg(string(data))
}
