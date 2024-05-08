package transfer

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func SimpleTransfer(ec *ethclient.Client, privateKey *ecdsa.PrivateKey, fromNonce *uint64, to common.Address, amount *big.Int) (*types.Transaction, error) {

	var (
		nonce uint64 = 0
		err   error
	)

	// 获取nonce
	if nil == fromNonce {
		// 获取nonce
		nonce, err = ec.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey))
		if nil != err {
			return nil, err
		}
	} else {
		nonce = *fromNonce
	}

	// 获取gasPrice
	gasPrice, err := ec.SuggestGasPrice(context.Background())
	if nil != err {
		return nil, err
	}

	// 获取chainID
	chainID, err := ec.NetworkID(context.Background())
	if nil != err {
		return nil, err
	}

	// 创建交易
	tx := types.NewTransaction(nonce, to, amount, 21000, gasPrice, nil)

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if nil != err {
		return nil, err
	}

	return signedTx, nil

}
