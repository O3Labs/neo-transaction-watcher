# neo-transaction-watcher
A headless daemon to connect to the NEO network to receive TX message.

Change the endpoint of the network that you want to connect to and run. You will be able to see transaction coming in in near real-time.

```
go run main.go
```

In this example, only INV message type TX (transaction) will be in the log.

You can fetch transaction detail by calling `getrawtransaction` JSON-RPC method.  
You can use [neo-utils](https://github.com/O3Labs/neo-utils) to parse the invocation script to get information about the transaction like script hash, smart contract's method and params.


#### Things that you can build when you have this  
- a notification service to notify user by email about the transaction
- watch particular smart contract
- watch particular NEO address
- etc.
