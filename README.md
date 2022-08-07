# btc-transaction

Comand line app to create Bitcoin transactions

Create Bitcoin transactions offline

## Create a signed transaction

### Example command

To create an offline transaction from a given json file:

```
go run ./main.go create -i input.json
```

### Example of input file

```
{
    "testnet": true,
    "privateKey": "5HusYj2b2x4nroApgfvaSfKYZhRbKFH41bVyPooymbC6KfgSXdD",
    "destination": "mkHS9ne12qx9pS9VojpwU5xtRd4T7X7ZUt",
    "amount": 10000,
    "fee": 250,
    "vout": 1,
    "lastHash": "5b8295afbdc7c682a34a9af0f397c2d90f6d7a625702d3af7ca5a63536722b37",
    "utxos":
      [
        {
          "txid": "b0c484f4c5d190b4849c2522be9442d975a711109e27b2e79c970b1a4ec662bd",
          "vout": 0,
          "status": {
            "confirmed": true,
            "block_height": 2315276,
            "block_hash": "000000007baa9471e6ca314ad0a64e62cb2e27ffdaae007f2e69a99ebcd9a38b",
            "block_time": 1659848137
          },
          "value": 10000
        },
        {
          "txid": "80d0a05f8b74654ccef60246b0229814e2f2a59b8a856c1f060df50cad3b82ca",
          "vout": 0,
          "status": {
            "confirmed": true,
            "block_height": 2315277,
            "block_hash": "00000000000000051d0e8ca10b96b79d6ed6266072a8df1e9232e88cbd532dab",
            "block_time": 1659848618
          },
          "value": 50000
        }
      ]
}
```

### UTXOs

To retrieve an array of currently available unspent transactions, use mempool:

```
curl -X GET https://mempool.space/testnet/api/address/mxEEptDqy8pHvcxr98v6HCkdTUb1xg5tYf/utxo | jq
```

## Fetch fees

```
curl -X GET https://api.blockcypher.com/v1/btc/test3 | jq
```
