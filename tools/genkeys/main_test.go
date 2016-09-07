package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"io/ioutil"
	"testing"
)

func TestGenKeys(t *testing.T) {
	name := "test_rsa"
	if err := genkeys(name); err != nil {
		t.Fatalf("could not generate keys: %s", err)
	}

	rawPriv, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatalf("could not read private key: %s", err)
	}
	priv, err := x509.ParsePKCS1PrivateKey(rawPriv)
	if err != nil {
		t.Fatalf("cannot parse private key: %s", err)
	}

	rawPub, err := ioutil.ReadFile(name + ".pub")
	if err != nil {
		t.Fatalf("could not read public key: %s", err)
	}
	interm, err := x509.ParsePKIXPublicKey(rawPub)
	if err != nil {
		t.Fatalf("cannot parse public key: %s", err)
	}
	pub, ok := interm.(*rsa.PublicKey)
	if !ok {
		t.Fatalf("cannot handle key type %T", interm)
	}

	h := crypto.SHA256.New()
	h.Write([]byte("hello world"))
	signed, err := rsa.SignPSS(rand.Reader, priv, crypto.SHA256, h.Sum(nil), &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash})
	if err != nil {
		t.Fatalf("cannot sign data: %s", err)
	}

	err = rsa.VerifyPSS(pub, crypto.SHA256, h.Sum(nil), signed, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash})
	if err != nil {
		t.Fatalf("could not verify data: %s", err)
	}
}
