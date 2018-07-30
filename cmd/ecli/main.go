package main

import (
	"github.com/EducationEKT/EKT/cmd/ecli/cmd"
	"github.com/EducationEKT/EKT/cmd/ecli/param"
	"github.com/spf13/cobra"
)

var (
	cmds = []*cobra.Command{}
)

func init() {
	cmds = append(cmds, cmd.TransactionCmd, cmd.AccountCmd, cmd.NodeCommand)
}

func main() {
	var RootCmd = &cobra.Command{
		Use: "ecli",
	}
	RootCmd.AddCommand(cmds...)
	RootCmd.PersistentFlags().BoolVar(&param.Localnet, "localnet", false, "localnet peers")
	RootCmd.PersistentFlags().BoolVar(&param.Testnet, "testnet", false, "testnet peers")
	RootCmd.PersistentFlags().BoolVar(&param.Mainnet, "mainnet", false, "mainnet peers")
	RootCmd.Execute()
}
