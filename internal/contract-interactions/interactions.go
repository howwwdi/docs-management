package contract_interactions

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

func LoadABI(filePath string) (abi.ABI, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to open ABI file: %v", err)
	}
	defer file.Close()

	var parsedABI abi.ABI
	err = json.NewDecoder(file).Decode(&parsedABI)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to parse ABI: %v", err)
	}

	return parsedABI, nil
}

func InvokeContractMethod(client *ethclient.Client, contractAddress common.Address, method string, args []interface{}, privateKey *ecdsa.PrivateKey, abi abi.ABI) (common.Hash, error) {
	data, err := abi.Pack(method, args...)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to pack method call: %v", err)
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get nonce: %v", err)
	}

	gasLimit := uint64(300000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get gas price: %v", err)
	}

	tx := types.NewTransaction(
		nonce,
		contractAddress,
		big.NewInt(0),
		gasLimit,
		gasPrice,
		data,
	)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get chain ID: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to sign transaction: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to send transaction: %v", err)
	}

	return signedTx.Hash(), nil
}

