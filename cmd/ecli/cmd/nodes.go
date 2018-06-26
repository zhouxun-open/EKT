package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/EducationEKT/EKT/crypto"
	"github.com/spf13/cobra"
)

var NodeCommand *cobra.Command

func init() {
	NodeCommand = &cobra.Command{
		Use:   "node",
		Short: "Create peerId",
	}
	NodeCommand.AddCommand([]*cobra.Command{
		&cobra.Command{
			Use:   "init",
			Short: "Generate peerId and private key.",
			Run:   GeneratePeerId,
		},
	}...)
}

func GeneratePeerId(cmd *cobra.Command, args []string) {
	pubKey, privateKey := crypto.GenerateKeyPair()
	fmt.Println("Please save your peer private key: ", hex.EncodeToString(privateKey))
	fmt.Println("Your peerId is: ", hex.EncodeToString(crypto.Sha3_256(pubKey)))
}
