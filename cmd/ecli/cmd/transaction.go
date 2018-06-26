package cmd

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/EducationEKT/EKT/cmd/ecli/param"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/EKT/io/ekt8/util"
	"github.com/spf13/cobra"
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
	}...)
}

func SendTransaction(cmd *cobra.Command, args []string) {
	pub, private := crypto.GenerateKeyPair()
	fmt.Println(hex.EncodeToString(private))
	address := common.FromPubKeyToAddress(pub)
	fmt.Println(hex.EncodeToString(address))
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
	from := common.FromPubKeyToAddress(pubKey)
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
	fmt.Print("Input the address who recieve this token:")
	input.Scan()
	to := input.Text()
	tx := common.Transaction{
		From:         hex.EncodeToString(from),
		To:           to,
		TimeStamp:    time.Now().UnixNano() / 1e6,
		Amount:       int64(amount),
		Nonce:        1,
		Data:         "",
		TokenAddress: tokenAddress,
	}
	sign, err := crypto.Crypto(crypto.Sha3_256([]byte(tx.String())), privKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	tx.Sign = hex.EncodeToString(sign)
	sendTransaction(tx)
}

func sendTransaction(tx common.Transaction) {
	for _, node := range param.GetPeers() {
		url := fmt.Sprintf(`http://%s:%d/transaction/api/newTransaction`, node.Address, node.Port)
		resp, err := util.HttpPost(url, tx.Bytes())
		fmt.Println(string(resp), err)
		if err == nil {
			break
		}
	}
}
