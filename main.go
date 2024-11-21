package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	solanaswapgo "github.com/franco-bianco/solanaswap-go/solanaswap-go"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

var rpcClient = rpc.New("https://mainnet.helius-rpc.com/?api-key=970d93f2-5af8-4d62-901f-8d5280ef7c86")

func parseTransactionHandler(c *gin.Context) {
	// Get the txHash from the query parameters
	txHash := c.Query("txHash")
	if txHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "txHash parameter is required"})
		return
	}

	// Parse the transaction signature
	txSig, err := solana.SignatureFromBase58(txHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid txHash format"})
		return
	}

	// Get transaction details
	var maxTxVersion uint64 = 0
	tx, err := rpcClient.GetTransaction(
		context.TODO(),
		txSig,
		&rpc.GetTransactionOpts{
			Commitment:                     rpc.CommitmentConfirmed,
			MaxSupportedTransactionVersion: &maxTxVersion,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching transaction", "details": err.Error()})
		return
	}

	// Parse the transaction
	parser, err := solanaswapgo.NewTransactionParser(tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating transaction parser", "details": err.Error()})
		return
	}

	transactionData, err := parser.ParseTransaction()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error parsing transaction", "details": err.Error()})
		return
	}

	// Process swap data
	swapData, err := parser.ProcessSwapData(transactionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error processing swap data", "details": err.Error()})
		return
	}

	// Create the final response
	response := gin.H{
		"transactionData": transactionData,
		"swapData":        swapData,
		"slot":            tx.Slot,
		"blockTime":       tx.BlockTime,
		"fee":             tx.Meta.Fee,
	}

	// Respond with JSON
	c.JSON(http.StatusOK, response)
}

func main() {
	// Create a Gin router
	router := gin.Default()

	// Define the API route
	router.GET("/parseTransaction", parseTransactionHandler)

	// Start the server
	router.Run(":8080") // Listen on port 8080
}
