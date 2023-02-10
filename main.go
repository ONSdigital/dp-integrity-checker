package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"

	"github.com/ONSdigital/dp-integrity-checker/checker"
	"github.com/ONSdigital/dp-integrity-checker/config"
)

const serviceName = "dp-integrity-checker"

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(ctx, "fatal runtime error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return errors.Wrap(err, "unable to retrieve service configuration")
	}
	log.Info(ctx, "config on startup", log.Data{"config": cfg, "build_time": BuildTime, "git-commit": GitCommit})

	chk := checker.Checker{
		ZebedeeRoot:                cfg.ZebedeeRoot,
		CheckPublishedPreviousDays: cfg.CheckPublishedPreviousDays,
	}

	// Run the checker in the background, using a result channel and an error channel for fatal errors
	errChan := make(chan error, 1)
	resultChan := make(chan *checker.Result, 1)
	go func() {
		result, err := chk.Run(ctx)
		if err != nil {
			errChan <- err
		}
		resultChan <- result
	}()

	// blocks until completion, an os interrupt or a fatal error occurs
	select {
	case err := <-errChan:
		log.Error(ctx, "checker error received", err)
	case sig := <-signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	case result := <-resultChan:
		log.Info(ctx, "integrity check complete", log.Data{"Result": result})
	}
	return nil // TODO close down the checker and confirm task completion state (err or nil)
}
