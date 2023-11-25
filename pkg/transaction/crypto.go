package transaction

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"go_bdb_driver/pkg/base58"
	"go_bdb_driver/pkg/base64url"
	"golang.org/x/crypto/sha3"
	"strconv"

	"crypto"
	"github.com/go-interledger/cryptoconditions"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ed25519"
	"strings"
)

type KeyPair struct {
	PrivateKey ed25519.PrivateKey `json:"privateKey"`
	PublicKey  ed25519.PublicKey  `json:"publicKey"`
}

// TODO add configurable seed to GenerateKey
func NewKeyPair() (*KeyPair, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)

	if err != nil {
		return nil, errors.Wrap(err, "Could not generate new ED25519 KeyPair")
	}

	return &KeyPair{
		PublicKey:  pubKey,
		PrivateKey: privKey,
	}, nil

}

func NewCondition(key ed25519.PublicKey) (*condition, error) {
	fingerprintContent, err := fingerprintContents(key)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(fingerprintContent)
	conitionUri := fmt.Sprintf("ni:///sha-256;%s?fpt=%s&cost=%s", base64url.Base64RemovePadding(base64url.Encode(hash[:])),
		ED25519Sha256, cost)
	res := condition{
		typeId: CTEd25519Sha256,
		fpt:    ED25519Sha256,
		cost:   cost,
		uri:    conitionUri,
		pub:    base58.BytesToB58(key),
	}
	return &res, nil
}

func fingerprintContents(publicKey []byte) ([]byte, error) {
	base := []byte{
		48, 34, 128, 32,
	}
	base = append(base, publicKey...)
	return base, nil
}

func (t *Transaction) Sign(keyPairs []*KeyPair) error {

	t.ID = nil

	signedTx := *t

	for idx, input := range signedTx.Inputs {
		var serializedTxn strings.Builder
		s, err := t.String()
		if err != nil {
			return err
		}
		serializedTxn.WriteString(s)

		keyPair := keyPairs[idx]

		if input.Fulfills != nil {
			serializedTxn.WriteString(input.Fulfills.TransactionID)
			serializedTxn.WriteString(strconv.Itoa(input.Fulfills.OutputIndex))
		}

		bytes_to_sign := []byte(serializedTxn.String())
		h3_256 := sha3.New256()
		h3_256.Write(bytes_to_sign)
		h3_256Hash := h3_256.Sum(nil)

		signature, err := keyPair.PrivateKey.Sign(rand.Reader, h3_256Hash, crypto.Hash(0))

		ed25519Fulfillment, err := cryptoconditions.NewEd25519Sha256(keyPair.PublicKey, signature)
		if err != nil {
			return errors.Wrap(err, "Could not create fulfillment")
		}

		ff, err := ed25519Fulfillment.Encode()
		if err != nil {
			return err
		}
		ffSt := base64url.Encode(ff)
		signedTx.Inputs[idx].Fulfillment = &ffSt

	}
	id, err := signedTx.createID()
	if err != nil {
		return errors.Wrap(err, "Could not create ID")
	}
	t.Inputs = signedTx.Inputs
	t.ID = &id

	return nil
}

func (t *Transaction) _removeSignatures() Transaction {
	tn := *t
	for i, input := range tn.Inputs {
		var inp Input
		inp.Fulfills = input.Fulfills
		inp.OwnersBefore = input.OwnersBefore
		tn.Inputs[i] = inp
	}
	return tn
}
