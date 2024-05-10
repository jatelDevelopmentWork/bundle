package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/jatelDevelopmentWork/bundle/authenticate"
	"github.com/jatelDevelopmentWork/bundle/transfer"
	"github.com/jatelDevelopmentWork/bundle/types"
)

var (
	aWorkPrivateKey *ecdsa.PrivateKey
	aWorkAddress    common.Address // 0x4307ffd08477668dC6d9f49f90b084B1f1CCC82b

	bWorkPrivateKey *ecdsa.PrivateKey
	bWorkAddress    common.Address // 0x9FD5bD701Fc8105E46399104AC4B8c1B391df760

	// token
	username = "jatel"
	password = "123456"

	// test net
	bscTestNet = "https://data-seed-prebsc-1-s1.binance.org:8545/"
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
	client, err := ethclient.Dial(bscTestNet)
	if err != nil {
		fmt.Println("Dial err", err)
		return
	}
	defer client.Close()

	// 设置认证方式
	authenticate.SetBasicAuth(client, username, password)

	// 获取区块高度
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		fmt.Println("HeaderByNumber err", err)
		return
	}

	fmt.Println("block number:", header.Number.String())

	// 获取nonce
	nonce, err := client.PendingNonceAt(context.Background(), aWorkAddress)
	if nil != err {
		fmt.Println("PendingNonceAt err", err)
		return
	}

	// build bundle
	signedTx, err := transfer.SimpleTransfer(client, aWorkPrivateKey, &nonce, bWorkAddress, big.NewInt(1.345e9))
	if nil != err {
		fmt.Println("SimpleTransfer err", err)
		return
	}

	fmt.Println("signedTx hash:", signedTx.Hash().Hex())

	// serialize transaction
	rawTxBytes, err := rlp.EncodeToBytes(signedTx)
	if nil != err {
		fmt.Println("EncodeToBytes err", err)
		return
	}
	txsList := []hexutil.Bytes{rawTxBytes}

	// next nonce
	nonce = nonce + 1
	signedTx, err = transfer.SimpleTransfer(client, aWorkPrivateKey, &nonce, bWorkAddress, big.NewInt(1.345e9))
	if nil != err {
		fmt.Println("SimpleTransfer err", err)
		return
	}

	fmt.Println("signedTx hash:", signedTx.Hash().Hex())

	// serialize transaction
	rawTxBytes, err = rlp.EncodeToBytes(signedTx)
	if nil != err {
		fmt.Println("EncodeToBytes err", err)
		return
	}
	txsList = append(txsList, rawTxBytes)

	// next nonce
	nonce = nonce + 1
	signedTx, err = transfer.SimpleTransfer(client, aWorkPrivateKey, &nonce, bWorkAddress, big.NewInt(1.345e9))
	if nil != err {
		fmt.Println("SimpleTransfer err", err)
		return
	}

	fmt.Println("signedTx hash:", signedTx.Hash().Hex())

	// serialize transaction
	rawTxBytes, err = rlp.EncodeToBytes(signedTx)
	if nil != err {
		fmt.Println("EncodeToBytes err", err)
		return
	}
	txsList = append(txsList, rawTxBytes)

	// send bundle
	args := types.SendBundleArgs{
		Txs:            txsList,
		MaxBlockNumber: header.Number.Uint64() + 50,
	}

	client.Close()
	client, err = ethclient.Dial("http://144.76.100.145:8575")
	if err != nil {
		fmt.Println("Dial err", err)
		return
	}

	// 设置认证方式
	authenticate.SetBasicAuth(client, username, password)
	err = types.SendBundle(client, context.Background(), &args)
	if nil != err {
		fmt.Println("SendBundle err", err)
		return
	}

	// build transaction
	nonce = nonce + 1
	client.Close()
	client, err = ethclient.Dial(bscTestNet)
	if err != nil {
		fmt.Println("Dial err", err)
		return
	}
	signedTx, err = transfer.SimpleTransfer(client, aWorkPrivateKey, &nonce, bWorkAddress, big.NewInt(1.345e9))
	if nil != err {
		fmt.Println("SimpleTransfer err", err)
		return
	}

	fmt.Println("signedTx hash:", signedTx.Hash().Hex())

	// send transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if nil != err {
		fmt.Println("SendTransaction err", err)
		return
	}

	// // wait
	// _, err = bind.WaitMined(context.Background(), client, signedTx)
	// if err != nil {
	// 	fmt.Println("wait transaction error!!!")
	// 	return
	// }

}
