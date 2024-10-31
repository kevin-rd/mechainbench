package engine

import (
	"context"
	"github.com/rs/zerolog/log"
	"log/slog"
	"mechainbench/app/config"
	"mechainbench/app/worker"
	"sync"
	"time"
)

// Engine is used to control the rate for send tx.
type Engine interface {
	// Run start the engine.
	Run(ctx context.Context)

	// Close the engine.
	Close()
}

type DefaultEngine struct {
	config  *config.Config
	workers []worker.Worker
}

func NewDefaultEngine(config *config.Config) Engine {

	var workers []worker.Worker
	for i := range config.Bench.Concurrency {
		if w, err := worker.NewDefaultWorker(i, config); err != nil {
			panic(err)
		} else {
			workers = append(workers, w)
		}
	}

	return &DefaultEngine{
		config:  config,
		workers: workers,
	}
}

func (e *DefaultEngine) Run(ctx context.Context) {
	timeoutCtx, cancelFunc := context.WithTimeout(ctx, time.Duration(e.config.Bench.Duration))
	defer cancelFunc()
	e.schedule(timeoutCtx)
}

func (e *DefaultEngine) Close() {
	log.Info().Msg("engine closed")
}

func (e *DefaultEngine) schedule(ctx context.Context) {
	var wg sync.WaitGroup
	for i := uint(0); i < e.config.Bench.Concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if err := e.workers[id].Run(ctx); err != nil {
				slog.Error("worker %d execute tx failed: %s", id, err)
			}
		}(int(i))
	}
	wg.Wait()
}
