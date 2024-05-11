package types

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
)

type SendBundleArgs struct {
	Txs               []hexutil.Bytes `json:"txs"`
	MaxBlockNumber    uint64          `json:"maxBlockNumber"`
	MinTimestamp      *uint64         `json:"minTimestamp"`
	MaxTimestamp      *uint64         `json:"maxTimestamp"`
	RevertingTxHashes []common.Hash   `json:"revertingTxHashes"`
}

func SendBundle(ec *ethclient.Client, ctx context.Context, args *SendBundleArgs) (common.Hash, error) {
	var result common.Hash
	err := ec.Client().CallContext(ctx, &result, "eth_sendBundle", args)
	if nil != err {
		return common.Hash{}, err
	}

	return result, nil
}
