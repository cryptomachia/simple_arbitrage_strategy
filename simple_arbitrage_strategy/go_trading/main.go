package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
   	"github.com/cryptomachia/simple_arbitrage_strategy/contract"
    	"github.com/cryptomachia/simple_arbitrage_strategy/dex"
)

const (
	profitThreshold = 0.01 // 1% profit threshold to trigger arbitrage
)

func main() {
	client, err := ethclient.Dial("https://arbitrum-mainnet.infura.io/v3/100b6f70e4f44f7287031457bdd26b0f") // Replace with the appropriate Ethereum WebSocket URL
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	// Fetch prices from the DEXes
	traderJoePrice, err := dex.TraderJoePrice(client)
	if err != nil {
		log.Fatalf("Failed to fetch Trader Joe price: %v", err)
	}

	curvePrice, err := dex.CurvePrice(client)
	if err != nil {
		log.Fatalf("Failed to fetch Curve price: %v", err)
	}

	zyberPrice, err := dex.ZyberPrice(client)
	if err != nil {
		log.Fatalf("Failed to fetch Zyber price: %v", err)
	}

	fmt.Printf("Trader Joe price: %s\n", traderJoePrice.Text('f', 18))
	fmt.Printf("Curve price: %s\n", curvePrice.Text('f', 18))
	fmt.Printf("Zyber price: %s\n", zyberPrice.Text('f', 18))

	// Implement the simple arbitrage strategy
	swapperAddress := common.HexToAddress("DEPLOYED_SWAPPER_CONTRACT_ADDRESS") // Replace with the deployed Swapper contract address
	swapper, err := contract.NewSwapper(swapperAddress, client)
	if err != nil {
		log.Fatalf("Failed to instantiate Swapper contract: %v", err)
	}

	privateKey, err := crypto.HexToECDSA("PRIVATE_KEY") // Replace with your private key
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(42161))
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}

	var (
		source      uint8
		destination uint8
	)

	// Determine the DEX with the highest price
	maxPrice := big.NewFloat(0)
	if traderJoePrice.Cmp(curvePrice) > 0 && traderJoePrice.Cmp(zyberPrice) > 0 {
		maxPrice = traderJoePrice
		source = 0
	} else if curvePrice.Cmp(traderJoePrice) > 0 && curvePrice.Cmp(zyberPrice) > 0 {
		maxPrice = curvePrice
		source = 1
	} else {
		maxPrice = zyberPrice
		source = 2
	}

	// Determine the DEX with the lowest price
	minPrice := big.NewFloat(0)
	if traderJoePrice.Cmp(curvePrice) < 0 && traderJoePrice.Cmp(zyberPrice) < 0 {
		minPrice = traderJoePrice
		destination = 0
	} else if curvePrice.Cmp(traderJoePrice) < 0 && curvePrice.Cmp(zyberPrice) < 0 {
		minPrice = curvePrice
		destination = 1
	} else {
		minPrice = zyberPrice
		destination = 2
	}

	// Calculate the profit ratio and percentage
	profitRatio := new(big.Float).Quo(maxPrice, minPrice)
	profitDifference := new(big.Float).Sub(profitRatio, big.NewFloat(1))
	profitPercentage := new(big.Float).Mul(profitDifference, big.NewFloat(100))
	fmt.Printf("Profit percentage: %s%%\n", profitPercentage.Text('f', 2))

	// Check if the profit threshold is met and execute the arbitrage strategy
	if profitPercentage.Cmp(big.NewFloat(profitThreshold)) > 0 {
		fmt.Println("Arbitrage opportunity detected!")

		amount := big.NewInt(1000000000000000000) // 1 WETH or 1e18, adjust the amount as needed

		tx, err := swapper.Swap(auth, source, destination, amount, big.NewInt(0))
		if err != nil {
			log.Fatalf("Failed to execute swap transaction: %v", err)
		}

		fmt.Printf("Swap transaction sent: %s\n", tx.Hash().Hex())

		// Wait for the transaction to be mined
		ctx := context.Background()
		receipt, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			log.Fatalf("Failed to wait for transaction receipt: %v", err)
		}

		if receipt.Status == 1 {
			fmt.Println("Arbitrage successful!")
		} else {
			fmt.Println("Arbitrage failed.")
		}
	} else {
		fmt.Println("No arbitrage opportunity.")
	}
}
