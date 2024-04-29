package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	// ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jatelDevelopmentWork/bundle/transfer"
	"github.com/jatelDevelopmentWork/bundle/types"
)

var (
	aWorkPrivateKey *ecdsa.PrivateKey
	aWorkAddress    common.Address // 0x4307ffd08477668dC6d9f49f90b084B1f1CCC82b

	bWorkPrivateKey *ecdsa.PrivateKey
	bWorkAddress    common.Address // 0x9FD5bD701Fc8105E46399104AC4B8c1B391df760
)

func init() {
	var err error

	// a work address 0x4307ffd08477668dC6d9f49f90b084B1f1CCC82b
	// a work address private key 47d790a96ca73b23fbb65a6b911b8b57a1d915d364f12e2bc7fae83c196c9c97
	aWorkPrivateKey, err = crypto.HexToECDSA("47d790a96ca73b23fbb65a6b911b8b57a1d915d364f12e2bc7fae83c196c9c97")
	if nil != err {
		panic(err)
	}
	aWorkPublicKey := aWorkPrivateKey.Public()
	aWorkPublicKeyECDSA, ok := aWorkPublicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	aWorkAddress = crypto.PubkeyToAddress(*aWorkPublicKeyECDSA)
	fmt.Println("aWorkAddress: ", aWorkAddress.Hex())

	// b work address 0x9FD5bD701Fc8105E46399104AC4B8c1B391df760
	// b work private key e3166f9f62f109d19fb1b73f8d9c647530153cd03822d3091951081bac7f7c5e
	bWorkPrivateKey, err = crypto.HexToECDSA("e3166f9f62f109d19fb1b73f8d9c647530153cd03822d3091951081bac7f7c5e")
	if nil != err {
		panic(err)
	}
	bWorkPublicKey := bWorkPrivateKey.Public()
	bWorkPublicKeyECDSA, ok := bWorkPublicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	bWorkAddress = crypto.PubkeyToAddress(*bWorkPublicKeyECDSA)
	fmt.Println("bWorkAddress: ", bWorkAddress.Hex())

}

func main() {
	// 连接服务器
	client, err := ethclient.Dial("http://localhost:8575")
	if err != nil {
		fmt.Println("Dial err", err)
		return
	}
	defer client.Close()

	// 获取区块高度
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		fmt.Println("HeaderByNumber err", err)
		return
	}

	fmt.Println("block number:", header.Number.String())

	// build bundle
	oneTx, err := transfer.SimpleTransfer(client, aWorkPrivateKey, nil, bWorkAddress, big.NewInt(1.345e9))
	if nil != err {
		fmt.Println("SimpleTransfer err", err)
		return
	}
	txsList := []hexutil.Bytes{oneTx}

	// send bundle
	args := types.SendBundleArgs{
		Txs:            txsList,
		MaxBlockNumber: header.Number.Uint64() + 5,
	}
	err = types.SendBundle(client, context.Background(), &args)
	if nil != err {
		fmt.Println("SendBundle err", err)
		return
	}

	// wait
	// var signedTx ethTypes.Transaction
	// err = rlp.DecodeBytes(oneTx, &signedTx)
	// if nil != err {
	// 	fmt.Println("DecodeBytes err", err)
	// 	return
	// }
	// _, err = bind.WaitMined(context.Background(), client, &signedTx)
	// if err != nil {
	// 	fmt.Println("wait transaction error!!!")
	// 	return
	// }
}
