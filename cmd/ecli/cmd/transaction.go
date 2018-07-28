package cmd

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/EducationEKT/EKT/cmd/ecli/param"
	"github.com/EducationEKT/EKT/core/types"
	"github.com/EducationEKT/EKT/core/userevent"
	"github.com/EducationEKT/EKT/crypto"
	"github.com/EducationEKT/EKT/util"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

var TransactionCmd *cobra.Command

func init() {
	TransactionCmd = &cobra.Command{
		Use:   "transaction",
		Short: "send transaction or search transaction",
	}
	TransactionCmd.AddCommand([]*cobra.Command{
		&cobra.Command{
			Use:   "send",
			Short: "Send transaction to nodes.",
			Run:   SendTransaction,
		},
		&cobra.Command{
			Use:   "benchtest",
			Short: "Bench test TPS",
			Run:   BenchTest,
		},
	}...)
}

func SendTransaction(cmd *cobra.Command, args []string) {
	fmt.Print("Input your private key: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	privateKey := input.Text()
	privKey, err := hex.DecodeString(privateKey)
	if err != nil {
		fmt.Println("Your private key is not right, exit.")
		os.Exit(-1)
	}
	pubKey, err := crypto.PubKey(privKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	from := types.FromPubKeyToAddress(pubKey)
	fmt.Print("Input token address (Press ENTER if you need send EKT): ")
	input.Scan()
	tokenAddress := input.Text()
	fmt.Print("Amount you want tranfer: ")
	input.Scan()
	amount, err := strconv.Atoi(input.Text())
	if err != nil {
		fmt.Println("You can only input int64 only, exit.")
		os.Exit(-1)
	}
	fmt.Print("Input the address who receive this token: ")
	input.Scan()
	receive := input.Text()
	to, err := hex.DecodeString(receive)
	if err != nil {
		fmt.Println("Error address")
		os.Exit(-1)
	}
	nonce := getAccountNonce(hex.EncodeToString(from))
	tx := userevent.NewTransaction(from, to, time.Now().UnixNano()/1e6, int64(amount), 510000, nonce, "", tokenAddress)
	tx.Signature(privKey)
	sendTransaction(*tx)
}

func BenchTest(cmd *cobra.Command, args []string) {
	fmt.Print("Input your private key: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	privateKey := input.Text()
	privKey, err := hex.DecodeString(privateKey)
	if err != nil {
		fmt.Println("Your private key is not right, exit.")
		os.Exit(-1)
	}
	pubKey, err := crypto.PubKey(privKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	from := types.FromPubKeyToAddress(pubKey)
	fmt.Print("Input the address who receive this token: ")
	input.Scan()
	receive := input.Text()
	to, err := hex.DecodeString(receive)
	if err != nil {
		fmt.Println("Error address")
		os.Exit(-1)
	}
	amount := 1000000
	nonce := getAccountNonce(hex.EncodeToString(from))
	tx := userevent.NewTransaction(from, to, time.Now().UnixNano()/1e6, int64(amount), 510000, nonce, "", "")
	testTPS(tx, privKey)
}

func testTPS(tx *userevent.Transaction, priv []byte) {
	max, min := tx.Nonce+2000, tx.Nonce
	for nonce := max; nonce >= min; nonce-- {
		tx.Nonce = int64(nonce)
		tx.Signature(priv)
		sendTransaction(*tx)
		fmt.Println(tx.String())
	}
	fmt.Println("finish")
}

func sendTransaction(tx userevent.Transaction) {
	for _, node := range param.GetPeers() {
		url := fmt.Sprintf(`http://%s:%d/transaction/api/newTransaction`, node.Address, node.Port)
		resp, err := util.HttpPost(url, tx.Bytes())
		fmt.Println(string(resp), err)
		if err == nil {
			break
		}
	}
}

func getAccountNonce(address string) int64 {
	for _, node := range param.GetPeers() {
		url := fmt.Sprintf(`http://%s:%d/account/api/nonce?address=%s`, node.Address, node.Port, address)
		respBody, err := util.HttpGet(url)
		if err != nil {
			continue
		} else {
			var resp x_resp.XRespBody
			err := json.Unmarshal(respBody, &resp)
			if err != nil {
				continue
			} else {
				nonce, ok := resp.Result.(float64)
				if ok {
					return int64(nonce) + 1
				} else {
					panic("Can not get nonce from remote peer")
				}
			}
		}
	}
	return -1
}
