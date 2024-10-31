package mechain

import "cosmossdk.io/math"

// TransferDetail includes the target address and amount for token transfer.
type TransferDetail struct {
	ToAddress string
	Amount    math.Int
}
