package dispatcher

import (
	"encoding/hex"
	"errors"

	"github.com/EducationEKT/EKT/io/ekt8/blockchain_manager"
	"github.com/EducationEKT/EKT/io/ekt8/core/common"
	"github.com/EducationEKT/EKT/io/ekt8/ctxlog"
)

func NewTransaction(log *ctxlog.ContextLog, transaction *common.Transaction) error {
	// 主币的tokenAddress为空
	if transaction.TokenAddress != "" {
		log.Log("tokenAdddress", transaction.TokenAddress)
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
	log.Log("transfer EKT", true)
	if !transaction.Validate() {
		log.Log("validate", "error signature")
		return errors.New("error signature")
	}
	log.Log("Validate Success", true)
	if !blockchain_manager.GetMainChain().NewTransaction(log, transaction) {
		log.Log("Error", true)
		return errors.New("error transaction")
	}
	log.Log("success", true)
	return nil
}
