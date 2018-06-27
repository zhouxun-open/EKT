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
            "address": "968b10ebc111ea3434de7333d82e54890c4a2d8c34577e0e54f3464eb88e3b2f",
            "amount": 50000000000000000
        },{
            "address": "ae0ec97c589ff55b856cbad8ba54586453ce2cd17cc202ee7fec30524f33d407",
            "amount": 30000000000000000
        },{
            "address": "04707c449c822a1172003f262ab77607d3f91cafa17e1abbab7a880807bdac0c",
            "amount": 20000000000000000
        }
    ]
}
