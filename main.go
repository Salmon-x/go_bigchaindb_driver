package main

import (
	"crypto/ed25519"
	"fmt"
	cl "go_bdb_driver/pkg/client"
	txn "go_bdb_driver/pkg/transaction"
)

func main() {
	baseUrl := fmt.Sprintf("http://127.0.0.1:9992/api/v1/")
	cfg := cl.ClientConfig{
		Url: baseUrl,
	}
	client, err := cl.New(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	asset := txn.Asset{
		Data: map[string]interface{}{
			"bicycle": map[string]interface{}{
				"serial_number": "1234",
				"manufacturer":  "1234",
			},
		},
	}
	metadata := txn.Metadata{
		"meta1": "1",
	}
	key, err := txn.NewKeyPair()
	if err != nil {
		fmt.Println(err)
		return
	}
	conition, err := txn.NewCondition(key.PublicKey)
	output, err := txn.NewOutput(*conition, "1")
	if err != nil {
		fmt.Println(err)
		return
	}
	newBdb, err := txn.NewCreateTransaction(asset, metadata, []txn.Output{output}, []ed25519.PublicKey{key.PublicKey})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = newBdb.Sign([]*txn.KeyPair{key})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = client.PostTransaction(newBdb)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("创建成功")
	// 转移链
	key2, err := txn.NewKeyPair()
	if err != nil {
		fmt.Println(err)
		return
	}
	newConition, err := txn.NewCondition(key2.PublicKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	newOutput, err := txn.NewOutput(*newConition, "1")
	tranTx, err := txn.NewTransferTransaction([]txn.Transaction{*newBdb}, []txn.Output{newOutput}, metadata)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = tranTx.Sign([]*txn.KeyPair{key})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = client.PostTransaction(tranTx)
	if err != nil {
		fmt.Println(err)
		return
	}
}
