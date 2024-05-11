package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/chad-chain/chadChain/core/crypto"
	"github.com/chad-chain/chadChain/core/initialize"
	"github.com/chad-chain/chadChain/core/types"
	"github.com/chad-chain/chadChain/core/utils"
)

// POST
// /sendTx?signed={signed}: Used to send a signed transaction where the signed message is sent via query params. The signed value here refers to the rlp encoding of the transaction struct mentioned above including the signature

// GET
// /blockNumber: Returns the recent most block number
// /block?number={number}: Given the block number, returns the contents of a block else nil (Response type: json struct of block type described above)
// /block?hash={hash}: Given the block hash, returns the contents of a block else nil (Response type: json struct of block type described above)
// /tx?hash={hash}: Given the transaction hash, returns the contents of a transaction else nil (Response type: json struct of block type described above). Note that the tx hash is hash of all contents of transaction struct (including v, r, and s).
// /getNonce?address={address}: Given the address, returns the current nonce of that account
// /getBalance?address={address}: Given the address, returns the current balance/amount of that account

var (
	HostAddr string // Host address of the node
)

func main() {
	// Get String from cmd line
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <IP_ADDR>")
		return
	}
	HostAddr = os.Args[1]
	fmt.Println("HostAddr:", HostAddr)
	initialize.Keys()
	SendTx("E97155581dd619246baeC16E832F07B6d9D68773", 100, 0)
}

func SendTx(To string, Value uint64, Nonce uint64) {
	// Send a signed transaction
	tnx := types.UnSignedTx{
		To:    [20]byte(crypto.HexStringToBytes(To)),
		Value: Value,
		Nonce: Nonce,
	}

	// Sign the transaction
	signedTx, err := crypto.SignTransaction(&tnx)

	if err != nil {
		fmt.Println(err)
		return
	}

	rlpEncodedTx, err := utils.EncodeData(signedTx, false)
	if err != nil {
		fmt.Println("Error encoding the signed transaction:", err)
		return
	}

	requestBody, err := json.Marshal(rlpEncodedTx)
	if err != nil {
		fmt.Println("Error marshalling the signed transaction:", err)
		return
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/sendTx", HostAddr), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error sending the signed transaction:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code:", resp.StatusCode)
		return
	}
}

func GetBlockNumber() {
	// Get the most recent block number
	resp, err := http.Get(fmt.Sprintf("http://%s/blockNumber", HostAddr))
	if err != nil {
		fmt.Println("Error getting the block number:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code:", resp.StatusCode)
		return
	}

	// Read the response
	var blockNumber uint64
	err = json.NewDecoder(resp.Body).Decode(&blockNumber)
	if err != nil {
		fmt.Println("Error decoding the block number:", err)
		return
	}

	fmt.Println("Block Number:", blockNumber)
}

func Faucet(address string) {
	// Get some coins from the faucet
	resp, err := http.Get(fmt.Sprintf("http://%s/faucet?address=%s", HostAddr, address))
	if err != nil {
		fmt.Println("Error getting coins from the faucet:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code:", resp.StatusCode)
		return
	}

	fmt.Println("Coins received from the faucet")
}
