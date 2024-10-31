package engine

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/zkMeLabs/mechain-go-sdk/client"
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

	mutex *sync.Mutex

	defaultClient client.IClient
}

func NewDefaultEngine(config *config.Config) Engine {
	// default client
	c, err := client.New(config.Chain.ChainId, config.Chain.RpcAddr, client.Option{})
	if err != nil {
		log.Panic().Msgf("new chain client failed: %s", err)
	}

	// all workers
	var workers []worker.Worker
	for i := range config.Bench.Concurrency {
		if w, err := worker.NewDefaultWorker(i, config); err != nil {
			panic(err)
		} else {
			workers = append(workers, w)
		}
	}

	return &DefaultEngine{
		config:        config,
		workers:       workers,
		mutex:         new(sync.Mutex),
		defaultClient: c,
	}
}

func (e *DefaultEngine) Run(ctx context.Context) {

	_, err := e.defaultClient.GetAccount(ctx, "0x00000Be6819f41400225702D32d3dd23663Dd690")
	if err != nil {
		log.Fatal().Msgf("get account failed: %s", err)
		return
	}

	// 设置执行时间
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
			if err := e.workers[id].Run(ctx, e.mutex); err != nil {
				slog.Error("worker %d execute tx failed: %s", id, err)
			}
		}(int(i))
	}
	wg.Wait()
}
