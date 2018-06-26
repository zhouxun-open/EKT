package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/spf13/cobra"
)

var AccountCmd *cobra.Command

func init() {
	AccountCmd = &cobra.Command{
		Use:   "account",
		Short: "Create account",
	}
	AccountCmd.AddCommand([]*cobra.Command{
		&cobra.Command{
			Use:   "new",
			Short: "Generate public key and private key.",
			Run:   NewAccount,
		},
	}...)
}

func NewAccount(cmd *cobra.Command, args []string) {
	pubKey, privateKey := crypto.GenerateKeyPair()
	fmt.Println("Please save your Private Key: ", hex.EncodeToString(privateKey))
	fmt.Println("Your address is: ", hex.EncodeToString(common.FromPubKeyToAddress(pubKey)))
}
