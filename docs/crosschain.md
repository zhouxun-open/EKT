# EKT跨公链报文协议

EKT是一个多链多共识的公链，在EKT中的多条链共享同一个用户系统，因此在EKT内部是“天然可以跨链交易”的，即A链的用户可以将自己的A1币发给在A链上没有币但是在B链上有币的用户。也就是说，EKT公链在完成了多链之间互相隔离的同时有保证了链与链之间的关系。同时，EKT还支持跨公链的交易，因为EKT采用DPOS共识，所以第三方可以直接调用主链委托人节点进行查询。具体报文如下。

## EKT跨链握手协议
EKT判断一个token是否存在只需要一个http接口即可，示例报文如下：
```
POST /crosschain/api/handshake HTTP/1.1
Host: 127.0.0.1:19951
User-Agent: curl/7.47.0
Accept: */*
Content-Type:application/json
Content-Length: 86

{"publicChainId": "000000000FFFFFFFF00000001", "tokenId": "000000000FFFFFFFF0FFF001F"}
```
如果存在则返回true，返回的body示例如下：
`{"exist": true}`
如果不存在则返回false，返回的body示例如下：
`{"exist": false}`
