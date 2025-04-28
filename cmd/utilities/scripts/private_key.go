package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func main() {
	// Generate 128 bits entropy
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		log.Fatal(err)
	}

	// Generate mnemonic
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mnemonic:", mnemonic)

	// Generate seed from mnemonic
	seed := bip39.NewSeed(mnemonic, "")

	// Generate master key
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		log.Fatal(err)
	}

	// Derive m/44'/60'/0'/0/0
	purpose, _ := masterKey.NewChildKey(44 + bip32.FirstHardenedChild)
	coinType, _ := purpose.NewChildKey(60 + bip32.FirstHardenedChild)
	account, _ := coinType.NewChildKey(0 + bip32.FirstHardenedChild)
	change, _ := account.NewChildKey(0)
	addressIndex, _ := change.NewChildKey(0)

	privateKey, err := crypto.ToECDSA(addressIndex.Key)
	if err != nil {
		log.Fatal(err)
	}

	// Print Private Key to import into MetaMask
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Printf("Private Key (hex): %x\n", privateKeyBytes)

	// Print Ethereum Address
	publicAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Printf("Public Address: %s\n", publicAddress.Hex())
}
