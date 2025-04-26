package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Generate a private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	// Convert it to bytes
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Printf("Private Key: %x\n", privateKeyBytes)

	// Generate the public address from private key
	publicAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Printf("Public Address: %s\n", publicAddress.Hex())
}
