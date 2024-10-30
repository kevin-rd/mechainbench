package work

import (
	"context"
	"log/slog"
	"mechainbench/app/core"
)

// Worker is the interface of worker node
type Worker interface {
	// Run call the workers to running.
	Run(ctx context.Context) error
}

type DefaultWorker struct {
	ctx *core.Context
}

func NewDefaultWorker(ctx *core.Context) Worker {
	return &DefaultWorker{ctx: ctx}
}

// Run call the workers to running.
func (worker *DefaultWorker) Run(ctx context.Context) error {
	for i := uint(0); i < worker.ctx.BatchNum; i++ {
		if err := worker.Do(ctx); err != nil {
			slog.Error("worker %d execute tx failed: %s", i, err)
		}
	}
	return nil
}

// Do execute tx
func (worker *DefaultWorker) Do(ctx context.Context) error {
	cancelCtx, cancelFunc := context.WithTimeout(ctx, worker.ctx.Timeout)
	defer cancelFunc()
	return worker.do(cancelCtx)
}

func (worker *DefaultWorker) do(ctx context.Context) error {
	slog.Info("----- do work")
	return nil
}
