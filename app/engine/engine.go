package engine

import (
	"context"
	"log/slog"
	"mechainbench/app/core"
	"mechainbench/app/work"
	"sync"
)

// Engine is used to control the rate for send tx.
type Engine interface {
	// Run start the engine.
	Run(ctx context.Context)

	// Close the engine.
	Close()
}

type DefaultEngine struct {

	// 并发量
	concurrency uint

	worker work.Worker
}

func NewDefaultEngine(ctx *core.Context) Engine {
	worker := work.NewDefaultWorker(ctx)
	return &DefaultEngine{
		concurrency: 10,
		worker:      worker,
	}
}

func (e *DefaultEngine) Run(ctx context.Context) {
	e.schedule(ctx)
}

func (e *DefaultEngine) Close() {

}

func (e *DefaultEngine) schedule(ctx context.Context) {
	var wg sync.WaitGroup
	for i := uint(0); i < e.concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if err := e.worker.Run(ctx); err != nil {
				slog.Error("worker %d execute tx failed: %s", id, err)
			}
		}(int(i))
	}
	wg.Wait()
}
