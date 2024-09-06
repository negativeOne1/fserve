package server

import (
	"context"
	net "net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gitlab.com/martin.kluth1/fserve/api/http"
	"gitlab.com/martin.kluth1/fserve/internal/config"
	"gitlab.com/martin.kluth1/fserve/internal/logging"
	"gitlab.com/martin.kluth1/fserve/signature"
	"gitlab.com/martin.kluth1/fserve/storage"
)

var (
	RunCmd = &cobra.Command{
		Use:    "run",
		Short:  "run",
		Run:    run,
		PreRun: pre_run,
	}
	cfg config.Config
)

func pre_run(cmd *cobra.Command, args []string) {
	var err error

	cfg, err = config.ReadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
	}

	if err := logging.Setup(&cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to init logging")
	}
}

func run(cmd *cobra.Command, args []string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-c
		cancel()
	}()

	fs := storage.NewFileStorage(cfg.Storage.BasePath)

	hmacValidator := signature.NewHMACValidator(cfg.Secret)

	router := http.NewRouter(fs, hmacValidator)
	v1 := net.NewServeMux()
	v1.Handle("/v1/", net.StripPrefix("/v1", router))

	server := http.NewHTTPServer("0.0.0.0", cfg.HTTP.Port, v1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Run(ctx); err != nil {
			log.Fatal().Err(err).Msg("Failed to run server")
		}
	}()
	wg.Wait()
}
