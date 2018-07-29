package dispatcher

import (
	"encoding/hex"
	"errors"

	"github.com/EducationEKT/EKT/blockchain_manager"
	"github.com/EducationEKT/EKT/core/types"
	"github.com/EducationEKT/EKT/core/userevent"
)

func NewTransaction(transaction *userevent.Transaction) error {
	// 主币的tokenAddress为空
	if transaction.TokenAddress != "" {
		tokenAddress, err := hex.DecodeString(transaction.TokenAddress)
		if err != nil {
			return err
		}
		currentBlock := blockchain_manager.GetMainChain().GetLastBlock()
		var token types.Token
		err = currentBlock.TokenTree.GetInterfaceValue(tokenAddress, &token)
		if err != nil || token.Name == "" || token.Decimals <= 0 || token.Total <= 0 {
			return err
		}
	}
	if !userevent.Validate(transaction) {
		return errors.New("error signature")
	}
	blockchain_manager.GetMainChain().NewTransaction(transaction)
	return nil
}
