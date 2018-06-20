package dispatcher

import (
	"encoding/hex"
	"errors"

	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
)

func NewTransaction(transaction *common.Transaction) error {
	// 主币的tokenAddress为空
	if transaction.TokenAddress != "" {
		tokenAddress, err := hex.DecodeString(transaction.TokenAddress)
		if err != nil {
			return err
		}
		currentBlock := blockchain_manager.GetMainChain().GetLastBlock()
		var token common.Token
		err = currentBlock.TokenTree.GetInterfaceValue(tokenAddress, &token)
		if err != nil || token.Name == "" || token.Decimals <= 0 || token.Total <= 0 {
			return err
		}
	}
	if !transaction.Validate() {
		return errors.New("error signature")
	}
	if !blockchain_manager.GetMainChain().NewTransaction(transaction) {
		return errors.New("error transaction")
	}
	return nil
}
