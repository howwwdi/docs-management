package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

type ContractData struct {
	ABI      interface{} `json:"abi"`
	Bytecode string      `json:"bytecode"`
}

func main() {
	godotenv.Load()

	rpcUrl := os.Getenv("RPC_URL")
	hexPrivateKey := os.Getenv("PRIVATE_KEY")

	if _, err := os.Stat(`artifacts\contracts\DocumentManagement.sol\DocumentManagement.json`); os.IsNotExist(err) {
		log.Fatalf("File not found")
	}

	fileContent, err := os.ReadFile(`artifacts\contracts\DocumentManagement.sol\DocumentManagement.json`)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var contractData ContractData
	err = json.Unmarshal(fileContent, &contractData)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	defer client.Close()

	privateKey, err := crypto.HexToECDSA(hexPrivateKey[2:])
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey))
	if err != nil {
		log.Fatalf("Failed to retrieve nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to retrieve gas price: %v", err)
	}

	chainID := big.NewInt(421614)

	data, err := hex.DecodeString(contractData.Bytecode[2:])
	if err != nil {
		log.Fatalf("Failed to decode contract bytecode: %v", err)
	}

	tx := types.NewContractCreation(nonce, big.NewInt(0), uint64(3000000), gasPrice, data)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	log.Printf("Contract deployment transaction sent. Tx hash: %s\n", signedTx.Hash().Hex())
}
