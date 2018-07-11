package cmd

import (
	"encoding/hex"
	"fmt"
	"os"
	"bufio"
	"time"

	"encoding/json"
	"github.com/EducationEKT/EKT/cmd/ecli/param"
	"github.com/EducationEKT/EKT/core/common"
	"github.com/EducationEKT/EKT/crypto"
	"github.com/EducationEKT/EKT/util"
	"github.com/spf13/cobra"
	"xserver-go/x_http/x_resp"
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
	/*	fmt.Print("Input your private key: ")
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
		to := input.Text()*/
	fmt.Print("Input your private key: ")
        input := bufio.NewScanner(os.Stdin)
        input.Scan()
        privKey := input.Text()
	to := "ae0ec97c589ff55b856cbad8ba54586453ce2cd17cc202ee7fec30524f33d407"
	tokenAddress := ""
	priv, _ := hex.DecodeString(privKey)
	pub, _ := crypto.PubKey(priv)
	from := common.FromPubKeyToAddress(pub)
	amount := 100000000
	nonce := getAccountNonce(hex.EncodeToString(from))
	tx := common.Transaction{
		From:         hex.EncodeToString(from),
		To:           to,
		TimeStamp:    time.Now().UnixNano() / 1e6,
		Amount:       int64(amount),
		Nonce:        nonce,
		Data:         "",
		TokenAddress: tokenAddress,
	}
	sign, err := crypto.Crypto(crypto.Sha3_256([]byte(tx.String())), priv)
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

func getAccountNonce(address string) int64 {
	for _, node := range param.GetPeers() {
		url := fmt.Sprintf(`http://%s:%d/account/api/nonce?address=%s`, node.Address, node.Port, address)
		respBody, err := util.HttpGet(url)
		fmt.Println(string(respBody))
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
