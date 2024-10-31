package worker

import (
	"context"
	"cosmossdk.io/math"
	"fmt"
	evmostypes "github.com/evmos/evmos/v12/sdk/types"
	"github.com/rs/zerolog/log"
	"sync"
)

func (worker *DefaultWorker) do(ctx context.Context, nonce *sync.Mutex) error {
	log.Debug().Msg("----- do work")

	nonce.Lock()
	txHash, err := worker.chainClient.Transfer(ctx, "0x32a91324730D77FC25cfFF5a21038f306b6a8a30", math.NewIntWithDecimal(100, 9), evmostypes.TxOption{})
	nonce.Unlock()
	if err != nil {
		return err
	}
	log.Debug().Msgf("tx send success, hash: %s", txHash)

	// check result
	txResp, err := worker.chainClient.WaitForTx(ctx, txHash)
	if err != nil {
		return fmt.Errorf("wait tx %s failed: %s", txHash, err)
	}
	if txResp.TxResult.Code != 0 {
		return fmt.Errorf("tx %s has failed with: %s", txHash, txResp.TxResult.Info)
	}
	log.Debug().Msgf("tx %s has success", txHash)
	return nil
}
