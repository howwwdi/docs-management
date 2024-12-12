package corefuncs

import (
	"context"
	contract_interactions "docs-managment/internal/contract-interactions"
	pinata_api "docs-managment/internal/pinata-api"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func UploadFile(client *ethclient.Client, fileName string, contractAddress common.Address, hexPrivateKey string, abi abi.ABI, data []byte) (*big.Int, string, string) {
	godotenv.Load()

	privateKey, err := crypto.HexToECDSA(hexPrivateKey[2:])
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	hash, err := pinata_api.UploadFile(fileName, data)
	if err != nil {
		log.Fatal(err)
	}

	methodName := "addDocument"
	methodArgs := []interface{}{
		hash,
		fileName,
	}

	txHash, err := contract_interactions.InvokeContractMethod(client, contractAddress, methodName, methodArgs, privateKey, abi)
	if err != nil {
		log.Fatalf("Failed to invoke contract method: %v", err)
	}

	log.Printf("Transaction sent: %s\n", txHash.Hex())
	var receipt *types.Receipt
	for {
		receipt, err = client.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			log.Println("Waiting for transaction to be confirmed...")
		} else {
			break
		}
	}

	for _, vLog := range receipt.Logs {
		uintValue := new(big.Int)
		uintValue, success := uintValue.SetString(vLog.Topics[1].String()[2:], 16)
		if !success {
			log.Println("Failed to convert hex to uint256")
			return big.NewInt(0), "", ""
		} else {
			log.Printf("Document id: %s\n", uintValue)
			return uintValue, txHash.Hex(), hash
		}
	}

	return big.NewInt(0), "", ""
}

func GetDocument(client *ethclient.Client, contractAddress common.Address, abi abi.ABI, docId uint64) (string, string, error) {
	data, err := abi.Pack("getDocument", big.NewInt(int64(docId)))
	if err != nil {
		return "", "", fmt.Errorf("failed to pack method call: %v", err)
	}

	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to call contract: %v", err)
	}

	var (
		ipfsHash  string
		fileName  string
		timestamp *big.Int
		owner     common.Address
	)

	err = abi.UnpackIntoInterface(&[]interface{}{&ipfsHash, &fileName, &timestamp, &owner}, "getDocument", result)
	if err != nil {
		return "", "", fmt.Errorf("failed to unpack response: %v", err)
	}

	fmt.Printf("Document ID: %d\n", docId)
	fmt.Printf("IPFS Hash: %s\n", ipfsHash)
	fmt.Printf("File Name: %s\n", fileName)
	fmt.Printf("Timestamp: %d\n", timestamp)
	fmt.Printf("Owner: %s\n", owner.Hex())

	return ipfsHash, fileName, nil
}

func DeleteFile(client *ethclient.Client, contractAddress common.Address, hexPrivateKey string, abi abi.ABI, docId uint64, ipfshash string) {
	privateKey, err := crypto.HexToECDSA(hexPrivateKey[2:])
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	methodName := "deleteDocument"
	methodArgs := []interface{}{
		new(big.Int).SetUint64(docId),
	}

	txHash, err := contract_interactions.InvokeContractMethod(client, contractAddress, methodName, methodArgs, privateKey, abi)
	if err != nil {
		log.Fatalf("Failed to invoke contract method: %v", err)
	}

	log.Printf("Transaction sent: %s\n", txHash.Hex())
	for {
		_, err := client.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			log.Println("Waiting for transaction to be confirmed...")
		} else {
			break
		}
	}
	log.Println("Successfully deleted document from contract")

	if err := pinata_api.DeleteFromPinata(ipfshash); err != nil {
		log.Fatalf("Failed to delete file from Pinata: %v", err)
	} else {
		log.Println("Successfully deleted file from Pinata")
	}
}
