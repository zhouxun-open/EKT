{
    "version": "v0.5",
    "dbPath": "/data/EKT/db",
    "logPath": "/data/EKT/log/ekt8.log",
    "debug": true,
    "env": "{{.env}}",
    "node": {
        "peerId": "{{.peerId}}",
        "address": "{{.addr}}",
        "port": {{.port}},
        "addressVersion": {{.addrVer}}
    },
    "privateKey": "{{.privateKey}}",
    "genesisBlock": [
        {
            "address": "{{.genesisAddr0}}",
            "amount": 50000000000000000
        },{
            "address": "{{.genesisAddr1}}",
            "amount": 30000000000000000
        },{
            "address": "{{.genesisAddr2}}",
            "amount": 20000000000000000
        }
    ]
}
