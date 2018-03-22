# EKT跨公链报文协议

EKT是一个多链多共识的公链，在EKT中的多条链共享同一个用户系统，因此在EKT内部是“天然可以跨链交易”的，即A链的用户可以将自己的A1币发给在A链上没有币但是在B链上有币的用户。也就是说，EKT公链在完成了多链之间互相隔离的同时有保证了链与链之间的关系。同时，EKT还支持跨公链的交易，因为EKT采用DPOS共识，所以第三方可以直接调用主链委托人节点进行查询。具体报文如下。

## EKT跨链握手协议
EKT判断一个token是否存在只需要一个http接口即可，示例报文如下：
```
POST /crosschain/api/handshake HTTP/1.1
Host: 127.0.0.1:19951
User-Agent: curl/7.47.0
Accept: application/json
Content-Type:application/json
Content-Length: 86

{"publicChainId": "000000000FFFFFFFF00000001", "tokenId": "000000000FFFFFFFF0FFF001F"}
```
如果存在则返回true，返回的body示例如下：
 `{"exist": true}`
如果不存在则返回false，返回的body示例如下：
 `{"exist": false}`

## EKT跨链注册协议
对于一个未在EKT中注册的公链如果要在EKT中注册跨链功能，需要双方技术先对对方的公链数据格式进行对接，双方增加自己的operation code之后再调用接口进行注册，注册之后双方即可进行跨链交易。（双方对接完成之后，会对公链的各个数据进行记录，包括新增的跨链操作代码和bootnode等信息，代码新增完成之后需要调用接口激活，当所有DPOS节点都同意之后接口返回成功，否则返回失败）注册接口的示例报文如下：
```
POST /crosschain/api/regist HTTP/1.1
Host: localhost:19951
User-Agent: curl/7.47.0
Accept: application/json
Content-Type:application/json
Content-Length: 114

{"publicChainId": "000000000FFFFFFFF00000001", "tokenId": "000000000FFFFFFFF0FFF001F", "opCodeId": "0000FFFF00FF"}
```
如果注册成功则返回true，返回的body示例如下：
 `{"exist": true}`
如果注册失败则返回false，返回的body示例如下：
 `{"exist": false}`

## EKT跨链操作伪代码

发送跨链交易的伪代码
```
func (sender Sender) SendCrossChain(destChain []byte, to []byte, value int64) error {
    if balances[sender.Address] < value {
        return errors.New("no enough balance")
    }
    crosschain.SendEvent(destChain, event.Event{Name: "CrossChainEvent", To: hex.EncodeToString(to), Value: value})
    balances[sender.Address] -= value
    balances[destChain] += value
}
```

处理跨链交易的伪代码
```
func RecieveCrossChainEvent(sendChain []byte, to []byte, value int64) {
    if event.Name == "CrossChainEvent" {
        if balances[sendChain] < value {
	    return errors.New("no enough balance")
	}
	balances[to] +=value
        balances[sendChain] -= value
    }
}
```

当今的区块链技术百花齐放，各成一派，彼此之间无法通讯，EKT作为一个先进的区块链公链，为大家提供了一个通用的跨链协议，助力让区块链世界更加美好。
