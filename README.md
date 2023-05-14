# Eddie Huang Arbitrum Simple Arbitrage

## A. Swapper Smart Contract

A smart contract named Swapper was created using the Solidity programming language. The purpose of this contract is to facilitate buying and selling tokens on the WETH/USDT pool on Arbitrum using the following DEXes:

- TraderJoe V2.1: 0xd387c40a72703b38a5181573724bcaf2ce6038a5
- Curve V2: 0x960ea3e3C7FB317332d990873d354E18d7645590
- Zyber V3: 0xc4e9fce31518d4233c224772dc5532d49c5354c0

The Swapper smart contract has the following key features:

1. Functions for interacting with each DEX: Each DEX has its specific protocol for executing swaps. The smart contract contains functions to handle these protocols and perform the required operations on each DEX.
2. Token transfer handling: The contract includes functions to approve, transfer, and retrieve the balance of tokens involved in the swap operations.
3. Swap execution: The contract contains a performSwap function, which receives the source and destination DEXes, the amount of tokens to swap, and the minimum output amount. This function orchestrates the token transfers and interactions with the DEXes to execute the arbitrage trade.

## B. Go Framework

A Go framework was developed to interact with the Swapper smart contract and perform the arbitrage strategy. The framework is organized into two main parts:

1. DEX-specific Go files: Three separate Go files (traderjoe.go, curve.go, zyber.go) are included within the dex package. These files are responsible for fetching real-time prices from the respective DEXes using websockets.
2. Main Go file: The main.go file, located in the go_trading package, subscribes to the websockets provided by the DEX-specific Go files, calculates the arbitrage opportunities, and performs swaps by interacting with the Swapper smart contract.

The Go framework performs the following steps:

a) Fetch real-time prices: The framework subscribes to the websockets provided by the DEX-specific Go files to fetch real-time prices from each DEX.
b) Calculate profit percentage: It then calculates the profit percentage for each possible arbitrage opportunity between the DEXes.
c) Identify source and destination DEXes: The framework determines the DEX with the highest price (source) and the DEX with the lowest price (destination) for the arbitrage trade.
d) Execute the arbitrage strategy: If the profit percentage is above a specified threshold, the framework proceeds to execute the arbitrage strategy. It does this by submitting a swap transaction using the Swapper smart contract, providing the source and destination DEXes, the amount of tokens to swap, and the minimum output amount.
e) Monitor the transaction: After submitting the swap transaction, the framework monitors the transaction for completion using websockets to track on-chain transaction logs.
f) Report the result: Once the transaction is completed, the framework reports the result as either successful or failed.

By following this logic, the arbitrage strategy aims to capitalize on price discrepancies between different DEXes in the WETH/USDT pool on the Arbitrum network, generating profits from these opportunities.

## Setup and Execution

Before you begin, ensure you have the following software installed on your machine:

- Node.js
- Go
- Truffle
- Ganache (Optional, for local testing)

Additionally, you will need access to the following:

- A wallet with some ETH for gas fees
- An arbitrum node URL, such as the one provided by Infura

## Limitations
Gas fees and transaction latency: Executing trades on the Arbitrum network involves paying gas fees for each transaction. High gas fees can reduce the profit potential of arbitrage opportunities. Moreover, network congestion can lead to increased transaction latency, causing opportunities to disappear before the transaction is confirmed.

Slippage handling: My current implementation does not provide a mechanism to account for slippage when executing arbitrage trades. In real-world scenarios, slippage can significantly affect the profitability of arbitrage opportunities, particularly in illiquid markets or when executing large orders.

Price feed latency: My Go framework fetches real-time prices from the DEXes using the Arbitrum node URL provided by Infura. Network latency and update frequency can introduce a delay in receiving the most up-to-date prices. This delay may result in missed arbitrage opportunities or executing trades based on outdated price information.

Hardcoded DEX versions: My Swapper smart contract and Go framework are built to interact with specific versions of the DEXes (TraderJoe V2.1, Curve V2, and Zyber V3). If the DEXes release new versions or undergo protocol changes, the existing implementation may become incompatible and require updates to continue functioning properly.

Error handling: The current Go framework implementation does not have comprehensive error handling, making it difficult to diagnose and resolve issues that may arise during the execution of the arbitrage strategy. Improved error handling and logging can help identify potential problems and enhance the robustness of the strategy.

Concurrency: My current Go framework does not take advantage of concurrency features to fetch prices from DEXes in parallel. Implementing concurrent price fetching can improve the speed and efficiency of the arbitrage strategy by reducing the time spent waiting for prices from each DEX.

Liquidity and trade size: My Go framework and Swapper smart contract do not account for liquidity constraints or trade size optimization when executing the arbitrage strategy. In practice, liquidity and trade size can significantly impact the profitability and feasibility of arbitrage opportunities.

Risk management: The current implementation of my Go framework does not include risk management features, such as setting maximum trade sizes or maximum allowable losses. Implementing risk management features can help protect the capital involved in executing the arbitrage strategy.

### Steps to execute the arbitrage strategy:

1. Clone the repository containing the arbitrage strategy:

   `git clone https://github.com/yourusername/yourprojectname.git`
   `cd yourprojectname`

2. Deploy the Swapper smart contract to the Arbitrum network. You can use Truffle, Remix, or any other preferred method for contract deployment. After deploying the contract, note down the Swapper contract address.

3. Open the `main.go` file in the `simple_arbitrage_strategy/go_trading` directory. Update the Swapper contract address with the one you noted down in the previous step:

   `swapperAddress := common.HexToAddress("DEPLOYED_SWAPPER_CONTRACT_ADDRESS") // Replace with the deployed Swapper contract address`

4. Build the Go trading strategy:

   `go build ./simple_arbitrage_strategy/go_trading`

5. Run the built Go trading strategy to start executing the arbitrage strategy:

   `./go_trading`

By following these steps, the Go framework will interact with the Swapper smart contract to execute the arbitrage strategy, taking advantage of price discrepancies between the different DEXes in the WETH/USDT pool on the Arbitrum network.
