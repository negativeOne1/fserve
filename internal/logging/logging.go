package logging

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/martin.kluth1/fserve/internal/config"
)

func Setup(cfg *config.Config) error {
	l, err := zerolog.ParseLevel(strings.ToLower(cfg.Log.Level))
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(l)

	if strings.ToLower(cfg.Log.Format) == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	return nil
}
