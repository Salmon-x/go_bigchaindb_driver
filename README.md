## go-bigchaindb-driver

Go BigchainDB driver，This is an implementation of a BigchainDB client written in golang.

## Installation
Run `go get -u github.com/Salmon-x/go_bigchaindb_driver`

## example

#### Create a transaction
```go
package main

import (
	"crypto/ed25519"
	"fmt"
	cl "github.com/Salmon-x/go_bigchaindb_driver/pkg/client"
	txn "github.com/Salmon-x/go_bigchaindb_driver/pkg/transaction"
)

func main()  {
	// create a client
	baseUrl := fmt.Sprintf("http://127.0.0.1:9992/api/v1/")
	cfg := cl.ClientConfig{
		Url: baseUrl,
	}
	client, err := cl.New(cfg)
	if err != nil {
		panic(err)
	}
    // create a asset
	asset := txn.Asset{
		Data: map[string]interface{}{
			"bicycle": map[string]interface{}{
				"serial_number": "1234",
				"manufacturer":  "1234",
			},
		},
	}
	// create a metadata
	metadata := txn.Metadata{
		"meta1": "1",
	}
	// create a keypair
	alice, err := txn.NewKeyPair()
	if err != nil {
		panic(err)
	}
	// create a condition
	conition, err := txn.NewCondition(alice.PublicKey)
	output, err := txn.NewOutput(*conition, "1")
	if err != nil {
		panic(err)
	}
	// create a transaction
	tx, err := txn.NewCreateTransaction(asset, metadata, []txn.Output{output}, []ed25519.PublicKey{alice.PublicKey})
	if err != nil {
		panic(err)
	}
	// Sign the transaction
	err = tx.Sign([]*txn.KeyPair{key})
	if err != nil {
		panic(err)
	}
	// send transaction
	err = client.PostTransaction(tx)
	if err != nil {
		panic(err)
	}
}

```

At this point, a transaction creation has been completed.

#### transfer a transaction
```go
package main

func main()  {
        // ... Refer to the previous text
    // Create a keypair again. It will be the new owner
	bob, err := txn.NewKeyPair()
	if err != nil {
		panic(err)
	}
	// Create a new condition. Bob is the new owner
	newCondition, err := txn.NewCondition(bob.PublicKey)
	if err != nil {
		panic(err)
	}
	newOutput, err := txn.NewOutput(*newCondition, "1")
	tranTx, err := txn.NewTransferTransaction([]txn.Transaction{*tx}, []txn.Output{newOutput}, metadata)
	if err != nil {
		panic(err)
	}
	// Alice will sign to confirm the transfer of assets to Bob
	err = tranTx.Sign([]*txn.KeyPair{alice})
	if err != nil {
		panic(err)
	}
	err = client.PostTransaction(tranTx)
	if err != nil {
		panic(err)
	}
}
```