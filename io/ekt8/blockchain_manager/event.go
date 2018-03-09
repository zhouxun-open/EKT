package blockchain_manager

//
//func NewEvent(block *blockchain.Block, event pkg_event.Event) error {
//	if !event.ValidateEvent() {
//		return errors.New("Invalid event")
//	}
//	if event.EventType == pkg_event.NewAccountEvent {
//		param := event.EventParam.(pkg_event.NewAccountParam)
//		if block.ExistAddress(param.Address) {
//			return errors.New("Address Exist")
//		} else {
//			if err := block.NewAccount(param.Address, param.PubKey); err == nil {
//
//			}
//		}
//	}
//	return nil
//}
