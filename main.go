package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sol-swap-parse/parse"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

/*
Example Transactions:
- Orca: 2kAW5GAhPZjM3NoSrhJVHdEpwjmq9neWtckWnjopCfsmCGB27e3v2ZyMM79FdsL4VWGEtYSFi1sF1Zhs7bqdoaVT
- Pumpfun: 4Cod1cNGv6RboJ7rSB79yeVCR4Lfd25rFgLY3eiPJfTJjTGyYP1r2i1upAYZHQsWDqUbGd1bhTRm1bpSQcpWMnEz
- Jupiter: DBctXdTTtvn7Rr4ikeJFCBz4AtHmJRyjHGQFpE59LuY3Shb7UcRJThAXC7TGRXXskXuu9LEm9RqtU6mWxe5cjPF
- Meteora DLMM: 125MRda3h1pwGZpPRwSRdesTPiETaKvy4gdiizyc3SWAik4cECqKGw2gggwyA1sb2uekQVkupA2X9S4vKjbstxx3
- Rayd V4: 5kaAWK5X9DdMmsWm6skaUXLd6prFisuYJavd9B62A941nRGcrmwvncg3tRtUfn7TcMLsrrmjCChdEjK3sjxS6YG9
- Meteora Pools Program: 4uuw76SPksFw6PvxLFkG9jRyReV1F4EyPYNc3DdSECip8tM22ewqGWJUaRZ1SJEZpuLJz1qPTEPb2es8Zuegng9Z //! not supported yet
- Multiple AMMs: 46Jp5EEUrmdCVcE3jeewqUmsMHhqiWWtj243UZNDFZ3mmma6h2DF4AkgPE9ToRYVLVrfKQCJphrvxbNk68Lub9vw //! not supported yet
*/

func main() {
	rpcClient := rpc.New(rpc.MainNetBeta.RPC)
	txSig := solana.MustSignatureFromBase58("5kaAWK5X9DdMmsWm6skaUXLd6prFisuYJavd9B62A941nRGcrmwvncg3tRtUfn7TcMLsrrmjCChdEjK3sjxS6YG9")

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
		log.Fatalf("error getting tx: %s", err)
	}

	parser, err := parse.NewTransactionParser(tx)
	if err != nil {
		log.Fatalf("error creating orca parser: %s", err)
	}

	transactionData, err := parser.ParseTransaction()
	if err != nil {
		log.Fatalf("error parsing transaction: %s", err)
	}

	// marshalledData, _ := json.MarshalIndent(transactionData, "", "  ")
	// fmt.Println(string(marshalledData))

	swapData, err := parser.ProcessSwapData(transactionData)
	if err != nil {
		log.Fatalf("error processing swap data: %s", err)
	}

	marshalledSwapData, _ := json.MarshalIndent(swapData, "", "  ")
	fmt.Println(string(marshalledSwapData))

}