package worker

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/zkMeLabs/mechain-go-sdk/client"
	"github.com/zkMeLabs/mechain-go-sdk/types"
	"mechainbench/app/config"
	"sync"
	"time"
)

// Worker is the interface of worker node
type Worker interface {
	// Run call the workers to running.
	Run(ctx context.Context, nonce *sync.Mutex) error
}

type DefaultWorker struct {
	Id          uint
	config      *config.Config
	chainClient client.IClient
}

func NewDefaultWorker(id uint, config *config.Config) (Worker, error) {
	// new chain client
	c, err := client.New(config.Chain.ChainId, config.Chain.RpcAddr, client.Option{})
	if err != nil {
		return nil, err
	}

	if account, err := types.NewAccountFromPrivateKey("", config.Option.DefaultPrivateKey); err != nil {
		return nil, fmt.Errorf("invalid private key: %s", err)
	} else {
		c.SetDefaultAccount(account)
	}

	return &DefaultWorker{
		Id:          id,
		config:      config,
		chainClient: c,
	}, nil
}

// Run call the workers to running.
func (worker *DefaultWorker) Run(ctx context.Context, nonce *sync.Mutex) error {
	log.Debug().Msgf("worker %d start running", worker.Id)
	for deadline, ok := ctx.Deadline(); ok && time.Now().Before(deadline); {
		if err := worker.Do(ctx, nonce); err != nil {
			log.Error().Msgf("worker %d execute tx failed: %s", worker.Id, err)
		}
	}
	return nil
}

// Do execute tx
func (worker *DefaultWorker) Do(ctx context.Context, nonce *sync.Mutex) error {
	cancelCtx, cancelFunc := context.WithTimeout(ctx, time.Duration(worker.config.Bench.Timeout))
	defer cancelFunc()
	return worker.do(cancelCtx, nonce)
}
